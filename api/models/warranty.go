package models

import (
	"gorm.io/gorm"
	"strconv"
	"time"
)

type Warranty struct {
	gorm.Model
	ID        uint      `json:"id" gorm:"primary_key"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Soxe      string    `json:"soxe"`
	Tenxe     string    `json:"tenxe"`
	LoaiKinh  string    `json:"loai_kinh"`
	NgayDan   time.Time `json:"ngay_dan"`
	AgencyID  int       `json:"nha_phan_phoi"`
	DailyID   int       `json:"dai_ly"`
	UserID    int       `json:"nguoi_nhap"`
	MaSoBH    string    `json:"ma_so_bh"`
	DoiXe     string    `json:"doi_xe"`
	MauXe     string    `json:"mau_xe"`
	LoaiXe    string    `json:"loai_xe"`
	Status    string    `json:"status"`
	ViTriDan  string    `json:"vi_tri_dan"`
	MaSo      string    `json:"ma_so"`
	Email     string    `json:"email"`
	Note      string    `json:"note"`
	Agency    Agency
	Daily     Daily
	User      User
	CreatedAt time.Time
	UpdatedAt time.Time
	Deleted   gorm.DeletedAt
}

func CreateWarranty(db *gorm.DB, Warranty *Warranty) (err error) {
	err = db.Create(Warranty).Error
	if err != nil {
		return err
	}
	return nil
}

func GetWarrantyList(db *gorm.DB, Warranty *[]Warranty, agency string, daily string, phone string, soxe string, maso string, ngaydanfrom string, ngaydanto string) (err error) {
	result := db.Debug().Preload("Agency").Preload("Daily").Preload("User")
	if agency != "" {
		intVar, _ := strconv.Atoi(agency)
		result = result.Where("agency_id = ?", intVar)
	}
	if daily != "" {
		intVar, _ := strconv.Atoi(daily)
		result = result.Where("daily_id = ?", intVar)
	}
	if phone != "" {
		result = result.Where("phone LIKE ?", "%"+phone+"%")
	}
	if soxe != "" {
		result = result.Where("soxe LIKE ?", "%"+soxe+"%")
	}
	if maso != "" {
		result = result.Where("ma_so LIKE ?", "%"+maso+"%")
	}
	if ngaydanfrom != "" {
		result = result.Where("ngay_dan >= ?", ngaydanfrom+" 00:00:00")
	}
	if ngaydanto != "" {
		result = result.Where("ngay_dan <= ?", ngaydanto+" 23:00:00")
	}

	result = result.Where("status != ?", "deleted").Find(Warranty).Limit(1000)
	err = result.Error
	if err != nil {
		return err
	}
	return nil
}

func SearchWarranty(db *gorm.DB, Warranty *[]Warranty) (err error) {
	err = db.Preload("Agency").Preload("Daily").Preload("User").Where("status != ?", "deleted").Find(Warranty).Error
	if err != nil {
		return err
	}
	return nil
}

func GetWarranty(db *gorm.DB, Warranty *Warranty, id string, agency string) (err error) {

	if agency == "" {
		err = db.Preload("Agency").Preload("Daily").Preload("User").Where("id = ?", id).Where("status != ?", "deleted").First(Warranty).Error
	} else {
		err = db.Preload("Agency").Preload("Daily").Preload("User").Where("id = ?", id).Where("agency_id = ?", agency).Where("status != ?", "deleted").First(Warranty).Error
	}

	if err != nil {
		return err
	}
	return nil
}

func GetWarrantyByDaiLy(db *gorm.DB, Daily *Daily, id string, agency string) (err error) {

	if agency == "" {
		err = db.Preload("Daily").Where("id = ?", id).Where("status != ?", "deleted").First(Daily).Error
	} else {
		err = db.Preload("Daily").Where("id = ?", id).Where("daily = ?", agency).First(Daily).Where("status != ?", "deleted").Error
	}

	if err != nil {
		return err
	}
	return nil
}

func GetWarrantyByCode(db *gorm.DB, Warranty *Warranty, id string) (err error) {
	err = db.Preload("Agency").Preload("Daily").Preload("User").Where("ma_so_bh = ?", id).Where("status != ?", "deleted").First(Warranty).Error
	if err != nil {
		return err
	}
	return nil
}

//update user
func UpdateWarranty(db *gorm.DB, Warranty *Warranty) (err error) {
	db.Debug().Save(Warranty)
	return nil
}

//delete user
func DeleteWarranty(db *gorm.DB, Warranty *Warranty, id string) (err error) {
	db.Debug().Model(Warranty).Where("ma_so_bh = ?", id).UpdateColumn("status", "deleted")
	db.Where("ma_so_bh = ?", id).Delete(Warranty)
	return nil
}
