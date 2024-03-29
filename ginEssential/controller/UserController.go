package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"oceanlearn.teach/ginessential/common"
	"oceanlearn.teach/ginessential/dto"
	"oceanlearn.teach/ginessential/model"
	"oceanlearn.teach/ginessential/response"
	"oceanlearn.teach/ginessential/util"
)

func Register(ctx *gin.Context) {
	DB := common.GetDB()
	//获取参数
	name := ctx.PostForm("name")
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")
	//数据验证
	if len(telephone) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"})
		return
	}
	log.Println(name, telephone, password)

	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码必须大于6位")

		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码必须大于6位"})
		return
	}
	log.Println(name, telephone, password)

	if len(name) == 0 {
		name = util.RandomString(10)
	}
	log.Println(name, telephone, password)

	if isTelephone(DB, telephone) {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户已经存在")

		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户已经存在"})
		return
	}

	//密码加密
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 500, nil, "加密错误")

		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 500, "msg": "加密错误"})
		return

	}
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hasedPassword),
	}
	DB.Create(&newUser)
	log.Println(name, telephone, password)
	//ctx.JSON(200, gin.H{
	//	"code":    200,
	//	"message": "success",
	//})
	response.Success(ctx, nil, "注册成功")
}

func Login(ctx *gin.Context) {
	//获取参数
	db := common.GetDB()
	telephone := ctx.PostForm("telephone")
	psd := ctx.PostForm("password")
	//	数据验证
	if len(telephone) != 11 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"})
		return
	}
	log.Println(telephone, psd)

	if len(psd) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码必须大于6位")

		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码必须大于6位"})
		return
	}

	//	判断手机号是否存在
	var user model.User
	db.Where("telephone=?", telephone).First(&user)

	log.Println(user)
	if user.ID == 0 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "用户不存在",
		})
	}

	//	判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(psd)); err != nil {
		response.Response(ctx, http.StatusBadRequest, 400, nil, "密码错误")

		//ctx.JSON(http.StatusBadRequest,gin.H{"code":400,"msg":"密码错误"})
	}

	//	发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(ctx, http.StatusBadRequest, 500, nil, "系统异常")

		//ctx.JSON(http.StatusInternalServerError,gin.H{"code":500,"msg":"系统异常"})
		log.Printf("token generate err: %v", err)
		return
	}
	//	返回结果

	//ctx.JSON(200,gin.H{
	//	"code":200,
	//	"data":gin.H{"token":token},
	//	"msg":"登录成功",
	//})
	response.Success(ctx, gin.H{"token": token}, "登录成功")

}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": dto.ToUserDto(user.(model.User))}})
}
func isTelephone(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone=?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
