package controller

import (
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"go_Ocean/common"
	"go_Ocean/dto"
	"go_Ocean/model"
	"go_Ocean/response"
)

func Register(c *gin.Context) {
	DB := common.GetDB()

	// 使用map獲取參數
	var requestUser = model.User{}
	c.Bind(&requestUser)

	// 獲取參數
	name := requestUser.Name
	password := requestUser.Password
	telephone := requestUser.Telephone

	// 數據驗證
	if len(name) == 0 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "名稱不能為空")
		return
	}
	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "密碼不能少於六位")
		return
	}
	if len(telephone) != 10 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "手號碼必須為10位")
		return
	}

	log.Println(name, password, telephone)

	// 判斷手機號碼是否存在
	if isTelephoneExist(DB, telephone) {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "手機號碼已經存在")
		return
	}

	// 創建用戶
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity, 500, nil, "加密錯誤")
		return
	}
	newUser := model.User{
		Name:      name,
		Password:  string(hashedPassword),
		Telephone: telephone,
	}
	DB.Create(&newUser)

	// 發放Token
	token, err := common.ReleaseToken(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "系統異常"})
		log.Printf("token generate error: %v", err)
		return
	}

	// 返回結果
	response.Success(c, gin.H{"token": token}, "註冊成功")
}

func Login(c *gin.Context) {
	DB := common.GetDB()

	// 獲取參數
	password := c.PostForm("password")
	telephone := c.PostForm("telephone")

	// 數據驗證
	if len(password) < 6 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密碼不能少於六位"})
		return
	}
	if len(telephone) != 10 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手號碼必須為10位"})
		return
	}

	// 判斷手機號是否存在
	var user model.User
	DB.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用戶不存在"})
		return
	}

	// 判斷密碼是否正確
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "密碼錯誤"})
		return
	}

	// 發放Token
	token, err := common.ReleaseToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "系統異常"})
		log.Printf("token generate error: %v", err)
		return
	}

	// 返回結果
	response.Success(c, gin.H{"token": token}, "登入成功")
}

func Info(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": dto.ToUserDto(user.(model.User))}})
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
