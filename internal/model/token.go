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

func (token Token) Create(user User, jwtSecret string) (Token, error) {
	token, err := GenerateToken(user, jwtSecret)
	if err != nil {
		return token, err
	}

	result := client.PostgreSqlClient.Create(&token)
	return token, result.Error
}

func (token Token) Update(column string, value interface{}) error {
	result := client.PostgreSqlClient.Model(&token).Update(column, value)
	return result.Error
}

func (token Token) Updates(data Token) error {
	result := client.PostgreSqlClient.Model(&token).Updates(data)
	return result.Error
}

func (token Token) Find(query ...interface{}) (Token, error) {
	result := client.PostgreSqlClient.First(&token, query...)
	if token.UserID == 0 {
		return token, ErrRecordNotFound
	}
	return token, result.Error
}

func (token Token) Count(column string, value interface{}) int64 {
	var counter int64
	client.PostgreSqlClient.Model(&token).Where(column, value).Count(&counter)
	return counter
}

func GenerateToken(user User, jwtSecret string) (Token, error) {
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
		"id":    user.ID,
	}

	jwtBytes, _ := claims.HMACSign(jwt.HS256, []byte(jwtSecret))
	token.Hash = jwtBytes

	return token, nil
}
