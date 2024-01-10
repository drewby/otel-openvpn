package pireceiver

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/receiver/receivertest"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	"github.com/drewby/pireceiver/internal/metadata"
)

func TestNewScraper(t *testing.T) {
	logger := zap.NewNop()
	metricsBuilder := &metadata.MetricsBuilder{}
	path := "/path/to/scrape"

	s := newScraper(metricsBuilder, path, logger)

	assert.NotNil(t, s)
	assert.Equal(t, s.raspberryPi.path, path)
	assert.Equal(t, s.logger, logger)
	assert.Equal(t, s.metricsBuilder, metricsBuilder)
}

func TestScrape(t *testing.T) {
	logger := zaptest.NewLogger(t)
	metricsBuilder := metadata.NewMetricsBuilder(metadata.DefaultMetricsBuilderConfig(), receivertest.NewNopCreateSettings())
	path := "testdata"

	s := newScraper(metricsBuilder, path, logger)

	ctx := context.Background()
	metrics, err := s.scrape(ctx)

	assert.NotNil(t, metrics)
	assert.Nil(t, err)

	assert.Equal(t, 1, metrics.MetricCount())
	assert.Equal(t, 3, metrics.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().Len())
	assert.Equal(t, 49.0, metrics.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).DoubleValue())
}
