package config

import (
	"os"
	"testing"
)

type simple struct {
	A int    `cli:"A" env:"A"`
	B string `cli:"B" env:"B"`
}

func TestParseOverridesFlat(t *testing.T) {
	s := simple{}
	m := map[string]string{"B": "CLI_B"}
	os.Setenv("QAAS_A", "5")
	os.Setenv("QAAS_B", "ENV_B")
	defer os.Unsetenv("A")
	defer os.Unsetenv("B")

	if err := ParseOverrides(&s, m); err != nil {
		t.Fatalf("Error parsing overrides: %s", err)
	} else if s.A != 5 {
		t.Fatalf("Should override (struct).A with 5, have: %d", s.A)
	} else if s.B != "CLI_B" {
		t.Fatalf("Should override (struct).B with CLI_B, have: %s", s.B)
	}
}

type complex struct {
	A int    `cli:"A" env:"A"`
	B string `cli:"B" env:"B"`
	Simple struct {
		C int    `cli:"C" env:"C"`
		D string `cli:"D" env:"D"`
		E int64 `cli:"E" env:"E"`
	}
}

func TestParseOverridesNested(t *testing.T) {
	s := complex{}
	m := map[string]string{"C": "10"}
	os.Setenv("QAAS_A", "5")
	os.Setenv("QAAS_B", "ENV_B")
	os.Setenv("QAAS_C", "5")
	os.Setenv("QAAS_D", "ENV_D")
	os.Setenv("QAAS_E", "9")
	defer os.Unsetenv("A")
	defer os.Unsetenv("B")
	defer os.Unsetenv("C")
	defer os.Unsetenv("D")
	defer os.Unsetenv("E")

	if err := ParseOverrides(&s, m); err != nil {
		t.Fatalf("Error parsing overrides: %s", err)
	} else if s.A != 5 {
		t.Fatalf("Should override (struct).A with 5, have: %d", s.A)
	} else if s.B != "ENV_B" {
		t.Fatalf("Should override (struct).B with ENV_B, have: %s", s.B)
	} else if s.Simple.C != 10 {
		t.Fatalf("Should override (struct).Simple.C with 10, have: %d", s.Simple.C)
	} else if s.Simple.D != "ENV_D" {
		t.Fatalf("Should override (struct).Simple.D with ENV_D, have: %s", s.Simple.D)
	} else if s.Simple.E != 9 {
		t.Fatalf("Should override (struct).Simple.E with 9, have: %d", s.Simple.E)
	}
}
