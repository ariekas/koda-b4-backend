package controller

import (
	"back-end-coffeShop/lib/config"
	"back-end-coffeShop/lib/middelware"
	"back-end-coffeShop/models"
	"back-end-coffeShop/respository"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matthewhartstonge/argon2"
)

type AuthController struct {
	Pool *pgxpool.Pool
}

// Register godoc
// @Summary Register a new user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "User registration data"
// @Success 201 {object} models.Response "Success register"
// @Failure 400 {object} models.Response "Bad request"
// @Failure 500 {object} models.Response "Internal server error"
// @Router /register [post]
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
		Data:    user,
	})
}

// Login godoc
// @Summary Login user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login credentials"
// @Success 201 {object} models.Response "Login success"
// @Failure 401 {object} models.Response "Wrong password"
// @Failure 404 {object} models.Response "User not found"
// @Failure 500 {object} models.Response "Internal server error"
// @Router /login [post]
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
		Data: gin.H{
			"token": token,
			"role":  users.Role,
		},
	})
}

func (ac AuthController) ForgetPassword(ctx *gin.Context) {
	var Input struct {
		Email string `json:"email"`
	}

	if err := ctx.BindJSON(&Input); err != nil {
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

func (ac AuthController) VerifCodeOtp(ctx *gin.Context) {
	var Input struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	if err := ctx.BindJSON(&Input); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid JSON",
		})
		return
	}

	if Input.Email == "" || Input.OTP == "" {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "email & otp are required",
		})
		return
	}

	data, exists := models.OtpForget[Input.Email]
	if !exists {
		ctx.JSON(404, models.Response{
			Success: false,
			Message: "OTP not found",
		})
		return
	}

	if time.Now().After(data.ExpiresAt) {
		delete(models.OtpForget, Input.Email)
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "OTP expired",
		})
		return
	}

	if data.Code != Input.OTP {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "Invalid OTP",
		})
		return
	}

	data.Verified = true
	models.OtpForget[Input.Email] = data

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "OTP verified successfully",
	})
}

func (ac AuthController) CreateNewPassword(ctx *gin.Context) {
	var Input struct {
		Email       string `json:"email"`
		NewPassword string `json:"new_password"`
	}

	if err := ctx.ShouldBindJSON(&Input); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid JSON",
		})
		return
	}

	if Input.Email == "" || Input.NewPassword == "" {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "email & new_password are required",
		})
		return
	}

	otpData, exists := models.OtpForget[Input.Email]
	if !exists || !otpData.Verified {
		ctx.JSON(403, models.Response{
			Success: false,
			Message: "OTP is not verified",
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
		Message: "Password updated successfully",
	})
}