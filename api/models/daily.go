package models

import (
	"gorm.io/gorm"
	"time"
)

type Daily struct {
	gorm.Model
	ID        uint   `json:"id" gorm:"primary_key"`
	Name      string `json:"name"`
	MaSoDL    string `json:"ma_so_dl"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	AgencyID  int
	Total     int    `json:"total" gorm:"default:0"`
	Province  string `json:"province"`
	Status    string `json:"status"`
	Agency    Agency
	CreatedAt time.Time
	UpdatedAt time.Time
	Deleted   gorm.DeletedAt
}

func CreateDaily(db *gorm.DB, data *Daily) (err error) {
	err = db.Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func GetListDaily(db *gorm.DB, data *[]Daily, agency string) (err error) {
	if agency == "" {
		err = db.Preload("Agency").Find(data).Error
	} else {
		err = db.Preload("Agency").Where("agency_id = ?", agency).Find(data).Error
	}

	if err != nil {
		return err
	}
	return nil
}

func GetDaily(db *gorm.DB, data *Daily, id string) (err error) {
	err = db.Preload("Agency").Where("id = ?", id).First(data).Error
	if err != nil {
		return err
	}
	return nil
}

//update user
func UpdateDaily(db *gorm.DB, data *Daily) (err error) {
	db.Save(data)
	return nil
}

//delete user
func DeleteDaily(db *gorm.DB, data *Daily, id string) (err error) {
	db.Where("id = ?", id).Delete(data)
	return nil
}

func AddDailyTotal(db *gorm.DB, data *Daily, id string) (err error) {
	db.Debug().Model(data).Where("id = ?", id).UpdateColumn("total", gorm.Expr("total + ?", 1))
	return nil
}

func DeductDailyTotal(db *gorm.DB, data *Daily, id string) (err error) {
	db.Debug().Model(data).Where("id = ?", id).UpdateColumn("total", gorm.Expr("total - ?", 1))
	return nil
}
