package config

import (
	"os"
	"testing"
)

func TestValues(t *testing.T) {
	_ = os.Setenv("DATABASE_URL", "mydb://")

	// check value from environment
	if Values().Database.Url != "mydb://" {
		t.Errorf("Incorrect Database.Url configuration: %s", Values().Database.Url)
	}

	// check default value
	if Values().Apiserver.Port != 8080 {
		t.Errorf("Incorrect Apiserver.Port configuration: %v", Values().Apiserver.Port)
	}
}
