package main

import (
	"testing"
	"booking-system/config"
	)

func Test_Env(t *testing.T) {
	got := config.LoadEnv()
	want := ".Env File Found"
	if got != want {
		t.Errorf("%s", got)
	}
}

