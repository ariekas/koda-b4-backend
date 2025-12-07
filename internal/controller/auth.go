package controller

import (
	"back-end-coffeShop/internal/config"
	"back-end-coffeShop/internal/middelware"
	"back-end-coffeShop/internal/models"
	"back-end-coffeShop/internal/respository"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matthewhartstonge/argon2"
)

type AuthController struct {
	Pool *pgxpool.Pool
}

func (ac AuthController) Register(ctx *gin.Context) {
	user, err := respository.Register(ctx, ac.Pool)

	if err != nil {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	
	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Success register",
		Result:    user,
	})
}

func (ac AuthController) Login(ctx *gin.Context) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := ctx.BindJSON(&loginData)
	jwtToken := config.ReadENV()

	if err != nil {
		fmt.Println("Error : Failed type much json")
	}

	if loginData.Email == "" {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "requared email",
		})
	}

	if loginData.Password == "" {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "requared password",
		})
	}

	users, err := respository.FindUserEmail(ac.Pool, loginData.Email)
	if err != nil {
		ctx.JSON(404, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if !respository.VerifPassword(users.Password, loginData.Password) {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "Wrong password",
		})
		return
	}

	token, err := middelware.GenerateToken(jwtToken, users.Role, users.Id)
	if err != nil {
		fmt.Println("Error: Failed to generate token")
	} 

	
	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Login success",
		Result: gin.H{
			"token": token,
			"role":  users.Role,
		},
	})
}
func (ac AuthController) ForgetPassword(ctx *gin.Context) {
	var Input struct {
		Email string `json:"email" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&Input); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid JSON",
		})
		return
	}

	if Input.Email == "" {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "email is required",
		})
		return
	}

	_, err := respository.FindUserEmail(ac.Pool, Input.Email)
	if err != nil {
		ctx.JSON(404, models.Response{
			Success: false,
			Message: "email not found",
		})
		return
	}

	otp := fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)

	models.OtpForget[Input.Email] = struct {
		Code      string
		ExpiresAt time.Time
		Verified  bool
	}{
		Code:      otp,
		ExpiresAt: time.Now().Add(5 * time.Minute),
		Verified:  false,
	}

	err = respository.SendOtpEmail(Input.Email, otp)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to send OTP email: " + err.Error(),
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "OTP has been sent to your email",
	})
}

func (ac AuthController) ResetPassword(ctx *gin.Context) {
	var Input struct {
		Email           string `json:"email" binding:"required"`
		OTP             string `json:"otp" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required"`
		ConfirmPassword string `json:"confirm_password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&Input); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid JSON",
		})
		return
	}

	if Input.Email == "" || Input.OTP == "" || Input.NewPassword == "" || Input.ConfirmPassword == "" {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "all fields are required",
		})
		return
	}

	if Input.NewPassword != Input.ConfirmPassword {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "password and confirm password do not match",
		})
		return
	}

	if len(Input.NewPassword) <= 6 {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "password must be more than 6 characters",
		})
		return
	}

	otpData, exists := models.OtpForget[Input.Email]
	if !exists {
		ctx.JSON(404, models.Response{
			Success: false,
			Message: "OTP not found or expired",
		})
		return
	}

	if time.Now().After(otpData.ExpiresAt) {
		delete(models.OtpForget, Input.Email)
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "OTP expired",
		})
		return
	}

	if otpData.Code != Input.OTP {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "Invalid OTP",
		})
		return
	}

	user, err := respository.FindUserEmail(ac.Pool, Input.Email)
	if err != nil {
		ctx.JSON(404, models.Response{
			Success: false,
			Message: "email not found",
		})
		return
	}

	ok, _ := argon2.VerifyEncoded([]byte(user.Password), []byte(Input.NewPassword))
	if ok {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "new password cannot be the same as old password",
		})
		return
	}

	if err := respository.UpdatePassword(ac.Pool, Input.Email, Input.NewPassword); err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	delete(models.OtpForget, Input.Email)

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Password has been reset successfully",
	})
}