package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"go_Ocean/common"
)

func main() {
	// Db
	db := common.InitDB()
	defer db.Close()

	// Gin
	r := gin.Default()
	r = CollectRoute(r)
	panic(r.Run())
}
