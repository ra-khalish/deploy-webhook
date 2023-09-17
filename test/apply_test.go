package test

import (
	"strings"
	"testing"

	"github.com/ra-khalish/deploy-webhook/services"
)

var ns = "test"

func TestApply(t *testing.T) {
	result := strings.Contains(services.Apply(&ns), "created")
	if result != true {
		t.Errorf("Apply not changes, got: %t, want: %t", result, true)
	}
}

func TestDelete(t *testing.T) {
	result := strings.Contains(services.Delete(&ns), "deleted")
	if result != true {
		t.Errorf("Delete object fail, got: %t, want: %t", result, true)
	}
}
