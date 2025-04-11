package tracing

import (
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

// InitJaeger menginisialisasi Jaeger tracer
func InitJaeger(serviceName string) (opentracing.Tracer, io.Closer, error) {
	cfg := jaegercfg.Configuration{
		ServiceName: "task-crud",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "jaeger:6831", // ✅ gunakan nama servicenya di docker
		},
	}
	

	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		return nil, nil, fmt.Errorf("cannot initialize Jaeger: %w", err)
	}

	opentracing.SetGlobalTracer(tracer)
	return tracer, closer, nil
}
