package swu

import (
	"log"
	"os"
	"testing"
)

func TestNewSWU(t *testing.T) {
	api := New("key")
	if api == nil {
		t.Error("New should not return nil")
	}
}

func TestTemplates(t *testing.T) {
	api := New(os.Getenv("SWU_KEY"))
	err := api.Emails()
	if err != nil {
		t.Error(err)
	}
}
