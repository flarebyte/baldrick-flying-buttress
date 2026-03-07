package textutil

import (
	"reflect"
	"testing"
)

func TestParseKeyValue(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		entry     string
		wantKey   string
		wantValue string
		wantOK    bool
	}{
		{name: "valid", entry: "a=b", wantKey: "a", wantValue: "b", wantOK: true},
		{name: "valid trims spaces", entry: "  key = value  ", wantKey: "key", wantValue: "value", wantOK: true},
		{name: "invalid missing equals", entry: "ab", wantOK: false},
		{name: "invalid empty key", entry: "=v", wantOK: false},
		{name: "invalid empty value", entry: "k=", wantOK: false},
		{name: "invalid blank key", entry: "   =v", wantOK: false},
		{name: "invalid blank value", entry: "k=   ", wantOK: false},
		{name: "keeps value with extra equals", entry: "k=v=x", wantKey: "k", wantValue: "v=x", wantOK: true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gotKey, gotValue, gotOK := ParseKeyValue(tc.entry)
			if gotKey != tc.wantKey || gotValue != tc.wantValue || gotOK != tc.wantOK {
				t.Fatalf("unexpected parse result key=%q value=%q ok=%v", gotKey, gotValue, gotOK)
			}
		})
	}
}

func TestSplitCSV(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		input string
		want  []string
	}{
		{name: "empty", input: "", want: []string{}},
		{name: "blank", input: "   ", want: []string{}},
		{name: "single", input: "a", want: []string{"a"}},
		{name: "multiple", input: "a,b,c", want: []string{"a", "b", "c"}},
		{name: "trims and drops empties", input: " a, ,b,  , c ", want: []string{"a", "b", "c"}},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := SplitCSV(tc.input)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("unexpected split result\nwant: %#v\n got: %#v", tc.want, got)
			}
		})
	}
}

func TestSplitNonEmptyLines(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		input string
		want  []string
	}{
		{name: "empty", input: "", want: []string{}},
		{name: "blank", input: " \n \n", want: []string{}},
		{name: "single", input: "a", want: []string{"a"}},
		{name: "multiple", input: "a\nb\nc", want: []string{"a", "b", "c"}},
		{name: "trims and drops empties", input: " a \n\n b \n  \n c ", want: []string{"a", "b", "c"}},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := SplitNonEmptyLines(tc.input)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("unexpected split result\nwant: %#v\n got: %#v", tc.want, got)
			}
		})
	}
}
