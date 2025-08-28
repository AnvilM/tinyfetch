package config

import (
	"os"
	"tinyfetch/internal/utils/logger"

	"gopkg.in/yaml.v3"
)

var configPath = os.Getenv("HOME") + "/.config/tinyfetch/config.yml"

// Allowed terminal colors
var allowedColors = []string{
	"black", "red", "green", "yellow",
	"blue", "magenta", "cyan", "white",
}

// Custom type for color
type Color string

// Generic pointer helper
func Ptr[T any](v T) *T { return &v }

// Configuration structures
type Title struct {
	FilePath *string `yaml:"filePath"`
	Color    *Color  `yaml:"color"`
}

type Module struct {
	Type            *string `yaml:"type"`
	InfoColor       *Color  `yaml:"infoColor"`
	Label           *string `yaml:"label"`
	LabelColor      *Color  `yaml:"labelColor"`
	Prefix            *string `yaml:"prefix"`
	PrefixColor       *Color  `yaml:"prefixColor"`
}

type Container struct {
	MarginLeft   *int   `yaml:"marginLeft"`
	MarginRight  *int   `yaml:"marginRight"`
	PaddingLeft  *int   `yaml:"paddingLeft"`
	PaddingRight *int   `yaml:"paddingRight"`
	BorderColor  *Color `yaml:"borderColor"`
}

type Config struct {
	Title     *Title      `yaml:"title"`
	Modules   *[]Module   `yaml:"modules"`
	Container *Container  `yaml:"container"`
}

// Default module types
var allowedTypes = []string{
	"user",
	"hostname",
	"os",
	"kernel",
	"uptime",
	"shell",
	"packages",
	"memory",
	"colors",
}

// LoadConfig reads the configuration and returns Config
func LoadConfig() Config {
	var cfg Config

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return GetDefaultConfig()
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		logger.Fatal("Error reading file: %v", err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		logger.Fatal("Error parsing YAML: %v", err)
	}

	checkRequired(&cfg)

	return cfg
}

// checkRequired ensures required fields and sets defaults
func checkRequired(cfg *Config) {
	defaultCfg := GetDefaultConfig()

	// Title
	if cfg.Title != nil {
		if cfg.Title.FilePath == nil {
			logger.Fatal("title.filePath is required!")
		}

		if cfg.Title.Color == nil || !contains(allowedColors, string(*cfg.Title.Color)) {
			cfg.Title.Color = Ptr(*defaultCfg.Title.Color)
		}
	}

	// Container
	if cfg.Container != nil {
		if cfg.Container.BorderColor == nil {
			cfg.Container.BorderColor = Ptr(*defaultCfg.Container.BorderColor)
		}
		if cfg.Container.MarginLeft == nil {
			cfg.Container.MarginLeft = Ptr(*defaultCfg.Container.MarginLeft)
		}
		if cfg.Container.MarginRight == nil {
			cfg.Container.MarginRight = Ptr(*defaultCfg.Container.MarginRight)
		}
		if cfg.Container.PaddingLeft == nil {
			cfg.Container.PaddingLeft = Ptr(*defaultCfg.Container.PaddingLeft)
		}
		if cfg.Container.PaddingRight == nil {
			cfg.Container.PaddingRight = Ptr(*defaultCfg.Container.PaddingRight)
		}
	} else {
		cfg.Container = Ptr(*defaultCfg.Container)
	}

	// Modules
	if cfg.Modules == nil || len(*cfg.Modules) < 1 {
		logger.Fatal("modules must contain at least one element!")
	}

	for i := range *cfg.Modules {
		module := &(*cfg.Modules)[i]

		if module.Type == nil {
			logger.Fatal("module must contain type")
		}

		if !contains(allowedTypes, *module.Type) {
			logger.Fatal("invalid type: " + *module.Type)
		}

		if module.InfoColor == nil {
			module.InfoColor = Ptr(*(*defaultCfg.Modules)[0].InfoColor)
		}

		if module.Label == nil {
			module.Label = module.Type
		}
	}
}

// contains checks if slice contains a string
func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
