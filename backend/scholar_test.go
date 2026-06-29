package main

import "testing"

func TestCleanAcademicDegrees(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"Dr. apt. Eka Fitrianda, M.Farm", "Eka Fitrianda"},
		{"Prof. Dr. Ir. Soekarno, M.Sc", "Soekarno"},
		{"Ns. Eka Fitrianda", "Eka Fitrianda"},
		{"H. Eka Fitrianda", "Eka Fitrianda"},
		{"Dr. apt. Eka Fitrianda", "Eka Fitrianda"},
		{"Muhammad, M.Si.", "Muhammad"},
		{"Apt. Roni", "Roni"},
	}

	for _, c := range cases {
		got := cleanAcademicDegrees(c.input)
		if got != c.expected {
			t.Errorf("cleanAcademicDegrees(%q) = %q; want %q", c.input, got, c.expected)
		}
	}
}
