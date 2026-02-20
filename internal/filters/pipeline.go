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


func (p *FilterPipeline) Process(context *log.LogContext){
	for _, f := range p.filters {
		f.Apply(context) 
	}
}
