package main

import (
	"testing"

	"github.com/spf13/viper"
)

func TestSetup(t *testing.T) {
	v := viper.New()
	setup(v)

	if c == nil {
		t.Fatalf("c is nil")
	}

	if duration == 0 {
		t.Fatalf("expected non zero duration. got %v", duration)
	}
}
