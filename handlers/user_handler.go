package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/malachi190/paylode-backend/config"
	"github.com/malachi190/paylode-backend/types"
	"golang.org/x/crypto/bcrypt"
)

func CreatePin(d *config.Deps) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqBody types.CreatePinBody

		if err := ctx.ShouldBindJSON(&reqBody); err != nil {
			config.ErrorLogger.Printf("error while binding request: %v\n", err.Error())

			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Something went wrong, please try again",
			})
			return
		}

		// GET USER FROM CONTEXT
		user, err := config.GetLoggedInUser(ctx)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}

		// CONFIRM PIN
		if reqBody.Pin != reqBody.PinConfirmation {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Pin mismatch, please confirm pin again",
			})
			return
		}

		// HASH PIN
		hashedPin, err := bcrypt.GenerateFromPassword([]byte(reqBody.Pin), 10)

		if err != nil {
			config.ErrorLogger.Printf("error while hashing pin: %v\n", err.Error())

			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Something went wrong, please try again",
			})
			return
		}

		pin := string(hashedPin)

		// CREATE PIN FOR USER
		if err := d.Models.Users.CreateUserPin(user.ID, pin); err != nil {
			config.ErrorLogger.Printf("error while creating pin: %v\n", err.Error())

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to create pin, please try again",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Pin created successfully",
		})
	}
}

func GetUserWallet(d *config.Deps) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := config.GetLoggedInUser(ctx)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}

		wallet, err := d.Models.Wallets.GetWallet(user.ID)

		if err != nil {
			config.ErrorLogger.Printf("error while fetching wallet: %v\n", err.Error())

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Wallet fetched successfully",
			"wallet":  wallet,
		})
	}
}

func FetchUserTransactions(d *config.Deps) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := config.GetLoggedInUser(ctx)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}

		// Get user transactions
		transactions, err := d.Models.Transactions.GetTransactions(user.ID)

		if err != nil {
			config.ErrorLogger.Printf("error while fetching transactions: %v\n", err.Error())

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message":      "Transactions fetched successfully",
			"transactions": transactions,
		})
	}
}
