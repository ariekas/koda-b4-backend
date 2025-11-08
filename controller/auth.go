package controller

import (
	"back-end-coffeShop/models"
	"back-end-coffeShop/respository"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type AuthController struct {
	Conn *pgx.Conn
}

func (ac AuthController) Register(ctx *gin.Context) {
	user := respository.Register(ctx, ac.Conn)

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Success register",
		Data:    user,
	})
}

func (ac AuthController) Login(ctx *gin.Context) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := ctx.BindJSON(&loginData)

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
	
	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Login success",
		Data: gin.H{
			"fullname": users.Fullname,
			"email":    users.Email,
			"role":     users.Role,
		},
	})
}
