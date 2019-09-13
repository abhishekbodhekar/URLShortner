package main

import "testing"

func TestDectyption(t *testing.T) {
	decryptedVal, _ := decrypt("QmYEqVz0GGF2lARpic_J_fe8iEZr-otpRWx11svn4AQ")
	if decryptedVal != "Abhishek" {
		t.Errorf("Sum was incorrect, got: %s, want: %s.", decryptedVal, "Abhishek")
	}
}
