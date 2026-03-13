package validate

import (
	"embed"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	cueerrors "cuelang.org/go/cue/errors"
	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/load"
)

//go:embed app_schema.cue
var schemaFS embed.FS

func validateStructure(raw domain.RawApp) []domain.Diagnostic {
	schemaRoot, appValue, err := compileSchemaAndApp(raw)
	if err != nil {
		return []domain.Diagnostic{newDiagnostic(schemaValidationSource, "FBV001", err.Error(), "config")}
	}

	constraint := schemaRoot.LookupPath(cue.ParsePath("#RawApp"))
	if !constraint.Exists() {
		return []domain.Diagnostic{newDiagnostic(schemaValidationSource, "FBV001", "missing cue schema definition: #RawApp", "config")}
	}

	validated := constraint.Unify(appValue)
	err = validated.Validate(cue.Concrete(true))
	if err == nil {
		return []domain.Diagnostic{}
	}

	diagnostics := make([]domain.Diagnostic, 0)
	for _, cueErr := range cueerrors.Errors(err) {
		path := toCanonicalPath(cueerrors.Path(cueErr))
		d, ok := mapSchemaDiagnostic(path)
		if ok {
			diagnostics = append(diagnostics, d)
		}
	}

	return diagnostics
}

func compileSchemaAndApp(raw domain.RawApp) (cue.Value, cue.Value, error) {
	schemaBytes, err := schemaFS.ReadFile("app_schema.cue")
	if err != nil {
		return cue.Value{}, cue.Value{}, fmt.Errorf("read cue schema: %w", err)
	}

	cueCtx := cuecontext.New()
	schemaValue := cueCtx.CompileBytes(schemaBytes, cue.Filename("internal/validate/app_schema.cue"))
	if err := schemaValue.Err(); err != nil {
		return cue.Value{}, cue.Value{}, fmt.Errorf("compile cue schema: %w", err)
	}

	appValue, err := buildRawAppValue(cueCtx, raw)
	if err != nil {
		return cue.Value{}, cue.Value{}, err
	}

	return schemaValue, appValue, nil
}

func buildRawAppValue(cueCtx *cue.Context, raw domain.RawApp) (cue.Value, error) {
	if raw.ConfigPath == "" {
		data, err := json.Marshal(raw)
		if err != nil {
			return cue.Value{}, fmt.Errorf("encode raw app: %w", err)
		}
		value := cueCtx.CompileBytes(data, cue.Filename("raw-app.json"))
		if err := value.Err(); err != nil {
			return cue.Value{}, fmt.Errorf("encode raw app: %w", err)
		}
		return value, nil
	}

	value, _, err := load.CompileConfigValueInContext(cueCtx, raw.ConfigPath)
	if err != nil {
		return cue.Value{}, err
	}
	return value, nil
}

func toCanonicalPath(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	if parts[0] == "#RawApp" {
		parts = parts[1:]
	}
	var b strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}
		if _, err := strconv.Atoi(part); err == nil {
			b.WriteString("[")
			b.WriteString(part)
			b.WriteString("]")
			continue
		}
		if b.Len() > 0 {
			b.WriteString(".")
		}
		b.WriteString(part)
	}
	return b.String()
}

func mapSchemaDiagnostic(path string) (domain.Diagnostic, bool) {
	switch {
	case path == "source":
		return newDiagnostic(schemaValidationSource, "FBV000", "missing required field: source", "source"), true
	case path == "reports":
		return newDiagnostic(schemaValidationSource, "FBV100", "missing required collection: reports", "reports"), true
	case path == "notes":
		return newDiagnostic(schemaValidationSource, "FBV200", "missing required collection: notes", "notes"), true
	case path == "relationships":
		return newDiagnostic(schemaValidationSource, "FBV300", "missing required collection: relationships", "relationships"), true
	case strings.HasSuffix(path, ".title") && strings.Contains(path, ".sections["):
		return newDiagnostic(schemaValidationSource, "FBV104", "missing required field: section title", path), true
	case strings.HasSuffix(path, ".title") && strings.HasPrefix(path, "reports["):
		return newDiagnostic(schemaValidationSource, "FBV101", "missing required field: report title", path), true
	case strings.HasSuffix(path, ".filepath") && strings.HasPrefix(path, "reports["):
		return newDiagnostic(schemaValidationSource, "FBV102", "missing required field: report filepath", path), true
	case isReportSectionsPath(path):
		return newDiagnostic(schemaValidationSource, "FBV103", "missing required field: report sections", path), true
	case strings.HasSuffix(path, ".name") && strings.HasPrefix(path, "notes["):
		return newDiagnostic(schemaValidationSource, "FBV201", "missing required field: note name", path), true
	case strings.HasSuffix(path, ".title") && strings.HasPrefix(path, "notes["):
		return newDiagnostic(schemaValidationSource, "FBV202", "missing required field: note title", path), true
	case strings.HasSuffix(path, ".from") && strings.HasPrefix(path, "relationships["):
		return newDiagnostic(schemaValidationSource, "FBV301", "missing required field: relationship from", path), true
	case strings.HasSuffix(path, ".to") && strings.HasPrefix(path, "relationships["):
		return newDiagnostic(schemaValidationSource, "FBV302", "missing required field: relationship to", path), true
	case path != "":
		return newDiagnostic(schemaValidationSource, "FBV001", fmt.Sprintf("cue schema validation failed at %s", path), path), true
	default:
		return domain.Diagnostic{}, false
	}
}

func isReportSectionsPath(path string) bool {
	if !strings.HasPrefix(path, "reports[") || !strings.HasSuffix(path, ".sections") {
		return false
	}
	return strings.Count(path, ".") == 1
}
