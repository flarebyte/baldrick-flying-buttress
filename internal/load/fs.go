// purpose: Loads and resolves external config inputs in fs.go into raw app data for validation.
// responsibilities: resolve config paths and formats; read filesystem/CUE inputs; normalize load-time source metadata
// architecture_notes: Load logic isolates format/path resolution so validation and CLI code remain agnostic to input source details.
package load

import (
	"context"
	"fmt"
	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

type FSAppLoader struct {
	ConfigPath string
}

func (l FSAppLoader) Load(ctx context.Context) (domain.RawApp, error) {
	if l.ConfigPath == "" {
		return domain.RawApp{}, fmt.Errorf("config path is required")
	}
	if err := ctx.Err(); err != nil {
		return domain.RawApp{}, err
	}
	value, resolvedPath, err := CompileConfigValue(l.ConfigPath)
	if err != nil {
		return domain.RawApp{}, err
	}

	var raw domain.RawApp
	if err := value.Decode(&raw); err != nil {
		return domain.RawApp{}, fmt.Errorf("parse config %s: %w", resolvedPath, err)
	}
	if err := ctx.Err(); err != nil {
		return domain.RawApp{}, err
	}

	raw.ConfigPath = resolvedPath
	return raw, nil
}
