package item

import "testing"

func TestValidatePrice(t *testing.T) {
	result := validatePrice("40", "$21")
	if !result {
		t.Errorf("Failed")
	}
}
