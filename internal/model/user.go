package model

import (
	"errors"
	"gorm.io/gorm"
	"shoppingApp/internal/client"
	"shoppingApp/internal/data"
	"shoppingApp/internal/validator"
	"strings"
	"time"
)

type User struct {
	ID        int64          `json:"id" gorm:"primarykey"` // Unique integer ID for the product
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	FullName  string         `json:"full_name" gorm:"type:text;size:100;not null"`
	Email     string         `json:"email" gorm:"type:text;size:100;unique;not null"`
	Password  []byte         `json:"-" gorm:"type:bytea;size:100;not null"`
	Activated bool           `json:"activated" gorm:"type:bool;not null"`
}

// Custom ErrDuplicateEmail error.
var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

func (user User) Create() (User, error) {
	err := client.PostgreSqlClient.Create(&user)
	if err.Error != nil {
		switch {
		case strings.Contains(err.Error.Error(), "duplicate key value violates unique constraint \"users_email_key\""):
			return user, ErrDuplicateEmail
		default:
			return user, err.Error
		}
	}

	return user, nil
}

func (user User) Find(query ...interface{}) (User, error) {
	err := client.PostgreSqlClient.First(&user, query...)
	return user, err.Error
}

func (user User) Get(filter data.Filters, query ...interface{}) ([]User, data.Metadata) {
	var users []User
	postClient := client.PostgreSqlClient
	var totalCount int64

	postClient = postClient.Order(filter.SortColumn() + " " + filter.SortDirection())
	postClient = postClient.Limit(filter.Limit())
	postClient = postClient.Offset(filter.Offset())

	postClient.Find(&users, query...)
	postClient.Model(&User{}).Count(&totalCount)

	metadata := data.CalculateMetadata(totalCount, filter.Page, filter.PageSize)
	return users, metadata
}

func (user User) Update(column string, value interface{}) error {
	err := client.PostgreSqlClient.Model(&user).Update(column, value)
	return err.Error
}

func (user User) Updates(data User) error {
	err := client.PostgreSqlClient.Model(&user).Updates(data)
	return err.Error
}

func (user User) Delete(column string, value interface{}) error {
	err := client.PostgreSqlClient.Model(&user).Delete(column, value)
	return err.Error
}

func (user User) Count(column string, value interface{}) int64 {
	var counter int64
	client.PostgreSqlClient.Model(&user).Where(column, value).Count(&counter)
	return counter
}

// ValidateEmail validates the email address.
func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

// ValidatePasswordPlaintext validates the password.
func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 character long")
	v.Check(len(password) <= 30, "password", "must not be more than 30 character long")
}

// ValidateUser validates the user data.
func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.FullName != "", "full_name", "must be provided")
	v.Check(len(user.FullName) <= 100, "name", "must not be more than 100 bytes long")
	ValidateEmail(v, user.Email)
	ValidatePasswordPlaintext(v, string(user.Password))
}
