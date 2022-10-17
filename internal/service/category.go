package service

import (
	"gin-gorm-oj/define"
	"gin-gorm-oj/helper"
	"gin-gorm-oj/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// GetCategoryList
// @Tags 管理员私有方法
// @Summary 分类列表
// @Param authorization header string true "authorization token"
// @Param page query int false "page"
// @Param size query int false "size"
// @Param keyword query string false "keyword"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/category-list [get]
func GetCategoryList(c *gin.Context){
	size,_:= strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))

	page,err:= strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err!=nil{
		log.Println("GetProblemList Page strconv Parse Error:",err)
		return
	}

	page =(page -1)*size
	var count int64
	keyword :=c.Query("keyword")

	categoryList := make([]*models.CategoryBasic,0)
	err = models.DB.Model(new(models.CategoryBasic)).Where("name like ?","%"+keyword+"%").
		Count(&count).Limit(size).Offset(page).Find(&categoryList).Error
	if err !=nil{
		log.Println("GetCategoryList Error: ",err)
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"Get Category List Wrong",
		})
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"code":200,
		"data":map[string]interface{}{
			"list": categoryList,
			"count": count,
		},
	})
}

// CategoryCreate
// @Tags 管理员私有方法
// @Summary 分类创建
// @Param authorization header string true "authorization"
// @Param name formData string true "name"
// @Param parentId formData int false "parentId"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/category-create [post]
func CategoryCreate(c *gin.Context){
	name := c.PostForm("name")
	parentId,_ := strconv.Atoi(c.PostForm("parentId"))
	category := &models.CategoryBasic{
		Identity: helper.GetUUid(),
		Name:     name,
		ParentId: parentId,
	}
	err := models.DB.Create(category).Error
	if err != nil{
		log.Println("CategoryCreate Error: ",err)
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"CategoryCreate Error",
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"code":200,
		"msg":"Create success",
	})
}

// CategoryDelete
// @Tags 管理员私有方法
// @Summary 分类修改
// @Param authorization header string true "authorization"
// @Param identity query string true "identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/category-delete [delete]
func CategoryDelete(c *gin.Context){
	identity := c.Query("identity")
	if identity == ""{
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"Parameter identity empty",
		})
		return
	}

	var count int64
	err := models.DB.Model(new(models.ProblemCategory)).
		Where("category_id = (SELECT id FROM category_basic WHERE identity = ? LIMIT 1)",identity).Count(&count).Error
	if err != nil{
		log.Println("GetProblemCategory Error: ",err)
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"Get Problem Category Failed",
		})
		return
	}
	if count > 0{
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"Problem exist,can not delete",
		})
		return
	}
	err = models.DB.Where("identity = ?",identity).Delete(new(models.CategoryBasic)).Error
	if err != nil{
		log.Println("Delete Category Error: ",err)
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"Delete Category Error",
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"code":200,
		"msg":"Successful Delete",
	})
}

// CategoryModify
// @Tags 管理员私有方法
// @Summary 分类修改
// @Param authorization header string true "authorization"
// @Param identity formData string true "identity"
// @Param name formData string true "name"
// @Param parentId formData int false "parentId"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/category-modify [put]
func CategoryModify(c *gin.Context){
	identity := c.PostForm("identity")
	name := c.PostForm("name")
	parentId,_ := strconv.Atoi(c.PostForm("parentId"))
	if name == ""|| identity ==""{
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"Parameter Empty",
		})
		return
	}
	category := &models.CategoryBasic{
		Identity: identity,
		Name: name,
		ParentId: parentId,
	}
	err := models.DB.Model(new(models.CategoryBasic)).Where("identity = ?",identity).Updates(category).Error
	if err != nil{
		log.Println("CategoryModify Error: ",err)
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"CategoryModify Error",
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"code":200,
		"msg":"Modify success",
	})
}