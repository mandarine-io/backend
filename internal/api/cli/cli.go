package cli

import (
	"fmt"
	"github.com/mandarine-io/Backend/internal/api/config"
	"github.com/spf13/pflag"
	"os"

	"github.com/spf13/viper"
)

type Options struct {
	ConfigFilePath string
	EnvFilePath    string
}

var (
	helpOptionDesc = "Print usage and environment variables"

	configFilePathOptionDesc = "Configuration file path"
	defaultConfigFilePath    = getEnvWithDefault("MANDARINE_CONFIG__FILE", "config/config.yaml")

	envFilePathOptionDesc = "Environment variables file path"
	defaultEnvFilePath    = getEnvWithDefault("MANDARINE_ENV__FILE", "")
)

func MustParseCommandLine() *Options {
	var helpFlag bool
	var configFilePath string
	var envFilePath string

	pflag.BoolVarP(&helpFlag, "help", "h", false, helpOptionDesc)
	pflag.StringVarP(&configFilePath, "config", "c", defaultConfigFilePath, configFilePathOptionDesc)
	pflag.StringVarP(&envFilePath, "env", "e", defaultEnvFilePath, envFilePathOptionDesc)

	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if helpFlag {
		pflag.Usage()

		help := config.GetDescription()
		fmt.Println()
		fmt.Println(help)
		os.Exit(0)
	}

	return &Options{
		ConfigFilePath: configFilePath,
		EnvFilePath:    envFilePath,
	}
}

func getEnvWithDefault(envName, defaultValue string) string {
	if value, ok := os.LookupEnv(envName); ok {
		return value
	}
	return defaultValue
}
