package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func LoadViper() (*Config, error) {
	fs := pflag.NewFlagSet("user-svc-flags", pflag.ContinueOnError)
	fs.SortFlags = false

	// App
	fs.String("name", "user-svc", "service name")

	// Logging level
	fs.String("level", "info", "log level debug, info, warn, error, fatal or panic")

	// Config path/file
	currentWDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("get working directory: %w", err)
	}

	configPath := filepath.Join(currentWDir, "config")
	fs.String("config-path", configPath, "config dir path")
	fs.String("config", "local.yaml", "config file name")

	// Network
	fs.String("host", "", "Host to bind the server")
	fs.String("port", "8080", "Port on app will listhen")

	// NATS
	fs.String("nats-url", "nats://localhost:4222", "NATS connection url")

	// Database
	fs.String("db-host", "localhost", "database host")
	fs.String("db-port", "54432", "database port")
	fs.String("db-user", "user_svc", "database user")
	fs.String("db-name", "users_db", "database name")
	fs.String("db-password", "password_dev", "database name")
	fs.String("db-ssl-mode", "disable", "database ssl mode")
	fs.BoolP("dump", "d", false, "dump config")

	err = fs.Parse(os.Args[1:])
	if err != nil {
		switch {
		case errors.Is(err, pflag.ErrHelp):
			os.Exit(0)
		case err != nil:
			fmt.Fprintf(os.Stderr, "Error: %s\n\n", err.Error())
			fs.PrintDefaults()
			os.Exit(2)
		}
	}

	viper.BindPFlags(fs)
	viper.SetConfigType("yaml")

	hostname, _ := os.Hostname()
	viper.Set("hostname", hostname)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	configFileFullPath := filepath.Join(viper.GetString("config-path"), viper.GetString("config"))
	fmt.Println(configFileFullPath)
	if _, fileErr := os.Stat(configFileFullPath); fileErr == nil {
		configFileName := strings.Split(viper.GetString("config"), ".")[0]
		viper.SetConfigName(configFileName)
		viper.AddConfigPath(viper.GetString("config-path"))

		if readErr := viper.ReadInConfig(); readErr != nil {
			return nil, fmt.Errorf("Error reading config file: %w", readErr)
		}
	}

	os.Setenv("NATS_URL", viper.GetString("nats-url"))

	cfg := &Config{
		App: &Application{
			Name:     viper.GetString("name"),
			Host:     viper.GetString("host"),
			Port:     viper.GetString("port"),
			NatsUrl:  viper.GetString("nats-url"),
			LogLevel: viper.GetString("level"),
		},
		Db: &Database{
			Host:     viper.GetString("db-host"),
			Port:     viper.GetString("db-port"),
			User:     viper.GetString("db-user"),
			Name:     viper.GetString("db-name"),
			Password: viper.GetString("db-password"),
			SSLMode:  viper.GetString("db-ssl-mode"),
		},
	}

	if viper.GetBool("dump") {
		printConfig(viper.AllSettings())
	}

	return cfg, nil
}

func sensitiveConfigKey(k string) bool {
	lk := strings.ToLower(k)
	for _, marker := range []string{
		"password", "secret", "token", "credential", "passphrase",
		"private-key", "privatekey", "api-key", "apikey",
	} {
		if strings.Contains(lk, marker) {
			return true
		}
	}
	return false
}

func printConfig(m interface{}) {
	for k, v := range m.(map[string]any) {
		if !sensitiveConfigKey(k) {
			fmt.Printf("%s: %v\n", k, v)
			continue
		}

		fmt.Printf("%s: ***\n", k)
	}
}
