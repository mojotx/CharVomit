package arg

import "testing"

func TestVersion(t *testing.T) {
	if got, want := Version(), "dev"; got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}
