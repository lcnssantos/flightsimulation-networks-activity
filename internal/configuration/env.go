package configuration

import (
	"os"

	"github.com/go-playground/validator"
)

var Environment Env

type Env struct {
	mongoURL string `validate:"required"`
	mongoSSL string `validate:"required"`
}

func (e *Env) GetMongoSSL() bool {
	return e.mongoSSL == "true"
}

func (e *Env) Validate() error {
	err := validator.New().Struct(e)
	return err
}

func LoadEnv() {
	Environment = Env{
		mongoURL: os.Getenv("MONGO_URL"),
		mongoSSL: os.Getenv("MONGO_SSL"),
	}
}
