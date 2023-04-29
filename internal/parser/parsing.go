package parser

import (
	"errors"
	"strconv"
	"strings"
)

type Metrics struct {
	Valuecounter int64
	Valuegauge   float64
	Nametype     string
	Type         string
}

var ErrValue = errors.New("value is not correct")
var ErrType = errors.New("type metric is not correct")
var ErrNameMetric = errors.New("name metric is not exist")

func Parseuri(v string) (Metrics, error) {
	ressplit := strings.Split(v, "/")

	if len(ressplit) < 5 {
		return Metrics{}, ErrNameMetric
	}
	if ressplit[1] == "update" && (ressplit[2] == "counter" || ressplit[2] == "gauge") {
		if ressplit[2] == "counter" {
			valueint64, err := strconv.ParseInt(ressplit[4], 10, 64)
			if err != nil {
				return Metrics{}, ErrValue
			}
			if valueint64 < 0 {
				panic(errors.New("gauge cannot decrease in value"))
			}
			return Metrics{Valuecounter: valueint64, Nametype: ressplit[3], Type: ressplit[2]}, nil
		} else {

			valuefloat64, err := strconv.ParseFloat(ressplit[4], 64)
			if err != nil {
				return Metrics{}, ErrValue
			}
			if valuefloat64 < 0 {
				panic(errors.New("counter cannot decrease in value"))
			}
			return Metrics{Valuegauge: valuefloat64, Nametype: ressplit[3], Type: ressplit[2]}, nil
		}

	}

	return Metrics{}, ErrType
}
