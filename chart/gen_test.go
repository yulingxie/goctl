package chart

import (
	"testing"
)

func TestGenerator_GenCharts(t *testing.T) {
	generator, _ := NewGenerator("test1-service", "./", "charts")
	generator.GenSngServiceCharts()
}

func TestGenerator_GenSngGatewayCharts(t *testing.T) {
	gwGenerator, _ := NewGenerator("sapi-gw", "./", "")
	gwGenerator.GenSngGatewayCharts()
}
