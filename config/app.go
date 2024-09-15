package config

import (
	_ "embed"
	"github.com/in-rich/lib-go/deploy"
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

var App = deploy.LoadConfig[AppType](deploy.DevConfig(appDevFile))
