package ordering

import (
	"reflect"
	"testing"
)

func assertEqualSlices[T any](t *testing.T, got, want []T) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("length mismatch: got %d want %d", len(got), len(want))
	}
	for i := range got {
		if !reflect.DeepEqual(got[i], want[i]) {
			t.Fatalf("index %d mismatch: got %#v want %#v", i, got[i], want[i])
		}
	}
}
