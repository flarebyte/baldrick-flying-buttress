package safety

import (
	"errors"
	"strings"
	"testing"
)

func TestCheckConfigFileSize(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		size     int64
		wantErr  bool
		wantText string
	}{
		{name: "within limit", size: MaxConfigFileBytes, wantErr: false},
		{name: "exceeds limit", size: MaxConfigFileBytes + 1, wantErr: true, wantText: "config file too large"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := CheckConfigFileSize(tc.size)
			assertLimitCheck(t, err, tc.wantErr, tc.wantText)
		})
	}
}

func TestCheckReportsCount(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		count    int
		wantErr  bool
		wantText string
	}{
		{name: "within limit", count: MaxReportsCount, wantErr: false},
		{name: "exceeds limit", count: MaxReportsCount + 1, wantErr: true, wantText: "reports count"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := CheckReportsCount(tc.count)
			assertLimitCheck(t, err, tc.wantErr, tc.wantText)
		})
	}
}

func TestCheckNotesCount(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		count    int
		wantErr  bool
		wantText string
	}{
		{name: "within limit", count: MaxNotesCount, wantErr: false},
		{name: "exceeds limit", count: MaxNotesCount + 1, wantErr: true, wantText: "notes count"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := CheckNotesCount(tc.count)
			assertLimitCheck(t, err, tc.wantErr, tc.wantText)
		})
	}
}

func TestCheckRelationshipsCount(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		count    int
		wantErr  bool
		wantText string
	}{
		{name: "within limit", count: MaxRelationshipsCount, wantErr: false},
		{name: "exceeds limit", count: MaxRelationshipsCount + 1, wantErr: true, wantText: "relationships count"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := CheckRelationshipsCount(tc.count)
			assertLimitCheck(t, err, tc.wantErr, tc.wantText)
		})
	}
}

func TestCheckCSVFileSize(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		size     int
		wantErr  bool
		wantText string
	}{
		{name: "within limit", size: MaxCSVFileBytes, wantErr: false},
		{name: "exceeds limit", size: MaxCSVFileBytes + 1, wantErr: true, wantText: "csv file too large"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := CheckCSVFileSize(tc.size)
			assertLimitCheck(t, err, tc.wantErr, tc.wantText)
		})
	}
}

func TestCheckCSVRenderedRows(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		rows     int
		wantErr  bool
		wantText string
	}{
		{name: "within limit", rows: MaxCSVRowsRenderedPerNote, wantErr: false},
		{name: "exceeds limit", rows: MaxCSVRowsRenderedPerNote + 1, wantErr: true, wantText: "csv rendered rows"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := CheckCSVRenderedRows(tc.rows)
			assertLimitCheck(t, err, tc.wantErr, tc.wantText)
		})
	}
}

func TestCheckGraphRenderNodeCount(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		nodes    int
		wantErr  bool
		wantText string
	}{
		{name: "within limit", nodes: MaxGraphRenderNodesPerSection, wantErr: false},
		{name: "exceeds limit", nodes: MaxGraphRenderNodesPerSection + 1, wantErr: true, wantText: "graph render node count"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := CheckGraphRenderNodeCount(tc.nodes)
			assertLimitCheck(t, err, tc.wantErr, tc.wantText)
		})
	}
}

func TestIsLimitError(t *testing.T) {
	t.Parallel()

	limitErr := newLimitError("limit")
	if !IsLimitError(limitErr) {
		t.Fatal("expected true for limit error")
	}
	if !IsLimitError(errors.Join(errors.New("wrapper"), limitErr)) {
		t.Fatal("expected true for wrapped limit error")
	}
	if IsLimitError(errors.New("other")) {
		t.Fatal("expected false for non-limit error")
	}
}

func assertLimitCheck(t *testing.T, err error, wantErr bool, wantText string) {
	t.Helper()
	if wantErr && err == nil {
		t.Fatal("expected error")
	}
	if !wantErr && err != nil {
		t.Fatalf("did not expect error: %v", err)
	}
	if !wantErr {
		return
	}
	if !IsLimitError(err) {
		t.Fatalf("expected limit error, got: %T", err)
	}
	if wantText != "" && !contains(err.Error(), wantText) {
		t.Fatalf("expected error to contain %q, got %q", wantText, err.Error())
	}
}

func contains(text string, want string) bool {
	return strings.Contains(text, want)
}
