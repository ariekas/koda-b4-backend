package controller

import (
	"back-end-coffeShop/lib/config"
	"back-end-coffeShop/lib/middelware"
	"back-end-coffeShop/models"
	"back-end-coffeShop/respository"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type AuthController struct {
	Conn *pgx.Conn
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
	user := respository.Register(ctx, ac.Conn)

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

	users, err := respository.FindUserEmail(ac.Conn, loginData.Email)

	if err != nil {
		ctx.JSON(404, models.Response{
			Success: false,
			Message: "Not found users",
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

	token, err := middelware.GenerateToken(jwtToken, users.Role)
	if err != nil {
		fmt.Println("Error: Failed to generate token")
	} 

	
	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Login success",
		Data: fmt.Sprintf("Token Login : %s", token),
	})
}
