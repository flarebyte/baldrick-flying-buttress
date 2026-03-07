package load

import (
	"context"
	"fmt"
	"os"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/safety"
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
	info, err := os.Stat(l.ConfigPath)
	if err != nil {
		return domain.RawApp{}, fmt.Errorf("stat config %s: %w", l.ConfigPath, err)
	}
	if err := safety.CheckConfigFileSize(info.Size()); err != nil {
		return domain.RawApp{}, err
	}

	data, err := os.ReadFile(l.ConfigPath)
	if err != nil {
		return domain.RawApp{}, fmt.Errorf("read config %s: %w", l.ConfigPath, err)
	}
	if err := ctx.Err(); err != nil {
		return domain.RawApp{}, err
	}

	cueCtx := cuecontext.New()
	value := cueCtx.CompileBytes(data, cue.Filename(l.ConfigPath))
	if err := value.Err(); err != nil {
		return domain.RawApp{}, fmt.Errorf("parse config %s: %w", l.ConfigPath, err)
	}

	var raw domain.RawApp
	if err := value.Decode(&raw); err != nil {
		return domain.RawApp{}, fmt.Errorf("parse config %s: %w", l.ConfigPath, err)
	}
	if err := ctx.Err(); err != nil {
		return domain.RawApp{}, err
	}

	raw.ConfigPath = l.ConfigPath
	return raw, nil
}
