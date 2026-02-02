package library

import "testing"

func TestStableExternalIDDeterministic(t *testing.T) {
	a := stableExternalID("variant", "role", "preset", "tpl")
	b := stableExternalID("variant", "role", "preset", "tpl")
	if a != b {
		t.Fatalf("expected deterministic id, got %q != %q", a, b)
	}
	c := stableExternalID("variant", "role", "preset", "tpl2")
	if a == c {
		t.Fatalf("expected different ids for different input")
	}
}
