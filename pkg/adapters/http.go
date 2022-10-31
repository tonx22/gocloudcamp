package adapters

import (
	"context"
	"fmt"
	"github.com/tidwall/gjson"
	Models "gocloudcamp/pkg/models"
	"io"
	"net/http"
	"strconv"
)

func DecodeSetRequest(ctx context.Context, r *http.Request) (*Models.ConfigRequest, error) {
	var req Models.ConfigRequest
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, Models.ResponseError{ErrorDescr: "Reading input failure"}
	}

	json := string(b)
	fmt.Println(json)
	if !gjson.Valid(json) {
		return nil, Models.ResponseError{ErrorDescr: "Invalid input json", Status: http.StatusBadRequest}
	}

	service := gjson.Get(json, "service")
	if !service.Exists() {
		return nil, Models.ResponseError{ErrorDescr: "Invalid json: service field missing", Status: http.StatusBadRequest}
	} else if len(service.String()) == 0 {
		return nil, Models.ResponseError{ErrorDescr: "Invalid json: service field is empty", Status: http.StatusBadRequest}
	} else {
		req.Service = service.String()
	}

	data := gjson.Get(json, "data")
	if !data.Exists() {
		return nil, Models.ResponseError{ErrorDescr: "Invalid json: data field missing", Status: http.StatusBadRequest}
	} else if len(data.String()) == 0 {
		return nil, Models.ResponseError{ErrorDescr: "Invalid json: data field is empty", Status: http.StatusBadRequest}
	}

	req.Data = make(map[string]interface{})
	for _, w := range data.Array() {
		d, ok := w.Value().(map[string]interface{})
		if !ok {
			return nil, Models.ResponseError{ErrorDescr: "Invalid json: data array is invalid", Status: http.StatusBadRequest}
		}
		for k, v := range d {
			req.Data[k] = v
		}
	}
	return &req, nil
}

func DecodeGetRequest(_ context.Context, r *http.Request) (*Models.ConfigRequest, error) {
	var req Models.ConfigRequest

	service := r.URL.Query().Get("service")
	if len(service) == 0 {
		return nil, Models.ResponseError{ErrorDescr: "service parameter must be specified", Status: http.StatusBadRequest}
	}
	req.Service = service

	v := r.URL.Query().Get("version")
	if len(v) > 0 {
		version, err := strconv.Atoi(v)
		if err != nil {
			return nil, Models.ResponseError{ErrorDescr: "version parameter incorrect, must be a number", Status: http.StatusBadRequest}
		}
		req.Version = version
	}

	used := r.URL.Query().Get("used")
	if len(used) > 0 {
		switch used {
		case "true", "false":
			req.Used, _ = strconv.ParseBool(used)
		default:
			return nil, Models.ResponseError{ErrorDescr: "used parameter incorrect, must be a true or false", Status: http.StatusBadRequest}
		}
	}

	if r.URL.Query().Get("extended") == "true" {
		req.Extended = true
	}
	return &req, nil
}
