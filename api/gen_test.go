package api

import (
	"testing"
)

func TestNewGenerator(t *testing.T) {
	g, _ := NewGenerator(
		"antiaddiction-service",
		"./test/",
		"6490bff6ad4730d1aece480a3f606abfc2c6ecb5c803cab241ced0f9a8960900",
		[]int{},
		[]int{16776, 16784},
	)
	err := g.GenTestApi()
	if err != nil {
		t.Fatal(err)
	}
}
