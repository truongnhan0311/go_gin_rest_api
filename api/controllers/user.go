package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"warranty-api/database"
	"warranty-api/forms"
	"warranty-api/helpers"
	"warranty-api/models"
	"warranty-api/services"
)

type UserRepo struct {
	Db *gorm.DB
}

func UserController() *UserRepo {
	db := database.InitDb()
	db.AutoMigrate(&models.User{})
	return &UserRepo{Db: db}
}

func getAgency(c *gin.Context) string {
	agencyID, _ := c.Get("user-agency")
	var agency = ""
	if agencyID != nil {
		agency = fmt.Sprintf("%v", agencyID)
	}

	return agency
}

func getDaily(c *gin.Context) string {
	agencyID, _ := c.Get("user-daily")
	var agency = ""
	if agencyID != nil {
		agency = fmt.Sprintf("%v", agencyID)
	}

	return agency
}

//create user
func (repository *UserRepo) CreateUser(c *gin.Context) {
	var data forms.SignupUserCommand
	if c.BindJSON(&data) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Provide Relevant Field"})
		return
	}

	result, _ := models.GetUserByEmail(repository.Db, data.Email)
	if result.Email != "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Email is already in use"})
		return
	}

	user := models.User{
		Name:     data.Name,
		Email:    data.Email,
		AgencyID: data.Agency,
		DailyID:  data.Daily,
		Password: helpers.GeneratePasswordHash([]byte(data.Password))}
	err := models.CreateUser(repository.Db, &user)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Problem to create account"})
	}
	c.JSON(http.StatusOK, user)
}

//get users
func (repository *UserRepo) GetUsers(c *gin.Context) {
	var user []models.User
	err := models.GetUsers(repository.Db, &user, getAgency(c), getDaily(c))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

//get user by id
func (repository *UserRepo) GetMe(c *gin.Context) {
	User, _ := c.Get("User")
	c.JSON(http.StatusOK, gin.H{"user": User})
}

//get user by id
func (repository *UserRepo) GetUser(c *gin.Context) {
	id, _ := c.Params.Get("id")
	user, err := models.GetUser(repository.Db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

// update user
func (repository *UserRepo) UpdateUser(c *gin.Context) {
	var data forms.UpdateUserCommand

	id, _ := c.Params.Get("id")
	user, _ := models.GetUser(repository.Db, id)

	if c.BindJSON(&data) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Provide Relevant Field"})
		return
	}

	result, _ := models.GetUserByEmail(repository.Db, data.Email)
	if result.Email != "" && (result.ID != user.ID) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Email is already in use"})
		return
	}

	user.Email = data.Email
	user.Name = data.Name
	userID, _ := c.Get("user-id")

	if userID != user.ID {
		user.AgencyID = data.Agency
		user.Agency = models.Agency{}
	}
	user.DailyID = data.Daily
	user.Daily = models.Daily{}

	err := models.UpdateUser(repository.Db, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Problem to create account"})
	}
	c.JSON(http.StatusOK, user)
}

func (repository *UserRepo) UpdateMe(c *gin.Context) {
	var data forms.UpdateUserCommand

	if c.BindJSON(&data) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Provide Relevant Field"})
		return
	}

	userEmail, _ := c.Get("user-email")
	if userEmail != data.Email {
		result, _ := models.GetUserByEmail(repository.Db, data.Email)
		if result.Email != "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Email is already in use"})
			return
		}
	}
	user, _ := models.GetUserByEmail(repository.Db, fmt.Sprintf("%v", userEmail))
	user.Email = data.Email
	user.Name = data.Name

	err := models.UpdateUser(repository.Db, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Problem to create account"})
	}

	c.Set("user-email", user.Email)
	c.JSON(http.StatusOK, user)
}

func (repository *UserRepo) UpdatePassword(c *gin.Context) {
	var data forms.ChangePasswordCommand
	if c.BindJSON(&data) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Provide Relevant Field"})
		return
	}
	if data.Password != data.Confirm {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Provide Relevant Field"})
		return
	}
	userEmail, _ := c.Get("user-email")
	user, err := models.GetUserByEmail(repository.Db, fmt.Sprintf("%v", userEmail))
	if user.Email == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	hashedPassword := []byte(user.Password)
	password := []byte(data.CurrentPassword)

	err = helpers.PasswordCompare(password, hashedPassword)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid Credentials"})
		return
	}

	user.Password = helpers.GeneratePasswordHash([]byte(data.Password))
	updateErr := models.UpdateUser(repository.Db, user)
	if updateErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Problem to update password"})
	}
	c.JSON(http.StatusOK, user)
}

// delete user
func (repository *UserRepo) DeleteUser(c *gin.Context) {
	var user models.User
	id, _ := c.Params.Get("id")
	err := models.DeleteUser(repository.Db, &user, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// delete user
func (repository *UserRepo) Login(c *gin.Context) {
	var data forms.LoginUserCommand
	if c.BindJSON(&data) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Provide Relevant Field"})
		return
	}

	user, err := models.GetUserByEmail(repository.Db, data.Email)
	if user.Email == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Problem Login to Account"})
		return
	}

	hashedPassword := []byte(user.Password)
	password := []byte(data.Password)

	err = helpers.PasswordCompare(password, hashedPassword)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid Credentials"})
		return
	}

	jwtToken, jwtRefeshToken, err2 := services.GenerateToken(user.Email)
	if err2 != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Problem Login to Account. Try Again Later"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successfully", "token": jwtToken, "refresh_token": jwtRefeshToken})
}
