package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"varchar(110);not null;unique"`
	Password  string `gorm:"size:255;not null"`
}

func main() {
	// Db
	db := InitDB()
	defer db.Close()

	// Gin
	r := gin.Default()
	r.POST("/api/auth/register", func(c *gin.Context) {
		// 獲取參數
		name := c.PostForm("name")
		password := c.PostForm("password")
		telephone := c.PostForm("telephone")

		// 數據驗證
		if len(name) == 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "名稱不能為空"})
			return
		}
		if len(password) < 6 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密碼不能少於六位"})
			return
		}
		if len(telephone) != 10 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手號碼必須為10位"})
			return
		}

		log.Println(name, password, telephone)

		// 判斷手機號碼是否存在
		if isTelephoneExist(db, telephone) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手機號碼已經存在"})
			return
		}

		// 創建用戶
		newUser := User{
			Name:      name,
			Password:  password,
			Telephone: telephone,
		}
		db.Create(&newUser)

		// 返回結果

		c.JSON(200, gin.H{
			"msg": "register success",
		})
	})

	panic(r.Run())
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		// 如果用戶存在，回傳true
		return true
	}
	return false
}

func InitDB() *gorm.DB {
	db, err := gorm.Open("mysql", "root:@root@(localhost:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database, err: " + err.Error())
	}
	db.AutoMigrate(&User{})
	return db
}
