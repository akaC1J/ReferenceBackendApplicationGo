package initialization

import (
	"route256/loms/internal/app/initialization"
	"testing"
)

func TestNew(t *testing.T) {

	testConf, err := initialization.LoadDefaultConfig()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	_, err = initialization.New(testConf)
	if err != nil {
		t.Fatalf("failed to create cart: %v", err)
	}
	t.Log("[TestNew] Application initialization successful")
}
