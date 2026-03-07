package ordering

import "testing"

func assertEqualSlices[T comparable](t *testing.T, got, want []T) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("length mismatch: got %d want %d", len(got), len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("index %d mismatch: got %#v want %#v", i, got[i], want[i])
		}
	}
}
