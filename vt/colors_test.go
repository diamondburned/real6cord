package vt

import (
	"testing"
)

func TestGetRGBInt(t *testing.T) {
	if xt := GetRGBInt(0xFFFFFF); xt != 15 {
		t.Fatalf("Unexpected value: %d", xt)
	}

	if xt := GetRGBInt(0x00af00); xt != 34 {
		t.Fatalf("Unexpected value: %d", xt)
	}
}
