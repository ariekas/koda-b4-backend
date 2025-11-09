package controller

import (
	"back-end-coffeShop/models"
	"back-end-coffeShop/respository"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type UserController struct{
	Conn *pgx.Conn
}

func (uc UserController) GetUsers(ctx *gin.Context){
	users, err := respository.GetDataUsers(uc.Conn)

	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Failed to getting data users",
		})
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Success getting users data",
		Data: users,
	})
}

func (uc UserController) DeleteUser(ctx *gin.Context){
	err := respository.DeleteUser(uc.Conn, ctx)

	if err != nil {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "Error: Failed to delete user",
		})
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Success deleted",
	})
}

func (uc UserController) UpdateRole(ctx *gin.Context){
	id := ctx.Param("id")

	userId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Error : failed to converd type data")
	}

	err = ctx.ShouldBindJSON(&models.InputNewRole)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	err = respository.UpdateRole(uc.Conn,ctx, userId, models. InputNewRole.Role)

	if err != nil {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "Error: Failed to update role user",
		})
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Success update Role user",
	})

}
