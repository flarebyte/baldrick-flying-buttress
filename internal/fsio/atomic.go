package fsio

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

var (
	mkdirAllFn   = os.MkdirAll
	createTempFn = os.CreateTemp
	renameFn     = os.Rename
	chmodFn      = os.Chmod
)

func WriteFileAtomic(ctx context.Context, destination string, data []byte, perm os.FileMode) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	dir := filepath.Dir(destination)
	if err := mkdirAllFn(dir, 0o755); err != nil {
		return fmt.Errorf("create output directory %s: %w", dir, err)
	}

	tmpFile, err := createTempFn(dir, ".flyb-write-*")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()
	cleanup := true
	defer func() {
		if cleanup {
			_ = os.Remove(tmpPath)
		}
	}()

	if _, err := tmpFile.Write(data); err != nil {
		_ = tmpFile.Close()
		return err
	}
	if err := tmpFile.Sync(); err != nil {
		_ = tmpFile.Close()
		return err
	}
	if err := tmpFile.Close(); err != nil {
		return err
	}
	if err := chmodFn(tmpPath, perm); err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := renameFn(tmpPath, destination); err != nil {
		return err
	}
	cleanup = false

	syncDirBestEffort(dir)
	return nil
}

func syncDirBestEffort(dir string) {
	d, err := os.Open(dir)
	if err != nil {
		return
	}
	defer func() {
		_ = d.Close()
	}()
	_ = d.Sync()
}
