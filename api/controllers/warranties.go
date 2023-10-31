package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
	"time"
	"warranty-api/database"
	"warranty-api/forms"
	"warranty-api/models"
	"warranty-api/services"
)

type WarrantyRepo struct {
	Db *gorm.DB
}

func WarrantyController() *WarrantyRepo {
	db := database.InitDb()
	db.AutoMigrate(&models.Warranty{})
	return &WarrantyRepo{Db: db}
}

//create user
func (repository *WarrantyRepo) CreateWarranty(c *gin.Context) {
	var data forms.CreateWarrantyCommand
	if c.BindJSON(&data) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Provide Relevant Field"})
		return
	}
	jsonVitridan, _ := json.Marshal(data.ViTriDan)
	userID, _ := c.Get("user-id")
	unique := uuid.New()
	ngaydan, _ := time.Parse("2006-01-02", data.NgayDan)

	warranty := models.Warranty{
		Name:     data.Name,
		Phone:    data.Phone,
		Soxe:     data.SoXe,
		Tenxe:    data.TenXe,
		LoaiKinh: strings.Join(data.LoaiKinh, ","),
		NgayDan:  ngaydan,
		AgencyID: data.Agency,
		DailyID:  data.DaiLy,
		UserID:   userID.(int),
		MaSoBH:   unique.String(),
		ViTriDan: string(jsonVitridan),
		DoiXe:    data.DoiXe,
		MauXe:    data.MauXe,
		LoaiXe:   data.LoaiXe,
		Email:    data.Email,
		Note:     data.Note,
		Status:   "active"}
	err := models.CreateWarranty(repository.Db, &warranty)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var agency models.Agency
	var daily models.Daily
	_ = models.Get(repository.Db, &agency, strconv.Itoa(data.Agency))
	_ = models.GetDaily(repository.Db, &daily, strconv.Itoa(data.DaiLy))

	warranty.MaSo = daily.Province + daily.MaSoDL + "_" + agency.MaSoDL + "_" + fmt.Sprintf(warranty.NgayDan.Format("01022006")) + "_" + fmt.Sprintf("%06d", warranty.ID)

	err = models.UpdateWarranty(repository.Db, &warranty)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	_ = models.AddTotal(repository.Db, &agency, strconv.Itoa(int(agency.ID)))
	_ = models.AddDailyTotal(repository.Db, &daily, strconv.Itoa(int(daily.ID)))

	c.JSON(http.StatusOK, warranty)
}

//get list
func (repository *WarrantyRepo) GetWarrantyList(c *gin.Context) {
	var request = parseRequest(c)
	var data []models.Warranty
	err := models.GetWarrantyList(repository.Db, &data, request.Agency, request.Daily, request.Phone, request.Soxe, request.Maso, request.Ngaydanfrom, request.Ngaydanto)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (repository *WarrantyRepo) ExportWarrantyList(c *gin.Context) {
	var request = parseRequest(c)
	var data []models.Warranty
	err := models.GetWarrantyList(repository.Db, &data, request.Agency, request.Daily, request.Phone, request.Soxe, request.Maso, request.Ngaydanfrom, request.Ngaydanto)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	dowloadUrl := services.GenerateExcel(data)
	c.JSON(http.StatusOK, gin.H{"data": dowloadUrl})
}

func parseRequest(c *gin.Context) forms.Request {
	q := c.Request.URL.Query()

	var request forms.Request
	request.Phone = ""
	request.Soxe = ""
	request.Maso = ""
	request.Agency = ""
	request.Daily = ""
	request.Ngaydanfrom = ""
	request.Ngaydanto = ""

	if len(q["phone"]) != 0 {
		request.Phone = fmt.Sprintf("%v", q["phone"][0])
	}

	if len(q["soxe"]) != 0 {
		request.Soxe = fmt.Sprintf("%v", q["soxe"][0])
	}

	if len(q["maso"]) != 0 {
		request.Maso = fmt.Sprintf("%v", q["maso"][0])
	}

	if len(q["ngay_dan_from"]) != 0 {
		request.Ngaydanfrom = fmt.Sprintf("%v", q["ngay_dan_from"][0])
	}

	if len(q["ngay_dan_to"]) != 0 {
		request.Ngaydanto = fmt.Sprintf("%v", q["ngay_dan_to"][0])
	}

	agencyID, _ := c.Get("user-agency")
	if agencyID == nil {
		if len(q["agency"]) != 0 {
			request.Agency = fmt.Sprintf("%v", q["agency"][0])
		}
	} else {
		request.Agency = fmt.Sprintf("%v", agencyID)
	}

	dailyID, _ := c.Get("user-daily")
	if dailyID == nil {
		if len(q["daily"]) != 0 {
			request.Daily = fmt.Sprintf("%v", q["daily"][0])
		}
	} else {
		request.Daily = fmt.Sprintf("%v", dailyID)
	}

	return request
}

func (res *WarrantyRepo) SearchWarranty(c *gin.Context) {
	q := c.Request.URL.Query()
	phone, soxe, maso := "", "", ""
	if len(q["phone"]) != 0 {
		phone = fmt.Sprintf("%v", q["phone"][0])
	}

	if len(q["soxe"]) != 0 {
		soxe = fmt.Sprintf("%v", q["soxe"][0])
	}

	if len(q["maso"]) != 0 {
		maso = fmt.Sprintf("%v", q["maso"][0])
	}

	var data []models.Warranty
	err := models.GetWarrantyList(res.Db, &data, "", "", phone, soxe, maso, "", "")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

//get user by id
func (repository *WarrantyRepo) GetWarranty(c *gin.Context) {
	id, _ := c.Params.Get("id")
	agencyID, _ := c.Get("user-agency")

	var data models.Warranty
	var agency = ""
	if agencyID == nil {
		agency = ""
	} else {
		agency = fmt.Sprintf("%v", agencyID)
	}

	err := models.GetWarranty(repository.Db, &data, id, agency)
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

func (repository *WarrantyRepo) GetWarrantyByCode(c *gin.Context) {
	id, _ := c.Params.Get("id")

	var data models.Warranty
	err := models.GetWarrantyByCode(repository.Db, &data, id)
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

func (repository *WarrantyRepo) DownloadWarranty(c *gin.Context) {
	id, _ := c.Params.Get("id")
	var data models.Warranty
	err := models.GetWarrantyByCode(repository.Db, &data, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	dowloadUrl := services.GeneratePDF(data)
	c.JSON(http.StatusOK, gin.H{"data": dowloadUrl})

}

// update user
func (repository *WarrantyRepo) UpdateWarranty(c *gin.Context) {
	var data models.Warranty
	var requestData forms.UpdateWarrantyCommand
	if c.BindJSON(&requestData) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Provide Relevant Field"})
		return
	}

	id, _ := c.Params.Get("id")
	err := models.GetWarranty(repository.Db, &data, id, "")

	var daily models.Daily
	_ = models.GetDaily(repository.Db, &daily, strconv.Itoa(data.DailyID))
	_ = models.DeductDailyTotal(repository.Db, &daily, strconv.Itoa(int(daily.ID)))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	data.Name = requestData.Name
	data.Phone = requestData.Phone
	data.DoiXe = requestData.DoiXe
	data.MauXe = requestData.MauXe
	data.Tenxe = requestData.TenXe
	data.Email = requestData.Email
	data.DailyID = requestData.DaiLy
	data.Note = requestData.Note
	data.Daily = models.Daily{}

	if data.UserID == 0 {
		userID, _ := c.Get("user-id")
		data.UserID = userID.(int)
	}

	err = models.UpdateWarranty(repository.Db, &data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	_ = models.GetDaily(repository.Db, &daily, strconv.Itoa(data.DailyID))
	_ = models.AddDailyTotal(repository.Db, &daily, strconv.Itoa(int(data.DailyID)))
	c.JSON(http.StatusOK, data)
}

// delete warranty
func (repository *WarrantyRepo) DeleteWarranty(c *gin.Context) {
	var data models.Warranty
	var modelAgency models.Agency
	var modelDaily models.Daily
	id, _ := c.Params.Get("id")
	agency, _ := c.Params.Get("agency")
	daily, _ := c.Params.Get("daily")

	err := models.DeleteWarranty(repository.Db, &data, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	_ = models.DeductTotal(repository.Db, &modelAgency, agency)
	_ = models.DeductDailyTotal(repository.Db, &modelDaily, daily)
	c.JSON(http.StatusOK, gin.H{"message": "Warranty deleted successfully"})
}
