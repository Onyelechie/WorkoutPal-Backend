package test

import (
	"testing"
)

func TestFailingCIExample(t *testing.T) {
    t.Fatalf("This test is meant to fail and verify CI catches it.")
}