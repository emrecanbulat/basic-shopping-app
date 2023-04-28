package client

import "gorm.io/gorm"

var PostgreSqlClient *gorm.DB

func Connections() {
	PostgreSqlClient = GetPostgreSqlClient()
}
