package conversion

import (
	"github.com/buonotti/apisense/v2/validation/pipeline"
	"github.com/goccy/go-json"
)

// Json returns a new jsonConverter
func Json() Converter {
	return jsonConverter{}
}

type jsonConverter struct{}

func (jsonConverter) Convert(reports ...pipeline.Report) ([]byte, error) {
	if len(reports) == 1 {
		return json.Marshal(reports[0])
	} else {
		return json.Marshal(reports)
	}
}
