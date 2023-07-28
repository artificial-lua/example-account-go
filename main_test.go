package main_test

import (
	"fmt"
	"testing"

	"github.com/artificial-lua/example-account-go/env"
	"github.com/joho/godotenv"
)

func TestEnv(t *testing.T) {
	err := godotenv.Load(".env")
	if err != nil {
		t.Log(err)
	}

	fmt.Println(env.Getenv("DB_HOST"))
	fmt.Println(env.Getenv("DB_PORT"))
	fmt.Println(env.Getenv("DB_USER"))
	fmt.Println(env.Getenv("DB_PASS"))
	fmt.Println(env.Getenv("DB_NAME"))
}
