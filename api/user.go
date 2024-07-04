package api

import (
	"fmt"
	"net/http"
	"paysyncevets/models"
	"paysyncevets/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8"`
	Email    string `json:"email" validate:"required,email"`
	RoleName string `json:"role_name" validate:"required,oneof=ADMIN ARTIST PROMOTER NORMAL"`
}
type LoginRequest struct {
    Identifier string `json:"identifier" validate:"required"` // Can be either username or email
	Password string `json:"password" validate:"required"`
}

type loginUserResponse struct {
	User models.User `json:"user"`
}

func (server *Server) createUser(ctx *gin.Context) {
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

	// depending on the type of user (admin, artist, promoter, normal)

	switch userRequest.RoleName {
	case string(models.RoleArtist):
		// create artist
		artist := models.Artist{
			UserID:     user.ID,
			ArtistName: userRequest.Username,
			BookingFee: 100,
		}
		if err := server.db.Create(&artist).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	case string(models.RolePromoter):
		// create promoter
		promoter := models.Promoter{
			UserID:      user.ID,
			CompanyName: userRequest.Username,
		}
		if err := server.db.Create(&promoter).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	var createdUser models.User

	if err := server.db.Where("id = ?", user.ID).First(&createdUser).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"user": createdUser})
}


func (server *Server) userLogin(ctx *gin.Context) {
	var loginRequest LoginRequest
	fmt.Println(loginRequest)

	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// get user by username or email
	var user models.User
	if err := server.db.Where("email = ? OR username = ?", loginRequest.Identifier, loginRequest.Identifier).First(&user).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            ctx.JSON(http.StatusNotFound, gin.H{"error": "No user found with that email or username"})
        } else {
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query user"})
        }
        return
    }
	// check if password is correct
	if err := utils.CheckPassword(loginRequest.Password, user.HashedPassword); err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	resp := loginUserResponse{
		User: user,
	}
	ctx.JSON(http.StatusOK,resp)
}
