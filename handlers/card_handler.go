package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/malachi190/paylode-backend/config"
	"github.com/malachi190/paylode-backend/models"
	"github.com/malachi190/paylode-backend/service"
)

type CardRequestBody struct {
	CardHolderName string `json:"card_holder" binding:"required"`
	CardNumber     string `json:"card_number" binding:"required,max=16"`
	ExpiryDate     string `json:"expiry_date" binding:"required"`
	CVV            string `json:"cvv" binding:"required"`
}

type FundWalletRequestBody struct {
	Amount float64 `json:"amount" binding:"required"`
	Token  string  `json:"token" binding:"required"`
}

func GetCardBrand(pan string) string {
	if len(pan) < 4 {
		return "unknown"
	}

	d1 := int(pan[0] - '0')

	d2, err := strconv.Atoi(pan[:2])

	if err != nil {
		return err.Error()
	}

	d6, err := strconv.Atoi(pan[:6])

	if err != nil {
		return err.Error()
	}

	if d1 == 4 {
		return "visa"
	}

	if (222100 <= d6 && d6 <= 272099) || (51 <= d2 && d2 <= 55) {
		return "mastercard"
	}

	return "unknown"
}

func AddCard(d *config.Deps) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqBody CardRequestBody

		if err := ctx.ShouldBindJSON(&reqBody); err != nil {
			config.ErrorLogger.Printf("error while binding request: %v\n", err.Error())

			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Something went wrong, please try again",
			})
			return
		}

		user, err := config.GetLoggedInUser(ctx)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}

		expDate := strings.Split(reqBody.ExpiryDate, "/")
		expMonth := expDate[0]
		expYear := expDate[1]
		brandName := GetCardBrand(reqBody.CardNumber)

		// PASS CARD TO SERVICE
		token, lastFour, err := service.AddCard(reqBody.CardNumber, brandName, expMonth, expYear, reqBody.CVV)

		if err != nil {
			config.ErrorLogger.Printf("error while adding card: %v\n", err.Error())

			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Failed to add card, please try again",
			})
			return
		}

		// STORE TO DB
		card := models.Card{
			UserID:         user.ID,
			Brand:          brandName,
			LastFourDigits: lastFour,
			Token:          token,
			ExpiryMonth:    expMonth,
			ExpiryYear:     expYear,
		}

		if err := d.Models.Cards.CreateCard(&card); err != nil {
			config.ErrorLogger.Printf("error while adding card: %v\n", err.Error())

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error while adding card, please try again",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Card added successfully",
			"token":   token,
		})
	}
}

func FundWalletWithCard(d *config.Deps) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqBody FundWalletRequestBody

		if err := ctx.ShouldBindJSON(&reqBody); err != nil {
			config.ErrorLogger.Printf("error while binding request: %v\n", err.Error())

			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Something went wrong, please try again",
			})
			return
		}

		user, err := config.GetLoggedInUser(ctx)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}

		// Validate token
		ok, err := d.Models.Cards.ValidateCardToken(reqBody.Token)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Card does not exist",
			})
			return
		}

		transactionRef, err := service.ChargeCard(reqBody.Token)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		// FUND USER WALLET
		wallet, err := d.Models.Wallets.Fund(user.ID, reqBody.Amount)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to fund wallet, please try again",
			})
			return
		}

		// STORE TRANSACTION
		transaction := models.Transaction{
			UserID:               user.ID,
			TransactionType:      "credit",
			TransactionReference: transactionRef,
			Amount:               reqBody.Amount,
			PaymentMethod:        "card",
			Status:               "success",
		}

		if err := d.Models.Transactions.CreateTransaction(&transaction); err != nil {
			config.ErrorLogger.Printf("failed to create transaction: %v\n", err.Error())

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Wallet funding succesfull",
			"wallet":  wallet,
		})
	}
}

func FetchCards(d *config.Deps) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := config.GetLoggedInUser(ctx)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}

		// Get saved cards
		cards, err := d.Models.Cards.GetCards(user.ID)

		if err != nil {
			config.ErrorLogger.Printf("error while fetching cards: %v\n", err.Error())

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Cards fetched successfully",
			"cards":   cards,
		})
	}
}
