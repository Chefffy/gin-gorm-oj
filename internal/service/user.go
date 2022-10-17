package service

import (
	"gin-gorm-oj/define"
	"gin-gorm-oj/helper"
	"gin-gorm-oj/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"time"
)

// GetUserDetail
// @Tags 公共方法
// @Summary 用户详情
// @Param identity query string false "problem_identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /user-detail [get]
func GetUserDetail(c *gin.Context){
	identity := c.Query("identity")
	if identity == ""{
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"用户唯一标识不能为空",
		})
		return
	}
	data:=new(models.UserBasic)
	err := models.DB.Omit("password").Where("identity = ? ",identity).Find(&data).Error
	if err !=nil{
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"Get User Detail By Identity: "+identity +"Error: "+err.Error(),
		})
	}
	c.JSON(http.StatusOK,gin.H{
		"code":200,
		"data":data,
	})
}

// Login
// @Tags 公共方法
// @Summary 用户登录
// @Param username formData string false "username"
// @Param password formData string false "password"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /login [post]
func Login(c *gin.Context){
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username ==""||password == ""{
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"信息为空",
		})
	}

	//md5
	password = helper.GetMd5(password)

	data :=new(models.UserBasic)
	err := models.DB.Where("name = ? AND password = ?",username,password).Find(&data).Error
	if err != nil{
		if err == gorm.ErrRecordNotFound{
			c.JSON(http.StatusOK,gin.H{
				"code":-1,
				"msg":"用户名或密码错误",
			})
			return
		}
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"Get userBasic Error"+err.Error(),
		})
		return
	}

	token,err := helper.GenerateToken(data.Identity,data.Name,data.IsAdmin)
	if err !=nil {
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"Generation Error: "+err.Error(),
		})
	}

	c.JSON(http.StatusOK,gin.H{
		"code":200,
		"data":map[string]interface{}{
			"token":token,
		},
	})
}

// SendCode
// @Tags 公共方法
// @Summary 发送验证码
// @Param email formData string true "email"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /send-code [post]
func SendCode(c *gin.Context){
	email := c.PostForm("email")
	if email ==""{
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"email empty",
		})
		return
	}
	code := helper.GetRand()

	models.RDB.Set(c,email,code,time.Second*300)

	err := helper.SendCode(email,code)
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg": "Send Code Error: "+err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"code":200,
		"msg":"successfully sent",
	})
}


// Register
// @Tags 公共方法
// @Summary 用户注册
// @Param mail formData string true "mail"
// @Param code formData string true "code"
// @Param name formData string true "name"
// @Param password formData string true "password"
// @Param phone formData string false "phone"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /register [post]
func Register(c *gin.Context){
	mail := c.PostForm("mail")
	if mail == ""{
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"mail empty",
		})
		return
	}
	userCode := c.PostForm("code")
	if userCode == ""{
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"code empty",
		})
		return
	}
	name := c.PostForm("name")
	if name == ""{
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"name empty",
		})
		return
	}
	password := c.PostForm("password")
	if password == ""{
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"password empty",
		})
		return
	}
	phone := c.PostForm("phone")

	//验证码是否正确
	sysCode, err := models.RDB.Get(c,mail).Result()
	if err != nil{
		log.Printf("Get Code Error: %v \n",err)
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"Get Code Error "+ err.Error(),
		})
		return
	}

	if sysCode != userCode {
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"Wrong Verification Code "+ err.Error(),
		})
		return
	}

	//判断邮箱是否已存在
	var cnt int64
	err = models.DB.Where("mail = ?",mail).Model(new(models.UserBasic)).
		Count(&cnt).Error
	if err != nil{
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"Get User Error: "+err.Error(),
		})
		return
	}

	if cnt>0 {
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"Email has been used ",
		})
		return
	}

	//数据的插入
	userIdentity := helper.GetUUid()
	data := &models.UserBasic{
		Identity: userIdentity,
		Name:     name,
		Password: helper.GetMd5(password),
		Phone:    phone,
		Mail:     mail,
	}
	err = models.DB.Create(data).Error
	if err != nil{
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"Create User Error: "+err.Error(),
		})
		return
	}

	//生成token
	token,err := helper.GenerateToken(userIdentity, name,data.IsAdmin)
	if err != nil{
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"Generate Token Error: "+err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"code":200,
		"data":map[string]interface{}{
			"token":token,
		},
	})
}

// GetRankList
// @Tags 公共方法
// @Summary 用户排行榜
// @Param page query int false "请输入当前页，默认第一页"
// @Param size query int false "size"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /rank-list [get]
func GetRankList(c *gin.Context){
	size,_:= strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))

	page,err:= strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err!=nil{
		log.Println("GetProblemList Page strconv Parse Error:",err)
		return
	}

	page =(page -1)*size

	var count int64
	list := make([]*models.UserBasic,0)
	err = models.DB.Model(new(models.UserBasic)).Count(&count).Order("pass_num DESC,submit_num ASC").
		Offset(page).Limit(size).Find(&list).Error
	if err != nil{
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"Get Rank List Error: "+err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"code":200,
		"data":map[string]interface{}{
			"List":list,
			"count":count,
		},
	})
}