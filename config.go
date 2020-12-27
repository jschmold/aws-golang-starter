package server

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

var (
	_, caller, _, _ = runtime.Caller(0)

	// RootDir is the base path of the project
	RootDir = filepath.Dir(caller)
)

// ConfigPath is the full path to where we expect to find the config.toml file
var ConfigPath = fmt.Sprintf("%s/config.toml", RootDir)

const (
	// AWSKeyEnv is the env key string for the access key id for IAM
	AWSKeyEnv = "AWS_ACCESS_KEY_ID"

	// AWSSecretEnv is the secret access key for IAM
	AWSSecretEnv = "AWS_SECRET_ACCESS_KEY"

	// AWSRegionEnv is the env key for what region to use
	AWSRegionEnv = "AWS_REGION"

	// PgURLEnv is when we want to specify a local connection
	PgURLEnv = "PG_URL"

	// DBSecretEnv is the name of the secret containing the DB connection details
	DBSecretEnv = "DB_SECRET"
)

var viperKeys = []string{
	AWSSecretEnv,
	AWSRegionEnv,
	AWSKeyEnv,
	PgURLEnv,
	DBSecretEnv,
}

// LoadConfig tries to load config.toml and syncs the environment variables with the
// config.toml file
func LoadConfig() error {
	var err error

	viper.AutomaticEnv()

	path := fmt.Sprintf("%s/config.toml", RootDir)
	fmt.Printf("Searching for config: %s\n", path)
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Printf("Warning: Could not find config toml")
		return nil
	}

	if err != nil {
		return err
	}

	err = loadViper()
	return err
}

func loadViper() error {
	viper.SetConfigType("toml")
	viper.SetConfigFile(ConfigPath)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	for _, k := range viperKeys {
		fmt.Printf("  Sync key %s\n", k)
		syncViperKey(k)
	}

	return nil
}

func syncViperKey(key string) {
	val := viper.GetString(key)
	os.Setenv(key, val)
}
