package t39

func NewStraightPipeline(name string, filters ...Filter) *StraightPipeline {

	return &StraightPipeline{Name: name, Filters: &filters}
}

type StraightPipeline struct {
	Name    string
	Filters *[]Filter
}

func (sp *StraightPipeline) Process(data Request) (Response, error) {

	var ret interface{}
	var err error

	for _, filter := range *sp.Filters {
		ret, err = filter.Process(data)
		if err != nil {
			return nil, err
		}
		data = ret
	}
	return ret, nil
}
