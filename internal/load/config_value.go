package load

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	cueload "cuelang.org/go/cue/load"
	"cuelang.org/go/cue/parser"
	"github.com/flarebyte/baldrick-flying-buttress/internal/safety"
)

func CompileConfigValue(configPath string) (cue.Value, string, error) {
	return compileConfigValue(cuecontext.New(), configPath)
}

func CompileConfigValueInContext(cueCtx *cue.Context, configPath string) (cue.Value, string, error) {
	if cueCtx == nil {
		return cue.Value{}, "", fmt.Errorf("cue context is required")
	}
	return compileConfigValue(cueCtx, configPath)
}

func compileConfigValue(cueCtx *cue.Context, configPath string) (cue.Value, string, error) {
	resolvedPath, err := resolveConfigPath(configPath)
	if err != nil {
		return cue.Value{}, "", err
	}

	info, err := os.Stat(resolvedPath)
	if err != nil {
		return cue.Value{}, "", fmt.Errorf("stat config %s: %w", resolvedPath, err)
	}
	if err := safety.CheckConfigFileSize(info.Size()); err != nil {
		return cue.Value{}, "", err
	}

	if strings.EqualFold(filepath.Ext(resolvedPath), ".cue") {
		packageName, hasPackage, err := cuePackageName(resolvedPath)
		if err != nil {
			return cue.Value{}, "", fmt.Errorf("parse config %s: %w", resolvedPath, err)
		}
		if hasPackage {
			value, err := compileCuePackage(cueCtx, filepath.Dir(resolvedPath), packageName)
			if err != nil {
				return cue.Value{}, "", fmt.Errorf("parse config %s: %w", resolvedPath, err)
			}
			return value, resolvedPath, nil
		}
	}

	data, err := os.ReadFile(resolvedPath)
	if err != nil {
		return cue.Value{}, "", fmt.Errorf("read config %s: %w", resolvedPath, err)
	}

	value := cueCtx.CompileBytes(data, cue.Filename(resolvedPath))
	if err := value.Err(); err != nil {
		return cue.Value{}, "", fmt.Errorf("parse config %s: %w", resolvedPath, err)
	}
	return value, resolvedPath, nil
}

func resolveConfigPath(configPath string) (string, error) {
	if configPath == "" {
		return "", fmt.Errorf("config path is required")
	}

	info, err := os.Stat(configPath)
	if err != nil {
		return "", fmt.Errorf("stat config %s: %w", configPath, err)
	}
	if !info.IsDir() {
		return configPath, nil
	}

	return filepath.Join(configPath, "app.cue"), nil
}

func cuePackageName(configPath string) (string, bool, error) {
	file, err := parser.ParseFile(configPath, nil, parser.PackageClauseOnly)
	if err != nil {
		return "", false, err
	}

	packageName := file.PackageName()
	return packageName, packageName != "", nil
}

func compileCuePackage(cueCtx *cue.Context, dir, packageName string) (cue.Value, error) {
	insts := cueload.Instances([]string{"."}, &cueload.Config{
		Dir:     dir,
		Package: packageName,
	})
	if len(insts) != 1 {
		return cue.Value{}, fmt.Errorf("load cue package %s: expected 1 instance, got %d", dir, len(insts))
	}
	inst := insts[0]
	if inst.Err != nil {
		return cue.Value{}, inst.Err
	}

	var totalBytes int64
	for _, file := range inst.BuildFiles {
		info, err := os.Stat(file.Filename)
		if err != nil {
			return cue.Value{}, fmt.Errorf("stat config %s: %w", file.Filename, err)
		}
		totalBytes += info.Size()
	}
	if err := safety.CheckConfigFileSize(totalBytes); err != nil {
		return cue.Value{}, err
	}

	value := cueCtx.BuildInstance(inst)
	if err := value.Err(); err != nil {
		return cue.Value{}, err
	}
	return value, nil
}
