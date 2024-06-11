package routers

import (
	v1 "fisco-sgx-go/internal/routers/api/v1"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	fisco := v1.NewFisco()
	quote := v1.NewQuote()
	report := v1.NewReport()

	apiv1 := r.Group("/api/v1")
	{
		//这个路由负责后台启动fisco程序
		apiv1.GET("/startSgx", fisco.Get)
		apiv1.DELETE("/stopSgx", fisco.DELETE)

		//这个路由负责远程证明
		apiv1.GET("/quote", quote.Get)

		//这个路由负责获取报告
		apiv1.GET("/report", report.Get)
	}

	return r
}
