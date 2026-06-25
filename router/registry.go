// 调用service方法
// 聚合所有的 Handler，并准备把它们“注册（Register）”到 Gin 的路由引擎上。
package router

import (
	"errors"
	"net/http"
	"s9-tss-gin/model"
	"s9-tss-gin/services"

	"github.com/ethereum/go-ethereum/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gin-gonic/gin"
)

type Registry struct {
	signService services.SignService
}

// 注册到Gin上，依赖注入
func NewRegistry(signService services.SignService) *Registry {
	return &Registry{
		signService: signService,
	}
}

// "定义了一个名叫 KeygenHandler 的方法，它专门生产 Gin 所需的处理函数。"
// 当你执行 registry.KeygenHandler() 时，你拿到的是一个真正的 func(c *gin.Context)。
// Registry.signService.接口里面方法
// 命名 接口方法+Handler
func (registry Registry) KeygenHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		//把 HTTP 请求体里的 JSON 字符串，转成 Go 语言的内存结构体
		// 反序列化
		var request model.KeygenRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
			return
		}

		cpk, err := registry.signService.Keygen(request)
		if err != nil {
			c.String(http.StatusInternalServerError, "failed to sign state")
			log.Error("failed to sign state", "error", err)
			return
		}

		// 执行流程分为两个阶段：
		//
		//阶段一（服务启动时）：Gin 引擎在初始化时调用你的 SignTxHandler()，它立刻返回那个匿名函数。Gin 拿到这个匿名函数后，把它挂载到 /sign 路由上，
		//但不执行里面的代码。这时候它只是一段“待命逻辑”。
		//阶段二（请求到达时）：当用户真正发来 POST /sign 请求时，Gin 才会在后台执行你返回的那个匿名函数。
		if _, err = c.Writer.Write([]byte(cpk)); err != nil {
			c.String(http.StatusInternalServerError, "failed to sign state")
			log.Error("failed to sign state", "error", err)
		}

		// 也可以返回 JSON
		// 方法一：返回二进制（原样）
		// c.Writer.Write(signature)
		// 方法二：返回 JSON（如果你想改的话，完全合法！）
		// c.JSON(http.StatusOK, gin.H{
		// 	"signature": string(signature), // 转成字符串返回
		// })
	}
}

func (registry *Registry) SignTxHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request model.TransactionSignRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, errors.New("invalid request body"))
			return
		}
		if request.MessageHash == "" {
			c.JSON(http.StatusBadRequest, errors.New("StartBlock and OffsetStartsAtIndex must not be nil or negative"))
			return
		}
		signature, err := registry.signService.TransactionSign(request)
		if err != nil {
			c.String(http.StatusInternalServerError, "failed to sign state")
			log.Error("failed to sign state", "error", err)
			return
		}
		if _, err = c.Writer.Write(signature); err != nil {
			log.Error("failed to write signature to response writer", "error", err)
		}
	}
}

func (registry *Registry) PrometheusHandler() gin.HandlerFunc {
	h := promhttp.InstrumentMetricHandler(
		prometheus.DefaultRegisterer, promhttp.HandlerFor(
			prometheus.DefaultGatherer,
			promhttp.HandlerOpts{MaxRequestsInFlight: 3},
		),
	)

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
