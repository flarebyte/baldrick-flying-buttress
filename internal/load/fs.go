package load

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

type FSAppLoader struct {
	ConfigPath string
}

func (l FSAppLoader) Load() (domain.RawApp, error) {
	if l.ConfigPath == "" {
		return domain.RawApp{}, fmt.Errorf("config path is required")
	}

	data, err := os.ReadFile(l.ConfigPath)
	if err != nil {
		return domain.RawApp{}, fmt.Errorf("read config %s: %w", l.ConfigPath, err)
	}

	var raw domain.RawApp
	if err := json.Unmarshal(data, &raw); err != nil {
		return domain.RawApp{}, fmt.Errorf("parse config %s: %w", l.ConfigPath, err)
	}

	raw.ConfigPath = l.ConfigPath
	return raw, nil
}
