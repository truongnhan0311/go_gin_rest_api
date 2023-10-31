package models

import (
	"gorm.io/gorm"
	"time"
)

type UserRepo struct {
	Db *gorm.DB
}

type User struct {
	gorm.Model
	ID        int
	Name      string
	Email     string
	Status    string
	Role      string
	Password  string
	AgencyID  *int
	Agency    Agency
	DailyID   *int
	Daily     Daily
	CreatedAt time.Time
	UpdatedAt time.Time
	Deleted   gorm.DeletedAt
}

//create a user
func CreateUser(db *gorm.DB, User *User) (err error) {
	err = db.Create(User).Error
	if err != nil {
		return err
	}
	return nil
}

//get users
func GetUsers(db *gorm.DB, User *[]User, agency string, daily string) (err error) {

	result := db.Debug().Preload("Agency").Preload("Daily")
	if agency != "" {
		result = result.Where("agency_id = ?", agency)
	}
	if daily != "" {
		result = result.Where("daily_id = ?", daily)
	}

	result = result.Find(User).Limit(1000)
	err = result.Error
	if err != nil {
		return err
	}
	return nil
}

//get user by id
func GetUser(db *gorm.DB, id string) (user *User, err error) {
	err = db.Preload("Agency").Preload("Daily").Where("id = ?", id).First(&user).Error
	return user, err
}

//update user
func UpdateUser(db *gorm.DB, User *User) (err error) {
	db.Save(User)
	return nil
}

//delete user
func DeleteUser(db *gorm.DB, User *User, id string) (err error) {
	db.Where("id = ?", id).Delete(User)
	return nil
}

//get user by id
func GetUserByEmail(db *gorm.DB, email string) (user *User, err error) {
	err = db.Preload("Agency").Preload("Daily").Where("email = ?", email).First(&user).Error
	return user, err
}
