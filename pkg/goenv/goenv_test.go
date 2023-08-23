package goenv

import (
	"fmt"
	"testing"
)

func TestLoadString(t *testing.T) {
	t.Setenv("TEST_APP_STRING", "STRING_VALUE")
	if Load("TEST_APP_STRING", "DEFAULT_STRING_VALUE") != "STRING_VALUE" {
		t.Fatal("failed to load env var (string)")
	}
}

func TestLoadInt(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Setenv("TEST_APP_INT", fmt.Sprint(i))
		value := Load("TEST_APP_INT", int64(-1))
		if value != int64(i) {
			t.Fatal("failed to load env var (int64)")
		}
	}
}

func TestLoadFloat(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Setenv("TEST_APP_FLOAT", fmt.Sprintf("%d.75", i))
		value := Load("TEST_APP_FLOAT", float64(-1))
		if value != float64(i)+0.75 {
			t.Fatal("failed to load env var (float64)")
		}
	}
}

func TestLoadBool(t *testing.T) {
	t.Setenv("TEST_APP_BOOL", "true")
	value := Load("TEST_APP_BOOL", false)
	if value != true {
		t.Fatal("failed to load env var (bool)")
	}
}
