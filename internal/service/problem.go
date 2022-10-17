package service

import (
	"encoding/json"
	"errors"
	"gin-gorm-oj/define"
	"gin-gorm-oj/helper"
	"gin-gorm-oj/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

// GetProblemList
// @Tags 公共方法
// @Summary 问题列表
// @Param page query int false "page"
// @Param size query int false "size"
// @Param keyword query string false "keyword"
// @Param category_identity query string false "category_identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /problem-list [get]
func GetProblemList(c *gin.Context){
	size,_:= strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))

	page,err:= strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err!=nil{
		log.Println("GetProblemList Page strconv Parse Error:",err)
		return
	}

	page =(page -1)*size
	var count int64
	keyword :=c.Query("keyword")
	categoryIdentity :=c.Query("category_identity")

	list := make([] *models.ProblemBasic,0)
	tx := models.GetProblemList(keyword,categoryIdentity)
	err = tx.Count(&count).Omit("content").Offset(page).Limit(size).Find(&list).Error
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"Get Problem List Error:"+err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"code":200,
		"data":map[string]interface{}{
			"list":list,
			"count":count,
		},
	})
}

// GetProblemDetail
// @Tags 公共方法
// @Summary 问题详情
// @Param identity query string false "problem_identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /problem-detail [get]
func GetProblemDetail(c *gin.Context){
	identity :=c.Query("identity")
	if identity == ""{
		c.JSON(http.StatusOK,gin.H{
			"code": -1,
			"msg": "问题唯一标识不能为空",
		})
		return
	}
	data := new(models.ProblemBasic)
	err := models.DB.Where("identity = ?",identity).
		Preload("ProblemCategories").Preload("ProblemCategories.CategoryBasic").
		First(&data).Error
	if err != nil {
		if err ==gorm.ErrRecordNotFound{
			c.JSON(http.StatusOK,gin.H{
				"code": -1,
				"msg": "问题不存在",
			})
			return
		}
		c.JSON(http.StatusOK,gin.H{
			"code": -1,
			"msg": "Get Problem Detail Error :"+err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"code": 200,
		"data": data,
	})
}

// ProblemCreate
// @Tags 管理员私有方法
// @Summary 问题创建
// @Param authorization header string true "authorization"
// @Param title formData string true "title"
// @Param content formData string true "content"
// @Param max_runtime formData int false "max_runtime"
// @Param max_mem formData int false "max_mem"
// @Param category_ids formData []string false "category_id" collectionFormat(multi)
// @Param test_cases formData []string true "test_cases" collectionFormat(multi)
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/problem-create [post]
func ProblemCreate(c *gin.Context){
	title := c.PostForm("title")
	content := c.PostForm("content")
	maxRuntime,_ := strconv.Atoi(c.PostForm("max_runtime"))
	maxMem,_ := strconv.Atoi(c.PostForm("max_mem"))
	categoryIds := c.PostFormArray("category_ids")
	testCases := c.PostFormArray("test_cases")
	if title == ""|| content == ""|| len(categoryIds) ==0 || len(testCases)==0 || maxRuntime==0 || maxMem==0 {
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"Parameter can not be empty",
		})
		return
	}

	identity := helper.GetUUid()
	data := &models.ProblemBasic{
		Identity: identity,
		Title: title,
		Content: content,
		MaxRuntime: maxRuntime,
		MaxMem: maxMem,
	}

	//处理分类
	categoryBasics :=make([]*models.ProblemCategory,0)
	for _,id :=range categoryIds{
		categoryID,_ := strconv.Atoi(id)
		categoryBasics =append(categoryBasics,&models.ProblemCategory{
			ProblemId: data.ID,
			CategoryId: uint(categoryID),
		})
	}
	data.ProblemCategories = categoryBasics

	//处理测试用例
	testCaseBasics := make([]*models.TestCase,0)
	for _,testCase :=range testCases{
		caseMap :=make(map[string]string)
		err := json.Unmarshal([]byte(testCase),&caseMap)
		if err != nil{
			c.JSON(http.StatusOK,gin.H{
				"code":-1,
				"msg":"test format error",
			})
			return
		}
		if _,ok := caseMap["input"];!ok{
			c.JSON(http.StatusOK,gin.H{
				"code":-1,
				"msg":"test format error",
			})
			return
		}
		if _,ok := caseMap["output"];!ok{
			c.JSON(http.StatusOK,gin.H{
				"code":-1,
				"msg":"test format error",
			})
			return
		}
		testCaseBasic := &models.TestCase{
			Identity:        helper.GetUUid(),
			ProblemIdentity: identity,
			Input:           caseMap["input"],
			Output:          caseMap["output"],
		}
		testCaseBasics = append(testCaseBasics,testCaseBasic)
	}
	data.TestCases = testCaseBasics

	//创建问题
	err := models.DB.Create(data).Error
	if err != nil{
		c.JSON(http.StatusOK,gin.H{
			"code":200,
			"data":map[string]interface{}{
				"identity":data.Identity,
			},
		})
	}
}

// ProblemModify
// @Tags 管理员私有方法
// @Summary 问题修改
// @Param authorization header string true "authorization"
// @Param identity formData string true "identity"
// @Param title formData string true "title"
// @Param content formData string true "content"
// @Param max_runtime formData int false "max_runtime"
// @Param max_mem formData int false "max_mem"
// @Param category_ids formData []string false "category_id" collectionFormat(multi)
// @Param test_cases formData []string true "test_cases" collectionFormat(multi)
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/problem-modify [put]
func ProblemModify(c *gin.Context){
	identity := c.PostForm("identity")
	title := c.PostForm("title")
	content := c.PostForm("content")
	maxRuntime,_ := strconv.Atoi(c.PostForm("max_runtime"))
	maxMem,_ := strconv.Atoi(c.PostForm("max_mem"))
	categoryIds := c.PostFormArray("category_ids")
	testCases := c.PostFormArray("test_cases")
	if identity== "" || title == ""|| content == ""|| len(categoryIds) ==0 || len(testCases)==0 || maxRuntime==0 || maxMem==0 {
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"Parameter can not be empty",
		})
		return
	}

	if err := models.DB.Transaction(func (tx *gorm.DB) error{
		//问题基础信息保存
		problemBasic := models.ProblemBasic{
			Identity: identity,
			Title: title,
			Content: content,
			MaxRuntime: maxRuntime,
			MaxMem: maxMem,
		}
		err := tx.Where("identity = ?",identity).Updates(problemBasic).Error
		if err != nil{
			return err
		}

		//查询问题详情
		err = tx.Where("identity = ?",identity).Find(problemBasic).Error
		if err != nil{
			return err
		}

		//关联问题分类更新
		//删除已存在的关联关系
		err = tx.Where("problem_id = ? ",problemBasic.ID).Delete(new(models.ProblemCategory)).Error
		if err != nil{
			return err
		}
		//新增新的关联关系
		pcs := make([]*models.ProblemCategory,0)
		for _,id := range categoryIds{
			intId, _ := strconv.Atoi(id)
			pcs =append(pcs,&models.ProblemCategory{
				ProblemId: problemBasic.ID,
				CategoryId: uint(intId),
			})
		}
		err = tx.Create(&pcs).Error
		if err != nil{
			return err
		}

		//关联测试案例的更新
		//删除已存在的关联关系
		err = tx.Where("problem_identity = ? ",identity).Delete(new(models.TestCase)).Error
		if err != nil{
			return err
		}
		//新增新的关联关系
		tcs := make([]*models.TestCase,0)
		for _,testCases := range testCases{
			caseMap := make(map[string]string)
			err := json.Unmarshal([]byte(testCases),&caseMap)
			if err != nil{
				return err
			}
			if _,ok := caseMap["input"]; !ok{
				return errors.New("testcase input wrong")
			}
			if _,ok := caseMap["output"]; !ok{
				return errors.New("testcase output wrong")
			}
			tcs = append(tcs,&models.TestCase{
				Identity:        helper.GetUUid(),
				ProblemIdentity: identity,
				Input:           caseMap["input"],
				Output:          caseMap["output"],
			})
		}
		err = tx.Create(tcs).Error
		if err != nil{
			return err
		}
		return nil
	});err != nil{
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg": "Problem Modify Error: "+err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"code":200,
		"msg":"Problem Modify Success",
	})

}