package service

import (
	"bytes"
	"errors"
	"gin-gorm-oj/define"
	"gin-gorm-oj/helper"
	"gin-gorm-oj/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// GetSubmitList
// @Tags 公共方法
// @Summary 提交列表
// @Param page query int false "page"
// @Param size query int false "size"
// @Param status query int false "status"
// @Param problem_identity query string false "problem_identity"
// @Param user_identity query string false "user_identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /submit-list [get]
func GetSubmitList(c *gin.Context){
	size,_:= strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))

	page,err:= strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err!=nil{
		log.Println("GetProblemList Page strconv Parse Error:",err)
		return
	}

	page =(page -1)*size
	var count int64
	list := make([]models.SubmitBasic,0)


	problemIdentity := c.Query("problem_identity")
	userIdentity := c.Query("user_identity")
	status, _ :=strconv.Atoi(c.Query("status"))

	tx := models.GetSubmitList(problemIdentity,userIdentity,status)
	err = tx.Count(&count).Offset(page).Limit(size).Find(&list).Error
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"Get Submit List Error: "+err.Error(),
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

// Submit
// @Tags 用户私有方法
// @Summary 代码提交
// @Param authorization header string true "authorization"
// @Param problem_identity query string true "problem_identity"
// @Param code body string true "code"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /user/submit [post]
func Submit(c *gin.Context){
	problemIdentity := c.Query("problem_identity")
	code, err := ioutil.ReadAll(c.Request.Body)
	if err != nil{
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg": "Read Code Error: "+err.Error(),
		})
		return
	}

	//代码保存
	path, err := helper.CodeSave(code)
	if err !=nil {
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"Code Save Error: "+err.Error(),
		})
		return
	}

	u,_ := c.Get("user")
	userClaim := u.(*helper.UserClaims)
	sb:= &models.SubmitBasic{
		Identity:        helper.GetUUid(),
		ProblemIdentity: problemIdentity,
		UserIdentity:    userClaim.Identity,
		Path:            path,
	}

	// 代码判断
	pb := new(models.ProblemBasic)
	err = models.DB.Where("identity = ?",problemIdentity).Preload("TestCases").First(pb).Error
	if err != nil{
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"Get Problem Error: "+err.Error(),
		})
		return
	}

	//答案错误
	WA := make(chan int)
	//超内存
	OOM := make(chan int)
	//编译错误
	CE := make(chan int)
	//通过个数
	passCount := 0
	var lock sync.Mutex
	//提示信息
	var msg string

	for _, testCase:=range pb.TestCases{
		testCase :=testCase
		go func() {
			//执行测试
			cmd := exec.Command("go","run",path)
			var out,stderr bytes.Buffer
			cmd.Stderr = &stderr
			cmd.Stdout = &out
			stdinPipe, err := cmd.StdinPipe()
			if err != nil {
				log.Fatalln(err)
			}
			io.WriteString(stdinPipe,testCase.Input)

			var bm runtime.MemStats
			runtime.ReadMemStats(&bm)
			//根据测试的输入案例进行运行，拿到输出结果和标准的输出结果进行比对
			if err := cmd.Run();err != nil{
				log.Println(err,stderr.String())
				if err.Error() == "exit status 2"{
					msg = stderr.String()
					CE <-1
					return
				}
			}

			var em runtime.MemStats
			runtime.ReadMemStats(&em)
			//答案错误
			if testCase.Output != out.String(){
				msg ="WA"
				WA <- 1
				return
			}
			//运行超内存
			if (em.Alloc /1024 - bm.Alloc/1024) > uint64(pb.MaxMem){
				msg = "OOM"
				OOM <- 1
				return
			}

			lock.Lock()

			passCount ++

			lock.Unlock()
		}()
	}
	select {
	//-1-待判断，1-答案正确，2-答案错误，3-运行超时，4-运行超内存，5-编译错误
	case <-WA:
		sb.Status = 2
	case <-OOM:
		sb.Status = 4
	case <-time.After(time.Millisecond * time.Duration(pb.MaxRuntime)):
		if passCount == len(pb.TestCases){
			sb.Status = 1
		}else {
			sb.Status =3
		}
	case <-CE:
		sb.Status = 5
	}



	if err = models.DB.Transaction(func(tx *gorm.DB)error{
		err = tx.Create(sb).Error
		if err != nil{
			return errors.New("SubmitBasic Save Error: "+err.Error())
		}
		m := make(map[string]interface{})
		m["submit_num"]= gorm.Expr("submit_num + ?",1)
		if sb.Status == 1{
			m["pass_num"]= gorm.Expr("pass_num + ?",1)
		}
		// 更新user_basic
		err =tx.Model(new(models.UserBasic)).Where("identity = ? ",userClaim.Identity).Updates(m).Error
		if err != nil{
			return errors.New("UserBasic Modify Error: "+err.Error())
		}
		// 更新problem_basic
		err =tx.Model(new(models.ProblemBasic)).Where("identity = ? ",problemIdentity).Updates(m).Error
		if err != nil{
			return errors.New("ProblemBasic Modify Error: "+err.Error())
		}
		return nil
	}); err != nil{
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"Submit Error: "+err.Error(),
		})
		return
	}

	if err != nil{
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"Submit Error: "+err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"code":200,
		"data": map[string]interface{}{
			"status":sb.Status,
			"msg":msg,
		},
	})
}
