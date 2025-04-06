package main

import (
	"os"
	"testing"
)

func setupMain() {
}
func teardownMain() {
}
func TestMain(m *testing.M) {
	setupMain()
	code := m.Run()
	teardownMain()
	os.Exit(code)
}
func TestFormatPrice(t *testing.T) {
	expected := "$79,274.57"
	actual := FormatPrice(79274.57111071543)
	if expected != actual {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}
}
func TestInit(t *testing.T) {
	expected_nsec := "blah"
	expected_cmc := "blah2"
	if expected_nsec == nsec {
		t.Errorf("nsec should not be equal to %s", expected_nsec)
	}
	if expected_cmc == cmc {
		t.Errorf("cmc should not be equal to %s", expected_cmc)
	}
	os.Setenv("NSEC", expected_nsec)
	os.Setenv("CMC_API", expected_cmc)
	initConfig()
	if expected_nsec != nsec {
		t.Errorf("nsec should be equal to %s", expected_nsec)
	}
	if expected_cmc != cmc {
		t.Errorf("cmc should be equal to %s", expected_cmc)
	}
}

//func TestAdder(t *testing.T) {
//	expected := 5
//	actual := adder(2, 3)
//	if expected != actual {
//		t.Errorf("Expected %d, got %d", expected, actual)
//	}
//}
