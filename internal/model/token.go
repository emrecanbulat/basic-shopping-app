package model

import (
	"github.com/pascaldekloe/jwt"
	"shoppingApp/internal/client"
	"time"
)

// We did not use the REFERENCES user syntax because we use soft-delete

type Token struct {
	Hash   []byte    `json:"hash" gorm:"type:bytea;primarykey"`
	UserID int64     `json:"user_id"`
	Expiry time.Time `json:"expiry"`
}

func (token Token) Create(user User, role, jwtSecret string) (Token, error) {
	token, err := GenerateToken(user, role, jwtSecret)
	if err != nil {
		return token, err
	}

	result := client.PostgreSqlClient.Create(&token)
	return token, result.Error
}

func (token Token) Update(column string, value interface{}) error {
	err := client.PostgreSqlClient.Model(&token).Update(column, value)
	return err.Error
}

func (token Token) Updates(data Token) error {
	err := client.PostgreSqlClient.Model(&token).Updates(data)
	return err.Error
}

func (token Token) Find(query ...interface{}) (Token, error) {
	err := client.PostgreSqlClient.First(&token, query...)
	if token.UserID == 0 {
		return token, ErrRecordNotFound
	}
	return token, err.Error
}

func (token Token) Count(column string, value interface{}) int64 {
	var counter int64
	client.PostgreSqlClient.Model(&token).Where(column, value).Count(&counter)
	return counter
}

func GenerateToken(user User, role, jwtSecret string) (Token, error) {
	token := Token{
		UserID: user.ID,
		Expiry: time.Now().Add(24 * time.Hour),
	}

	var claims jwt.Claims
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	claims.Set = map[string]interface{}{
		"email": user.Email,
		"role":  role,
		"id":    user.ID,
	}

	jwtBytes, _ := claims.HMACSign(jwt.HS256, []byte(jwtSecret))
	token.Hash = jwtBytes

	return token, nil
}
