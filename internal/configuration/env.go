package configuration

import (
	"errors"
	"os"
	"strconv"
)

var Environment Env

type Env struct {
	MongoURL string `validate:"required"`
	mongoSSL string `validate:"required"`
	port     string `validate:"required"`
}

func (e *Env) GetMongoSSL() bool {
	return e.mongoSSL == "true"
}

func (e *Env) GetPort() int {
	port, err := strconv.Atoi(e.port)

	if err != nil {
		return 8080
	}

	return port
}

func (e *Env) Validate() error {
	if e.mongoSSL == "" || e.MongoURL == "" {
		return errors.New("INVALID ENVIRONMENT CONFIGURATION")
	}

	return nil
}

func LoadEnv() {
	Environment = Env{
		MongoURL: os.Getenv("MONGO_URL"),
		mongoSSL: os.Getenv("MONGO_SSL"),
		port:     os.Getenv("PORT"),
	}
}
