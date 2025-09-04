package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/malachi190/paylode-backend/config"
	"github.com/malachi190/paylode-backend/mailer"
	"github.com/malachi190/paylode-backend/models"
	"github.com/malachi190/paylode-backend/types"
	"golang.org/x/crypto/bcrypt"
)

func SendOtpToken(d *config.Deps) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqBody types.SendOtpEmailBody

		if err := ctx.ShouldBindJSON(&reqBody); err != nil {
			config.ErrorLogger.Printf("error while binding request: %v\n", err.Error())

			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Something went wrong, please try again",
			})
			return
		}

		// GENERATE OTP TOKEN
		otp := config.GenerateOtp()

		// STORE OTP IN REDIS
		c := ctx.Request.Context()

		err := d.Redis.Set(c, reqBody.Email, otp, 10*time.Minute).Err()

		if err != nil {
			config.ErrorLogger.Printf("error setting otp in redis: %v\n", err)

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to send OTP, please try again",
			})
			return
		}

		otpStr := strconv.Itoa(otp)

		m := mailer.New()

		// SEND EMAIL
		if err := m.Send(reqBody.Email, "Email Verification", map[string]string{
			"OTP": otpStr,
		}); err != nil {
			config.ErrorLogger.Printf("mailer error: %v\n", err.Error())

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong, please try again",
			})
			return
		}

		config.GeneralLogger.Println("OTP token sent!")

		ctx.JSON(http.StatusOK, gin.H{
			"message": "OTP token sent!",
		})
	}
}

func VerifyEmail(d *config.Deps) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqBody types.VerifyEmailBody

		if err := ctx.ShouldBindJSON(&reqBody); err != nil {
			config.ErrorLogger.Printf("error binding request: %v\n", err.Error())

			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Something went wrong, please try again",
			})
			return
		}

		// QUERY REDIS STORE
		c := ctx.Request.Context()
		value, err := d.Redis.Get(c, reqBody.Email).Result()

		if err != nil {
			config.ErrorLogger.Printf("error getting value from redis: %v\n", err.Error())

			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Something went wrong",
			})
			return
		}

		if reqBody.Otp != value {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid OTP token provided",
			})
			return
		}

		config.GeneralLogger.Println("Email successfully verified")

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Email verification successful",
		})
	}
}

func Register(d *config.Deps) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqBody types.RegisterBody

		if err := ctx.ShouldBindJSON(&reqBody); err != nil {
			config.ErrorLogger.Printf("error binding request body: %v\n", err.Error())

			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			return
		}

		// HASH PASSWORD
		hashed, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), 10)

		if err != nil {
			config.ErrorLogger.Printf("error generating hashed password: %v\n", err.Error())

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong, please try again",
			})
			return
		}

		// STORE USER
		user := models.User{
			FirstName:       reqBody.FirstName,
			LastName:        reqBody.LastName,
			Email:           reqBody.Email,
			EmailVerifiedAt: time.Now().Format(time.RFC3339),
			PhoneNumber:     reqBody.PhoneNumber,
			Password:        string(hashed),
		}

		res, err := d.Models.Users.CreateUser(&user)

		if err != nil {
			config.ErrorLogger.Printf("error creating user: %v\n", err.Error())

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to create account, please try again",
			})
			return
		}

		// CREATE WALLET FOR USER
		walletIdInt := config.GenerateWalletID()

		walletId := strconv.Itoa(walletIdInt)

		wallet := models.Wallet{
			UserID:        res.ID,
			WalletBalance: 0,
			WalletID:      walletId,
		}

		if err := d.Models.Wallets.CreateWallet(&wallet); err != nil {
			config.ErrorLogger.Printf("error creating wallet: %v\n", err.Error())

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong",
			})
			return
		}

		// RETURN RESPONSE
		ctx.JSON(http.StatusCreated, gin.H{
			"message": "Account created successfully",
		})
	}
}

func Login(d *config.Deps) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqBody types.LoginBody

		if err := ctx.ShouldBindJSON(&reqBody); err != nil {
			config.ErrorLogger.Printf("error binding request body: %v\n", err.Error())

			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong... please try again"})
			return
		}

		// Get user with email or phone number
		user, err := d.Models.Users.GetUserWithEmailOrPhone(reqBody.Email, reqBody.PhoneNumber)

		if err != nil {
			ctx.JSON(http.StatusConflict, gin.H{
				"message": "Invalid email or phone number",
			})
			return
		}

		// Compare password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password))

		if err != nil {
			ctx.JSON(http.StatusConflict, gin.H{
				"message": "Invalid email or password",
			})
			return
		}

		// Generate token
		token, err1 := config.GenerateAuthToken(user.ID)
		refresh, err2 := config.GenerateRefreshToken(user.ID)

		if err1 != nil || err2 != nil {
			config.ErrorLogger.Printf("error generating token: %v\n", err)

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong",
			})
			return
		}

		// HASH REFRESH TOKEN
		hash := sha256.Sum256([]byte(refresh))
		hashEncode := hex.EncodeToString(hash[:])

		// SAVE REFRESH TOKEN
		if err := d.Models.Sessions.SaveRefreshToken(hashEncode, user.ID); err != nil {
			config.ErrorLogger.Printf("error saving token: %v\n", err)

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message":       "Login successful",
			"user":          user,
			"access_token":  token,
			"refresh_token": refresh,
		})
	}
}

func Refresh(d *config.Deps) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqBody types.RefreshTokenBody

		if err := ctx.ShouldBindJSON(&reqBody); err != nil {
			config.ErrorLogger.Printf("error binding request body: %v\n", err.Error())

			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong... please try again"})
			return
		}

		claims, err := d.ValidateRefreshToken(reqBody.Token)

		if err != nil {
			config.ErrorLogger.Printf("invalid refresh token: %v\n", err)

			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "invalid refresh token"})
			return
		}

		// CREATE NEW ACCESS TOKEN
		newAccessToken, err1 := config.GenerateAuthToken(claims["user_id"].(float64))
		newRefreshToken, err2 := config.GenerateRefreshToken(claims["user_id"].(float64))

		if err1 != nil || err2 != nil {
			config.ErrorLogger.Printf("error generating token: %v\n", err)

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"access_token":  newAccessToken,
			"refresh_token": newRefreshToken,
		})
	}
}
