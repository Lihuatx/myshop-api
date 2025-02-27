package middlewares

import (
	"fmt"
	"myshop_api/goods-web/global"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

// Trace 是一个Gin框架的中间件函数，用于启用Jaeger分布式追踪。
// 它返回一个Gin的HandlerFunc类型的函数，该函数将Jaeger追踪集成到Gin的HTTP请求处理流程中。
func Trace() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cfg := jaegercfg.Configuration{
			Sampler: &jaegercfg.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: 1,
			},
			Reporter: &jaegercfg.ReporterConfig{
				LogSpans: true,
				//CollectorEndpoint: "http://172.26.240.124:14268/api/traces",
				CollectorEndpoint: fmt.Sprintf("http://%s:%d/api/traces", global.ServerConfig.JaegerInfo.Host, global.ServerConfig.JaegerInfo.Port),
				//LocalAgentHostPort: fmt.Sprintf("%s:%d", global.ServerConfig.JaegerInfo.Host, global.ServerConfig.JaegerInfo.Port),
			},
			ServiceName: global.ServerConfig.JaegerInfo.Name,
		}

		tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
		if err != nil {
			panic(err)
		}
		opentracing.SetGlobalTracer(tracer)
		defer closer.Close()

		startSpan := tracer.StartSpan(ctx.Request.URL.Path)
		defer startSpan.Finish()

		ctx.Set("tracer", tracer)
		ctx.Set("parentSpan", startSpan)
		ctx.Next()
	}
}
