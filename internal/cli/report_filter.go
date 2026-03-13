package cli

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/pipeline"
)

func withReportFilter(loader pipeline.AppLoader, reportIDs []string) pipeline.AppLoader {
	normalized, err := normalizeReportFilters(reportIDs)
	if err != nil {
		return pipeline.LoaderFunc(func(context.Context) (domain.RawApp, error) {
			return domain.RawApp{}, err
		})
	}
	if len(normalized) == 0 {
		return loader
	}

	return pipeline.LoaderFunc(func(ctx context.Context) (domain.RawApp, error) {
		raw, err := loader.Load(ctx)
		if err != nil {
			return domain.RawApp{}, err
		}
		filtered, err := filterRawReports(raw, normalized)
		if err != nil {
			return domain.RawApp{}, err
		}
		return filtered, nil
	})
}

func normalizeReportFilters(reportIDs []string) ([]string, error) {
	if len(reportIDs) == 0 {
		return nil, nil
	}
	seen := make(map[string]struct{}, len(reportIDs))
	normalized := make([]string, 0, len(reportIDs))
	for _, reportID := range reportIDs {
		reportID = strings.TrimSpace(reportID)
		if reportID == "" {
			return nil, fmt.Errorf("report filter must not be empty")
		}
		if _, exists := seen[reportID]; exists {
			continue
		}
		seen[reportID] = struct{}{}
		normalized = append(normalized, reportID)
	}
	slices.Sort(normalized)
	return normalized, nil
}

func filterRawReports(raw domain.RawApp, reportIDs []string) (domain.RawApp, error) {
	if len(reportIDs) == 0 {
		return raw, nil
	}

	allowed := make(map[string]struct{}, len(reportIDs))
	for _, reportID := range reportIDs {
		allowed[reportID] = struct{}{}
	}

	filtered := make([]domain.RawReport, 0, len(raw.Reports))
	available := make([]string, 0, len(raw.Reports))
	matched := make(map[string]struct{}, len(reportIDs))
	for _, report := range raw.Reports {
		reportID := rawReportFilterID(report)
		if reportID != "" {
			available = append(available, reportID)
		}
		if _, ok := allowed[reportID]; !ok {
			continue
		}
		matched[reportID] = struct{}{}
		filtered = append(filtered, report)
	}

	missing := make([]string, 0)
	for _, reportID := range reportIDs {
		if _, ok := matched[reportID]; !ok {
			missing = append(missing, reportID)
		}
	}
	if len(missing) > 0 {
		slices.Sort(available)
		return domain.RawApp{}, fmt.Errorf("unknown report filter: %s (available: %s)", strings.Join(missing, ","), strings.Join(available, ","))
	}

	raw.Reports = filtered
	return raw, nil
}

func rawReportFilterID(report domain.RawReport) string {
	reportID := domain.ReportIDFromFilepath(report.Filepath)
	if reportID != "" {
		return reportID
	}
	return strings.TrimSpace(report.Title)
}
