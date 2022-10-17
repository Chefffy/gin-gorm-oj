## 参考链接

GIN:https://gin-gonic.com/zh-cn/docs/

GORM:https://gorm.io/zh_CN/docs/



## 整合Swagger

参考文档：[https://github.com/swaggo/gin-swagger](https://gitee.com/link?target=https%3A%2F%2Fgithub.com%2Fswaggo%2Fgin-swagger)

接口访问链接：[http://localhost:8080/swagger/index.html](https://gitee.com/link?target=http%3A%2F%2Flocalhost%3A8080%2Fswagger%2Findex.html)

```
// GetProblemList
// @Tags 公共方法
// @Summary 问题列表
// @Param page query int false "page"
// @Param size query int false "size"
// @Success 200 {string} json "{"code":"200","msg","","data":""}"
// @Router /problem-list [get]
```



