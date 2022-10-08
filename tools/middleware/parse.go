package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

func defaults() (them map[string]interface{}) {
	them = map[string]interface{}{
		"size":   RANGE_SIZE_DEFAULT,
		"before": "",
	}

	return
}

func limitedSize(value string, limit int) (limted int, err error) {
	limted, err = strconv.Atoi(value)

	switch {
	case err != nil || limted < 0:
		err = fmt.Errorf("Incorrect value %s", value)
		return
	case limted > limit:
		limted = limit
	}

	return
}

func PaginationParams(request *http.Request) (modified *http.Request, ok bool, code int, _ map[string]interface{}, err error) {
	var parsed map[string]interface{} = defaults()
	var key string
	var value []string
	for key, value = range request.URL.Query() {
		if len(value) == 0 {
			continue
		}

		switch key {
		case "size":
			if parsed[key], err = limitedSize(value[0], RANGE_SIZE_LIMIT); err != nil {
				err = nil
				code = 400
				return
			}
		case "before":
			parsed[key] = value[0]
		}
	}

	ok = true
	modified = request.WithContext(context.WithValue(
		request.Context(),
		"query",
		parsed,
	))

	return
}

func ParseMultipart(request *http.Request) (_ *http.Request, ok bool, code int, r_map map[string]interface{}, err error) {
	if err = request.ParseMultipartForm(MULTIPART_MEM_MAX); err != nil {
		err = nil
		code = 400
		return
	}

	ok = true
	return
}
