package tracing

import (
	"context"
	"github.com/payloadops/lanyard/app/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTracer_WithCACert(t *testing.T) {
	cfg := &config.Config{
		OpenTelemetry: config.OpenTelemetryConfig{
			ProviderEndpoint: "localhost:4317",
			CACert: `-----BEGIN CERTIFICATE-----
MIIDazCCAlOgAwIBAgIUFT0qSh83impMko+VTd7gDrOH9KAwDQYJKoZIhvcNAQEL
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0yNDA1MTcwNTM2NDNaFw0yNTA1
MTcwNTM2NDNaMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEw
HwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwggEiMA0GCSqGSIb3DQEB
AQUAA4IBDwAwggEKAoIBAQDDf7GFhUtt7s9AUrsmSGteHicmaIKBjntrJj9ILvpu
8mEKHFQOtfqzuAQEnL5crQFcDsj4o/OzFYW0hewyC68VKgOyyNwokkk0GB4fD+vA
H9PipeJSy9cxnoL3uJLQ8uOvNh4ERpEIfGdSuIay8avBV8uu1m3zodiyNuX/d9ee
9mhUEnHZ5jION8moeUqwfez9wDTN12RQn1PrNhxhxXDd7Jhdo+6nFGtaTxrherYA
MbX7FoZ46g8HLqXlf1Dnct69hgaR7HcAT3YC9ny534oXXvH1vR36ITFhasqNRmCu
7HsaQ4KBcqOKe8pZu2ieUw8EpNAAFBPgW4q4KbMdyQP7AgMBAAGjUzBRMB0GA1Ud
DgQWBBQLD0Y8k/OiJB2XSh1vOtLM1YY4LjAfBgNVHSMEGDAWgBQLD0Y8k/OiJB2X
Sh1vOtLM1YY4LjAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQAP
lx16ijOdiJKMW8TXtG8a6M5dGpu78TeV/mjDN6ytrqTZP8crclldU/Uvs4ij/TCX
VouiTVbS0qfSymsiPvBkTYO4xjXJoqsRwaw8RRCJ5uTiaSH2NdKkvjl8Eh8gohUr
OjS/4Bh4Z1Xk4majyRRulQRHxk1yx+VvLosye1cueqVEkTDjaoLbV0vp+mco072w
PC69EImCtM4N1MC82pAbz6DVVsXGsBSkXYMKkdkQJ/tO225fjYJOwPN9hH93Ef9o
Fas8BkEK8qDaiWcMK+eydDLKhf72+pZrvh0IB3i4swj/szf6XIcrNF1xoDAY/3yH
9qKpcBScxVoYU/18dhsZ
-----END CERTIFICATE-----`,
		},
	}

	tp, err := NewTracer(context.Background(), cfg)
	assert.NoError(t, err)
	assert.NotNil(t, tp)
}

func TestNewTracer_WithoutCACert(t *testing.T) {
	cfg := &config.Config{
		OpenTelemetry: config.OpenTelemetryConfig{
			ProviderEndpoint: "localhost:4317",
			CACert:           "",
		},
	}

	tp, err := NewTracer(context.Background(), cfg)
	assert.NoError(t, err)
	assert.NotNil(t, tp)
}

func TestNewTracer_InvalidCACert(t *testing.T) {
	cfg := &config.Config{
		OpenTelemetry: config.OpenTelemetryConfig{
			ProviderEndpoint: "localhost:4317",
			CACert:           "invalid-cert",
		},
	}

	tp, err := NewTracer(context.Background(), cfg)
	assert.Error(t, err)
	assert.Nil(t, tp)
}
