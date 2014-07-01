package swu

import (
	"testing"
)

func TestNewSWU(t *testing.T) {
	api := New("key")
	if api == nil {
		t.Error("New should not return nil")
	}
}
