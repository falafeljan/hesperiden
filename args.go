package main

import (
	"encoding/json"
	"fmt"
	"github.com/caarlos0/env"
	"github.com/falafeljan/docker-recreate"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Args describe arguments we expect from the environment
type Args struct {
	HTTPPort     string `env:"HTTP_PORT" envDefault:"80"`
	InProduction bool   `env:"PRODUCTION" envDefault:"false"`
	RateLimit    string `env:"RATE_LIMIT" envDefault:"30-M"`
	RedisHost    string `env:"REDIS_HOST" envDefault:"127.0.0.1"`
	RedisPort    string `env:"REDIS_PORT" envDefault:"6379"`
	RedisPrefix  string `env:"REDIS_PREFIX" envDefault:"token"`
	Registries   []recreate.RegistryConf
}

func getRegistries() (registries []recreate.RegistryConf) {
	emptyRegistries := []recreate.RegistryConf{}
	ex, err := os.Executable()
	if err != nil {
		return emptyRegistries
	}

	cwd := filepath.Dir(ex)
	filePath := strings.Join([]string{
		cwd,
		"registries.json"},
		"/")

	file, err := os.Open(filePath)
	if err != nil {
		return emptyRegistries
	}

	defer file.Close()

	var parsedRegistries []recreate.RegistryConf
	byteValue, _ := ioutil.ReadAll(file)
	err = json.Unmarshal(byteValue, &parsedRegistries)

	if err != nil {
		return emptyRegistries
	}

	return parsedRegistries
}

func getArgs() Args {
	args := Args{}
	err := env.Parse(&args)

	if err != nil {
		panic(fmt.Sprintf("Error while parsing configuration: %+v\n", err))
	}

	args.Registries = getRegistries()

	return args
}
