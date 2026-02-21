package filters

import "github.com/PriestFaria/lingo/internal/analyzer/log"

// FilterPipeline runs a sequence of LogFilters and collects all issues.
type FilterPipeline struct {
	filters []LogFilter
}

// NewFilterPipeline creates a pipeline that will execute the given filters in order.
func NewFilterPipeline(filters []LogFilter) *FilterPipeline {
	return &FilterPipeline{
		filters: filters,
	}
}

// Process runs every filter in the pipeline against the provided LogContext
// and returns the aggregated list of issues.
func (p *FilterPipeline) Process(context *log.LogContext) []FilterIssue {
	var allIssues []FilterIssue
	for _, f := range p.filters {
		issues := f.Apply(context)
		allIssues = append(allIssues, issues...)
	}
	return allIssues
}
