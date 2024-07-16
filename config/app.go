package config

import (
	_ "embed"
	"github.com/goccy/go-yaml"
	"os"
	"time"
)

//go:embed app.dev.yaml
var appDevFile []byte

type NoteTierInformation struct {
	MaxEdits       int            `yaml:"max-edits"`
	CountEditsOver *time.Duration `yaml:"count-edits-over"`
}

type TierInformation struct {
	Notes NoteTierInformation `yaml:"notes"`
}

type AppType struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Postgres struct {
		DSN string `yaml:"dsn"`
	} `yaml:"postgres"`
	FreeTier TierInformation `yaml:"free-tier"`
}

var App AppType

func init() {
	switch os.Getenv("ENV") {
	case "prod":
		panic("not implemented")
	case "staging":
		panic("not implemented")
	default:
		if err := yaml.Unmarshal(appDevFile, &App); err != nil {
			panic(err)
		}
	}
}
