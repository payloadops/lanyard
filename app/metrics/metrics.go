package metrics

// https://aws-otel.github.io/docs/getting-started/go-sdk/manual-instr
import (
	"context"
	"crypto/x509"
	"fmt"
	"github.com/payloadops/plato/app/config"
	"go.opentelemetry.io/contrib/detectors/aws/ecs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// NewMeter initializes the OpenTelemetry metrics provider.
func NewMeter(ctx context.Context, cfg *config.Config) (*metric.MeterProvider, error) {
	options := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(cfg.OpenTelemetry.ProviderEndpoint),
		otlpmetricgrpc.WithDialOption(grpc.WithBlock()),
	}

	// TODO: Should we rigidly enforce this?
	if cfg.OpenTelemetry.CACert != "" {
		// Create a certificate pool and add the CA certificate
		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM([]byte(cfg.OpenTelemetry.CACert)) {
			return nil, fmt.Errorf("failed to append CA certificate to pool")
		}

		options = append(options, otlpmetricgrpc.WithTLSCredentials(
			credentials.NewClientTLSFromCert(certPool, ""),
		))
	} else {
		options = append(options, otlpmetricgrpc.WithInsecure())
	}

	me, err := otlpmetricgrpc.New(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP metric exporter: %v", err)
	}

	res, err := ecs.NewResourceDetector().
		Detect(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP trace resource: %v", err)
	}

	mp := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(me)),
		metric.WithResource(res),
	)

	otel.SetMeterProvider(mp)
	return mp, nil
}
