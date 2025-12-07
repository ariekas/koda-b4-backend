package controller

import (
	"back-end-coffeShop/internal/middelware"
	"back-end-coffeShop/internal/models"
	"back-end-coffeShop/internal/respository"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserController struct {
	Pool *pgxpool.Pool
}

func (uc UserController) GetUsers(ctx *gin.Context) {
	pageQuery := ctx.Query("page")
	page := 1
	if pageQuery != "" {
		p, err := strconv.Atoi(pageQuery)
		if err == nil && p > 0 {
			page = p
		}
	}

	users, err := respository.GetDataUsers(uc.Pool, page)

	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Failed to getting data users",
		})
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Success getting users data",
		Result:  users,
	})
}

func (uc UserController) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	userId, _ := strconv.Atoi(id)
	err := respository.DeleteUser(uc.Pool, userId)

	if err != nil {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Success deleted",
	})
}

func (uc UserController) UpdateRole(ctx *gin.Context) {
	id := ctx.Param("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Error : failed to converd type data")
	}

	err = ctx.ShouldBindJSON(&models.InputNewRole)

	if err != nil {
		fmt.Println("Error: failed type request ", err)
	}

	err = respository.UpdateRole(uc.Pool, userId, models.InputNewRole.Role)

	if err != nil {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Success update Role user",
	})
}

func (uc UserController) GetUserLogin(ctx *gin.Context) {
	userId := middelware.GetUserFromToken(ctx)

	if userId == 0 {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "Unauthorized: Invalid or missing token",
		})
		return
	}

	user, err := respository.GetUserByToken(uc.Pool, userId)
	if err != nil {
		ctx.JSON(404, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Success getting user",
		Result:  user,
	})
}

func (uc UserController) UpdateProfile(ctx *gin.Context) {
	userId := middelware.GetUserFromToken(ctx)

	userData := ctx.PostForm("data")
	var input models.UpdateProfileRequest
	if userData != "" {
		if err := json.Unmarshal([]byte(userData), &input); err != nil {
			ctx.JSON(400, models.Response{
				Success: false,
				Message: "Invalid JSON input",
			})
			return
		}
	}

	file, err := ctx.FormFile("pic")
	if err == nil {
		input.PicFile = file
	}

	err = respository.UpdateProfile(uc.Pool, userId, input)
	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: fmt.Sprintf("Failed to update profile: %v", err),
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Profile updated successfully",
	})
}
