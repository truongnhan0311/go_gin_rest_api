package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"warranty-api/database"
	"warranty-api/models"
)

type AgencyRepo struct {
	Db *gorm.DB
}

func AgencyController() *AgencyRepo {
	db := database.InitDb()
	db.AutoMigrate(&models.Agency{})
	return &AgencyRepo{Db: db}
}

func (repository *AgencyRepo) Create(c *gin.Context) {
	var data models.Agency
	c.BindJSON(&data)
	err := models.Create(repository.Db, &data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, data)
}

//get users
func (repository *AgencyRepo) GetList(c *gin.Context) {
	var data []models.Agency
	agencyID, _ := c.Get("user-agency")
	var agency = ""
	if agencyID == nil {
		agency = ""
	} else {
		agency = fmt.Sprintf("%v", agencyID)
	}

	err := models.GetList(repository.Db, &data, agency)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

//get user by id
func (repository *AgencyRepo) Get(c *gin.Context) {
	id, _ := c.Params.Get("id")
	var data models.Agency
	err := models.Get(repository.Db, &data, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

// update user
func (repository *AgencyRepo) Update(c *gin.Context) {
	var data models.Agency
	id, _ := c.Params.Get("id")
	err := models.Get(repository.Db, &data, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.BindJSON(&data)
	err = models.Update(repository.Db, &data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, data)
}

// delete warranty
func (repository *AgencyRepo) Delete(c *gin.Context) {
	var data models.Agency
	id, _ := c.Params.Get("id")
	err := models.Delete(repository.Db, &data, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
