package api

import (
	"net/http"
	"paysyncevets/models"
	"paysyncevets/utils"

	"github.com/gin-gonic/gin"
)


type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8"`
	Email    string `json:"email" validate:"required,email"`
	RoleName string `json:"role_name" validate:"required,oneof=ADMIN ARTIST PROMOTER NORMAL"`
}

func  (server *Server) createUser(ctx *gin.Context){
	var userRequest CreateUserRequest

	// bind the request body to the userRequest

	if err := ctx.ShouldBindJSON(&userRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// hash password

	hashedPassword, err := utils.HashedPassword(userRequest.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// create user

	// roleid from the UserRole table
	var role models.UserRole

	if err := server.db.Where("role_name = ?", userRequest.RoleName).First(&role).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user := models.User{
		Username:       userRequest.Username,
		HashedPassword: hashedPassword,
		Email:          userRequest.Email,
		RoleID:         role.ID,
	}

	if err := server.db.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// create the user depending on the role
	if userRequest.RoleName == string(models.RoleArtist) {
		artist := models.Artist{UserID: user.ID}
		if err := server.db.Create(&artist).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create artist"})
			return
		}
	} else if userRequest.RoleName == string(models.RolePromoter) {
		promoter := models.Promoter{UserID: user.ID}
		if err := server.db.Create(&promoter).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create promoter"})
			return
		}
	}
	ctx.JSON(http.StatusCreated, gin.H{"user": user})
}