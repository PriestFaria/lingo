package filters

import "lingo/internal/analyzer/log"


type FilterPipeline struct {
	filters []LogFilter
}

func NewFilterPipeline(filters []LogFilter) *FilterPipeline {
	return &FilterPipeline{
		filters: filters,
	}
}


func (p *FilterPipeline) Process(context *log.LogContext) []FilterIssue {
	var allIssues []FilterIssue
	for _, f := range p.filters {
		issues := f.Apply(context)
		allIssues = append(allIssues, issues...)
	}
	return allIssues
}
