package telemetry

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

// Mawinter 内で使用するサービス名の定数
const (
	ServiceNameAPI = "mawinter-api"
)

// ShutdownFunc はトレーサープロバイダーの終了処理を表す
type ShutdownFunc func(ctx context.Context) error

// Init は指定されたサービス向けのグローバルなトレーサープロバイダーを初期化する
// endpoint が空文字の場合はトレースを無効化し、終了処理は no-op となる
func Init(ctx context.Context, endpoint, serviceName, version string) (ShutdownFunc, bool, error) {
	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		return func(context.Context) error { return nil }, false, nil
	}

	exporter, err := newHTTPExporter(ctx, endpoint)
	if err != nil {
		return nil, false, fmt.Errorf("failed to build OTLP exporter: %w", err)
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(version),
		),
	)
	if err != nil {
		return nil, false, fmt.Errorf("failed to build resource: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp.Shutdown, true, nil
}

func newHTTPExporter(ctx context.Context, endpoint string) (sdktrace.SpanExporter, error) {
	clientOpts, err := buildHTTPClientOptions(endpoint)
	if err != nil {
		return nil, err
	}

	client := otlptracehttp.NewClient(clientOpts...)
	return otlptrace.New(ctx, client)
}

func buildHTTPClientOptions(raw string) ([]otlptracehttp.Option, error) {
	urlPath := "/v1/traces"
	useTLS := true
	endpoint := raw

	if strings.Contains(raw, "://") {
		parsed, err := url.Parse(raw)
		if err != nil {
			return nil, fmt.Errorf("invalid OTLP endpoint: %w", err)
		}
		if parsed.Host == "" {
			return nil, fmt.Errorf("OTLP endpoint host is empty: %s", raw)
		}
		endpoint = parsed.Host
		if parsed.Path != "" && parsed.Path != "/" {
			urlPath = parsed.Path
		}
		switch parsed.Scheme {
		case "", "http":
			useTLS = false
		case "https":
			useTLS = true
		default:
			return nil, fmt.Errorf("unsupported OTLP endpoint scheme %q", parsed.Scheme)
		}
	}

	// スキームが無い場合は HTTP で送信
	if !strings.Contains(raw, "://") {
		useTLS = false
	}

	endpoint = strings.TrimRight(endpoint, "/")
	if endpoint == "" {
		return nil, fmt.Errorf("OTLP endpoint host is empty: %s", raw)
	}

	opts := []otlptracehttp.Option{
		otlptracehttp.WithURLPath(urlPath),
		otlptracehttp.WithEndpoint(endpoint),
	}
	if !useTLS {
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	return opts, nil
}
