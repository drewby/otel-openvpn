package pireceiver

import (
	"context"
	"time"

	"github.com/drewby/pireceiver/internal/metadata"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

type scraper struct {
	logger         *zap.Logger              // Logger to log events
	metricsBuilder *metadata.MetricsBuilder // MetricsBuilder to build metrics
	raspberryPi    *raspberryPi             // raspberryPi to get stats from /sys/class/thermal/thermal_zone*
}

// newScraper is a constructor function which returns a new scraper instance
func newScraper(metricsBuilder *metadata.MetricsBuilder, path string, logger *zap.Logger) *scraper {
	return &scraper{
		logger:         logger,
		metricsBuilder: metricsBuilder,
		raspberryPi:    newRaspberryPi(path, logger),
	}
}

func (s *scraper) scrape(ctx context.Context) (pmetric.Metrics, error) {
	s.logger.Debug("Scraping Raspberry Pi thermal zone information from path", zap.String("path", s.raspberryPi.path))

	thermalZoneList, err := s.raspberryPi.get()
	if err != nil {
		return pmetric.NewMetrics(), err
	}

	now := pcommon.NewTimestampFromTime(time.Now())

	// Iterate over the stats and add them to the metrics builder
	for _, thermalZone := range thermalZoneList {

		s.metricsBuilder.RecordRaspberryPiThermalZoneTemperatureDataPoint(now, thermalZone.temp, thermalZone.zone_type, int64(thermalZone.id))

	}

	metrics := s.metricsBuilder.Emit()

	return metrics, nil
}
