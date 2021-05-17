package models

import (
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User model
type User struct {
	gorm.Model
	Email        string `gorm:"type:varchar(100);unique_index" json:email`
	Fname        string `gorm:"size:100;not null" json:fname`
	Lname        string `gorm:"size:100;not null" json:lname`
	Password     string `gorm:"size:100;not null" json:password`
	ProfileImage string `gorm:"size:255" json:profile_image`
}

// HashPassword hashes password from user input

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) // 14 is the cost for hashing the password
	return string(bytes), err
}

// CheckPasswordHash checks password hash and password from user input if they match
func CheckPassword(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(hash))
	if err != nil {
		return errors.New("password incorrect")
	}
	return nil
}

// BeforeSave hashes user password
func (u *User) BeforeSave() error {
	password := strings.TrimSpace(u.Password)
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword
	return nil
}

// Prepare strips user inputs of any white spaces
func (u *User) Prepare() {
	u.Email = strings.TrimSpace(u.Email)
	u.Fname = strings.TrimSpace(u.Fname)
	u.Lname = strings.TrimSpace(u.Lname)
	u.ProfileImage = strings.TrimSpace(u.ProfileImage)
}

// Validate user input
func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "login":
		if u.Email == "" {
			return errors.New("Email is required")
		}
		if u.Password == "" {
			return errors.New("Password is required")
		}
	// Case for creating a User
	default:
		if u.Fname == "" {
			return errors.New("First Name is required")
		}
		if u.Lname == "" {
			return errors.New("Last Name is required")
		}
		if u.Email == "" {
			return errors.New("Email is required")
		}
		if u.Password == "" {
			return errors.New("Password is required")
		}
		return nil
	}
	return nil
}

// SaveUser adds a user to the database
func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	var err error

	// Debug a single operation, show detailed log for this operation
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

// GetUser returns a user based on email
func (u *User) GetUser(db *gorm.DB) (*User, error) {
	account := &User{}
	if err := db.Debug().Table("users").Where("email = ?", u.Email).First(account).Error; err != nil {
		return nil, err
	}
	return account, nil
}

// GetAllUsers returns a list of all the user
func GetAllUsers(db *gorm.DB) (*[]User, error) {
	users := []User{}
	if err := db.Debug().Table("users").Find(&users).Error; err != nil {
		return &[]User{}, err
	}
	return &users, nil
}
