package tracer

import (
	"fmt"
	"io"

	"ocean-marketing/internal/config"
	"ocean-marketing/internal/pkg/logger"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
)

var tracer opentracing.Tracer
var closer io.Closer

// Init 初始化链路追踪
func Init(cfg config.TracerConfig) {
	jaegerCfg := jaegerConfig.Configuration{
		ServiceName: cfg.ServiceName,
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: cfg.SampleRate,
		},
		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: fmt.Sprintf("%s:%d", cfg.AgentHost, cfg.AgentPort),
		},
	}

	var err error
	tracer, closer, err = jaegerCfg.NewTracer()
	if err != nil {
		logger.Error("初始化链路追踪失败", zap.Error(err))
		return
	}

	opentracing.SetGlobalTracer(tracer)
	logger.Info("链路追踪初始化成功", zap.String("service", cfg.ServiceName))
}

// GetTracer 获取tracer实例
func GetTracer() opentracing.Tracer {
	return tracer
}

// Close 关闭tracer
func Close() {
	if closer != nil {
		closer.Close()
	}
}

// StartSpan 开始一个新的span
func StartSpan(operationName string) opentracing.Span {
	return tracer.StartSpan(operationName)
}

// StartSpanFromContext 从上下文开始一个新的span
func StartSpanFromContext(ctx opentracing.SpanContext, operationName string) opentracing.Span {
	return tracer.StartSpan(operationName, opentracing.ChildOf(ctx))
}
