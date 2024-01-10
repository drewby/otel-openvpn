package pireceiver

import (
	"errors"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/receiver/scraperhelper"

	"github.com/drewby/pireceiver/internal/metadata"
)

const (
	defaultPath = "/sys/class/thermal"
)

// Config defines the configuration for the Raspberry Pi receiver.
type Config struct {
	Path                                    string                   `mapstructure:"path"` // Path to the file to be scraped for metrics (default: /sys/class/thermal)
	scraperhelper.ScraperControllerSettings `mapstructure:",squash"` // ScraperControllerSettings to configure scraping interval (default: 10s)
	metadata.MetricsBuilderConfig           `mapstructure:",squash"` // MetricsBuilderConfig to enable/disable specific metrics (default: all enabled)
}

func createDefaultConfig() component.Config {
	return &Config{
		Path:                      defaultPath,
		ScraperControllerSettings: scraperhelper.NewDefaultScraperControllerSettings(metadata.Type),
		MetricsBuilderConfig:      metadata.DefaultMetricsBuilderConfig(),
	}
}

func (c Config) Validate() error {
	if c.Path == "" {
		return errors.New("path cannot be empty")
	}

	return nil
}
