package filters


type FilterPipeline struct {
	filters []LogFilter
}


func (p *FilterPipeline) Process(context *LogContext){
	for _, f := range p.filters {
		f.Apply(context) 
	}
}