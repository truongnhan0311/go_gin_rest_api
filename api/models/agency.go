package models

import (
	"gorm.io/gorm"
	"time"
)

type Agency struct {
	gorm.Model
	ID        uint   `json:"id" gorm:"primary_key"`
	Name      string `json:"name"`
	MaSoDL    string `json:"ma_so_dl"`
	Address   string `json:"address"`
	Total     int    `json:"total" gorm:"default:0"`
	Status    string `json:"status"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Deleted   gorm.DeletedAt
}

func Create(db *gorm.DB, data *Agency) (err error) {
	err = db.Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func GetList(db *gorm.DB, data *[]Agency, agency string) (err error) {

	if agency == "" {
		err = db.Find(data).Error
	} else {
		err = db.Where("id = ?", agency).Find(data).Error
	}

	if err != nil {
		return err
	}
	return nil
}

func Get(db *gorm.DB, data *Agency, id string) (err error) {
	err = db.Where("id = ?", id).First(data).Error
	if err != nil {
		return err
	}
	return nil
}

//update user
func Update(db *gorm.DB, data *Agency) (err error) {
	db.Save(data)
	return nil
}

//delete user
func Delete(db *gorm.DB, data *Agency, id string) (err error) {
	db.Where("id = ?", id).Delete(data)
	return nil
}

func AddTotal(db *gorm.DB, data *Agency, id string) (err error) {
	db.Debug().Model(data).Where("id = ?", id).UpdateColumn("total", gorm.Expr("total + ?", 1))
	return nil
}

func DeductTotal(db *gorm.DB, data *Agency, id string) (err error) {
	db.Debug().Model(data).Where("id = ?", id).UpdateColumn("total", gorm.Expr("total - ?", 1))
	return nil
}
