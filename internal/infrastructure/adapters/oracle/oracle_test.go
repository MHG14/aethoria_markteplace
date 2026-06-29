package oracle_test

import (
	"testing"

	"github.com/MHG14/aethoria_marketplace/internal/infrastructure/adapters/oracle"
)

func TestMockOracle_ReturnsConfiguredPrice(t *testing.T) {
	o := oracle.NewMock()
	o.SetPrice(1, 1500)

	price, err := o.GetBasePrice(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if price != 1500 {
		t.Errorf("expected 1500, got %d", price)
	}
}

func TestMockOracle_ReturnsError_WhenDown(t *testing.T) {
	o := oracle.NewMock()
	o.SetAlwaysFail()

	_, err := o.GetBasePrice(1)
	if err == nil {
		t.Error("expected error from unavailable oracle")
	}
}

func TestMockOracle_ReturnsZero_WhenConfigured(t *testing.T) {
	o := oracle.NewMock()
	o.SetReturnZero(1)

	price, _ := o.GetBasePrice(1)
	if price != 0 {
		t.Errorf("expected 0, got %d", price)
	}
}

func TestMockOracle_ReturnsNegative_WhenConfigured(t *testing.T) {
	o := oracle.NewMock()
	o.SetReturnNegative(1)

	price, _ := o.GetBasePrice(1)
	if price != -100 {
		t.Errorf("expected -100, got %d", price)
	}
}

func TestMockOracle_RecoverAfterTransientFailure(t *testing.T) {
	o := oracle.NewMock()
	o.SetPrice(1, 1000)
	o.SetFailOnce(1)

	_, err := o.GetBasePrice(1)
	if err == nil {
		t.Error("expected first call to fail")
	}

	price, err := o.GetBasePrice(1)
	if err != nil {
		t.Fatalf("expected recovery after transient failure, got %v", err)
	}
	if price != 1000 {
		t.Errorf("expected 1000, got %d", price)
	}
}
