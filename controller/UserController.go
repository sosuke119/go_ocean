package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"go_Ocean/common"
	"go_Ocean/model"
)

func Register(c *gin.Context) {
	DB := common.GetDB()

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
	if isTelephoneExist(DB, telephone) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手機號碼已經存在"})
		return
	}

	// 創建用戶
	newUser := model.User{
		Name:      name,
		Password:  password,
		Telephone: telephone,
	}
	DB.Create(&newUser)

	// 返回結果

	c.JSON(200, gin.H{
		"msg": "註冊成功",
	})
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		// 如果用戶存在，回傳true
		return true
	}
	return false
}
