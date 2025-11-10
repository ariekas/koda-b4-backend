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

// GetUsers godoc
// @Summary Get all users
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.Response{data=[]models.User} "Success getting users data"
// @Failure 400 {object} models.Response "Failed to get users data"
// @Failure 401 {object} models.Response "Unauthorized"
// @Router /users [get]
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


// DeleteUser godoc
// @Summary Delete a user
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 201 {object} models.Response "User deleted successfully"
// @Failure 401 {object} models.Response "Unauthorized"
// @Failure 404 {object} models.Response "User not found"
// @Router /users/{id} [delete]
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


// UpdateRole godoc
// @Summary Update user role
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body models.InputNewRoleUser true "New Role Data"
// @Success 201 {object} models.Response "Role updated successfully"
// @Failure 400 {object} models.Response "Invalid request"
// @Failure 401 {object} models.Response "Unauthorized"
// @Router /users/role/{id} [patch]
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
