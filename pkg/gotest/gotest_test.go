package gotest

import "testing"

func TestSample(t *testing.T) {
	if !Sample() {
		t.Fatal("invalid result. Expected 'true', but got 'false'")
	}
}
