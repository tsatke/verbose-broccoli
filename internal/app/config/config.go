package config

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/spf13/viper"
)

// Config keys.
const (
	ListenerHost       = "app.address.host"
	ListenerPort       = "app.address.port"
	AWSCognitoPoolID   = "aws.cognito.pool.id"
	AWSCognitoClientID = "aws.cognito.client.id"
	AWSS3Bucket        = "aws.s3.bucket"
	PGEndpoint         = "aws.postgres.endpoint"
	PGPort             = "aws.postgres.port"
	PGUsername         = "aws.postgres.username"
	PGPassword         = "aws.postgres.password"
	PGDatabase         = "aws.postgres.database"
)

type Config struct {
	*viper.Viper
	AWS aws.Config
}

func Load() (Config, error) {
	return load(func(v *viper.Viper) error {
		return v.ReadInConfig()
	})
}

func LoadFromReader(rd io.Reader) (Config, error) {
	return load(func(v *viper.Viper) error {
		return v.ReadConfig(rd)
	})
}

func load(readFn func(*viper.Viper) error) (Config, error) {
	v := viper.New()

	// set default values
	v.SetDefault(ListenerHost, "localhost")
	v.SetDefault(ListenerPort, 8080)

	// bind env
	v.AutomaticEnv()

	// load file
	v.SetConfigFile("config.yaml")
	if err := readFn(v); err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}

	// check whether all required values are set
	notSet := func(s string) error {
		return fmt.Errorf("%v not set", s)
	}
	for _, required := range []string{
		AWSCognitoPoolID,
		AWSCognitoClientID,
		AWSS3Bucket,
		PGEndpoint,
		PGPort,
		PGUsername,
		PGPassword,
	} {
		if !v.IsSet(required) {
			return Config{}, notSet(required)
		}
	}

	cfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion("eu-central-1"),
	)
	if err != nil {
		return Config{}, fmt.Errorf("load aws config: %w", err)
	}

	return Config{
		Viper: v,
		AWS:   cfg,
	}, nil
}
