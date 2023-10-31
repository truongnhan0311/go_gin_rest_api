package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"warranty-api/database"
	"warranty-api/forms"
	"warranty-api/models"
)

type DailyRepo struct {
	Db *gorm.DB
}

func DailyController() *DailyRepo {
	db := database.InitDb()
	db.AutoMigrate(&models.Daily{})
	return &DailyRepo{Db: db}
}

//create user
func (repository *DailyRepo) Create(c *gin.Context) {
	var data forms.CreateDailyCommand
	c.BindJSON(&data)
	daily := models.Daily{
		Name:     data.Name,
		MaSoDL:   data.MaSoDL,
		Address:  data.Address,
		Total:    0,
		AgencyID: data.AgencyID,
		Province: data.Province}

	err := models.CreateDaily(repository.Db, &daily)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, data)
}

//get list daily
func (repository *DailyRepo) GetList(c *gin.Context) {
	var data []models.Daily
	agencyID, _ := c.Get("user-agency")
	var agency = ""
	if agencyID == nil {
		agency = ""
	} else {
		agency = fmt.Sprintf("%v", agencyID)
	}
	err := models.GetListDaily(repository.Db, &data, agency)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

//get user by id
func (repository *DailyRepo) Get(c *gin.Context) {
	id, _ := c.Params.Get("id")
	var data models.Daily
	err := models.GetDaily(repository.Db, &data, id)
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
func (repository *DailyRepo) Update(c *gin.Context) {
	var data models.Daily
	id, _ := c.Params.Get("id")
	err := models.GetDaily(repository.Db, &data, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.BindJSON(&data)
	err = models.UpdateDaily(repository.Db, &data)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, data)
}

// delete warranty
func (repository *DailyRepo) Delete(c *gin.Context) {
	var data models.Daily
	id, _ := c.Params.Get("id")
	err := models.DeleteDaily(repository.Db, &data, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Daily deleted successfully"})
}
