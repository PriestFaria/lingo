package filters

import (
	"testing"
)

func TestFilterPipeline_RunsAllFilters(t *testing.T) {
	pipeline := NewFilterPipeline([]LogFilter{
		&FirstLetterFilter{},
		&EnglishFilter{},
		&EmojiStrictFilter{},
		&SecurityFilter{},
	})

	ctx := makeCtx(makeParts("Password: ", true))
	issues := pipeline.Process(ctx)
	if len(issues) != 2 {
		t.Errorf("got %d issues, want 2 (FirstLetter + Security)", len(issues))
	}
}

func TestFilterPipeline_Empty(t *testing.T) {
	pipeline := NewFilterPipeline([]LogFilter{})
	ctx := makeCtx(makeParts("Starting server", true))
	issues := pipeline.Process(ctx)
	if len(issues) != 0 {
		t.Errorf("got %d issues, want 0 for empty pipeline", len(issues))
	}
}

func TestFilterPipeline_NoIssues(t *testing.T) {
	pipeline := NewFilterPipeline([]LogFilter{
		&FirstLetterFilter{},
		&EnglishFilter{},
		&EmojiStrictFilter{},
		&SecurityFilter{},
	})

	ctx := makeCtx(makeParts("server started successfully", true))
	issues := pipeline.Process(ctx)
	if len(issues) != 0 {
		t.Errorf("got %d issues, want 0 for clean message", len(issues))
	}
}

func TestFilterPipeline_AggregatesFromMultipleFilters(t *testing.T) {
	pipeline := NewFilterPipeline([]LogFilter{
		&FirstLetterFilter{},
		&EnglishFilter{},
	})

	ctx := makeCtx(makeParts("Запуск", true))
	issues := pipeline.Process(ctx)
	if len(issues) != 2 {
		t.Errorf("got %d issues, want 2 (FirstLetter + English)", len(issues))
	}
}
