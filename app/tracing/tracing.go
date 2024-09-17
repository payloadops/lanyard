package tracing

// https://aws-otel.github.io/docs/getting-started/go-sdk/manual-instr
import (
	"context"
	"crypto/x509"
	"fmt"
	"github.com/payloadops/lanyard/app/config"
	"go.opentelemetry.io/contrib/detectors/aws/ecs"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	// TracingSampleRatio represents the sample size of trace IDs to record (currently 10%).
	TracingSampleRatio = 0.10
)

// NewTracer initializes the OpenTelemetry tracer.
func NewTracer(ctx context.Context, cfg *config.Config) (*trace.TracerProvider, error) {
	options := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(cfg.OpenTelemetry.ProviderEndpoint),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
	}

	// TODO: Should we rigidly enforce this?
	if cfg.OpenTelemetry.CACert != "" {
		// Create a certificate pool and add the CA certificate
		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM([]byte(cfg.OpenTelemetry.CACert)) {
			return nil, fmt.Errorf("failed to append CA certificate to pool")
		}

		options = append(options, otlptracegrpc.WithTLSCredentials(
			credentials.NewClientTLSFromCert(certPool, ""),
		))
	} else {
		options = append(options, otlptracegrpc.WithInsecure())
	}

	te, err := otlptracegrpc.New(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP trace exporter: %v", err)
	}

	res, err := ecs.NewResourceDetector().
		Detect(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP trace resource: %v", err)
	}

	idg := xray.NewIDGenerator()
	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.TraceIDRatioBased(TracingSampleRatio)),
		trace.WithResource(res),
		trace.WithBatcher(te),
		trace.WithIDGenerator(idg),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(xray.Propagator{})
	return tp, nil
}
