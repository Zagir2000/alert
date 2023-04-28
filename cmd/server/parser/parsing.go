package parser

import (
	"net/http"
	"regexp"
)

func parseuri(res http.ResponseWriter, req *http.Request, v string) (string, string) {
	re, _ := regexp.Compile(`/*[\s\S]+?/`)
	resregex := re.FindAllString(v+"/", -1)
	if len(resregex) < 4 {
		res.WriteHeader(http.StatusNotFound)
		return "", ""
	}
	if resregex[0] == "/update/" && (resregex[1] == "counter/" || resregex[1] == "gauge/") {
		return resregex[1][:len(resregex[1])-1], resregex[3][:len(resregex[3])-1]
	}
	res.WriteHeader(http.StatusBadRequest)
	return "", ""
}
