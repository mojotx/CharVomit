package arg

import "testing"

func TestVersion(t *testing.T) {
	v := Version()
	if v != CharVomitVersion {
		t.Errorf("expected '%s', got '%s'", CharVomitVersion, v)
	}
}
