package controller

import (
	"back-end-coffeShop/internal/models"
	"back-end-coffeShop/internal/respository"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RatingController struct {
	Pool *pgxpool.Pool
}

func (rc *RatingController) CreateRating(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(500, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.CreateRatingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(401, gin.H{"error": err.Error()})
		return
	}

	hasPurchased, err := respository.CheckUserPurchased(rc.Pool, userID.(int), req.ProductsID, req.TransactionsID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to verify purchase"})
		return
	}

	if !hasPurchased {
		ctx.JSON(403, gin.H{"error": "You can only rate products you have purchased in this transaction"})
		return
	}

	alreadyRated, err := respository.CheckAlreadyRated(rc.Pool, userID.(int), req.ProductsID, req.TransactionsID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to check rating status"})
		return
	}

	if alreadyRated {
		ctx.JSON(400, gin.H{"error": "You have already rated this product in this transaction"})
		return
	}

	rating := &models.Rating{
		UsersID:        userID.(int),
		ProductsID:     req.ProductsID,
		TransactionsID: req.TransactionsID,
		Rating:         req.Rating,
		Review:         req.Review,
	}

	if err := respository.CreateRating(rc.Pool, rating); err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to create rating"})
		return
	}

	ctx.JSON(201, gin.H{
		"message": "Rating created successfully",
		"data":    rating,
	})
}

func (rc *RatingController) GetRatingByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(401, gin.H{"error": "Invalid rating ID"})
		return
	}

	rating, err := respository.DetailRating(rc.Pool, id)
	if err != nil {
		ctx.JSON(404, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{"data": rating})
}

func (rc *RatingController) GetRatingsByProductID(ctx *gin.Context) {
	productID, err := strconv.Atoi(ctx.Param("product_id"))
	if err != nil {
		ctx.JSON(401, gin.H{"error": "Invalid product ID"})
		return
	}

	ratings, err := respository.GetRatingProduct(rc.Pool, productID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to fetch ratings"})
		return
	}

	ctx.JSON(201, gin.H{
		"data":  ratings,
		"count": len(ratings),
	})
}

func (rc *RatingController) GetRatingsByTransactionID(ctx *gin.Context) {
	transactionID, err := strconv.Atoi(ctx.Param("transaction_id"))
	if err != nil {
		ctx.JSON(401, gin.H{"error": "Invalid transaction ID"})
		return
	}

	ratings, err := respository.GetRatingTransaction(rc.Pool, transactionID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to fetch ratings"})
		return
	}

	ctx.JSON(201, gin.H{
		"data":  ratings,
		"count": len(ratings),
	})
}

func (rc *RatingController) GetMyRatings(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(403, gin.H{"error": "User not authenticated"})
		return
	}

	ratings, err := respository.GetRatingUser(rc.Pool, userID.(int))
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to fetch ratings"})
		return
	}

	ctx.JSON(201, gin.H{
		"data":  ratings,
		"count": len(ratings),
	})
}

func (rc *RatingController) GetUnratedProducts(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(400, gin.H{"error": "User not authenticated"})
		return
	}

	transactionID, err := strconv.Atoi(ctx.Param("transaction_id"))
	if err != nil {
		ctx.JSON(401, gin.H{"error": "Invalid transaction ID"})
		return
	}

	products, err := respository.GetUnratedProductsByTransaction(rc.Pool, userID.(int), transactionID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to fetch unrated products"})
		return
	}

	ctx.JSON(201, gin.H{
		"data":  products,
		"count": len(products),
	})
}

func (rc *RatingController) UpdateRating(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(400, gin.H{"error": "User not authenticated"})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(401, gin.H{"error": "Invalid rating ID"})
		return
	}

	var req models.UpdateRatingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(401, gin.H{"error": err.Error()})
		return
	}

	existingRating, err := respository.DetailRating(rc.Pool, id)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "Rating not found"})
		return
	}

	if existingRating.UsersID != userID.(int) {
		ctx.JSON(400, gin.H{"error": "You can only update your own ratings"})
		return
	}

	rating := &models.Rating{
		ID:      id,
		UsersID: userID.(int),
		Rating:  req.Rating,
		Review:  req.Review,
	}

	if err := respository.UpdateRating(rc.Pool, rating); err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to update rating"})
		return
	}

	ctx.JSON(201, gin.H{"message": "Rating updated successfully"})
}

func (rc *RatingController) DeleteRating(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(400, gin.H{"error": "User not authenticated"})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(401, gin.H{"error": "Invalid rating ID"})
		return
	}

	rating, err := respository.DetailRating(rc.Pool, id)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "Rating not found"})
		return
	}

	if rating.UsersID != userID.(int) {
		ctx.JSON(400, gin.H{"error": "You can only delete your own ratings"})
		return
	}

	if err := respository.DeleteRatinh(rc.Pool, id); err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to delete rating"})
		return
	}

	ctx.JSON(201, gin.H{"message": "Rating deleted successfully"})
}
