package parser

import (
	"errors"
	"regexp"
	"strconv"
)

type Metrics struct {
	Valuecounter int64
	Valuegauge   float64
	Nametype     string
}

var ErrValue = errors.New("value is not correct")
var ErrType = errors.New("type metric is not correct")
var ErrNameMetric = errors.New("name metric is not exist")

func Parseuri(v string) (Metrics, error) {

	re, _ := regexp.Compile(`/*[\s\S]+?/`)
	resregex := re.FindAllString(v+"/", -1)
	if len(resregex) < 4 {
		return Metrics{}, ErrNameMetric
	}
	if resregex[0] == "/update/" && (resregex[1] == "counter/" || resregex[1] == "gauge/") {
		if resregex[1] == "counter/" {
			valueint64, err := strconv.ParseInt(resregex[3][:len(resregex[3])-1], 10, 64)
			if err != nil {

			}
			if valueint64 < 0 {
				panic(errors.New("gauge cannot decrease in value"))
			}
			return Metrics{Valuecounter: valueint64, Nametype: resregex[2][:len(resregex[2])-1]}, nil
		} else {

			valuefloat64, err := strconv.ParseFloat(resregex[3][:len(resregex[3])-1], 64)
			if err != nil {
				return Metrics{}, ErrValue
			}
			if valuefloat64 < 0 {
				panic(errors.New("counter cannot decrease in value"))
			}
			return Metrics{Valuegauge: valuefloat64, Nametype: resregex[2][:len(resregex[2])-1]}, nil
		}

	}

	return Metrics{}, nil
}
