package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"gocloudcamp/pkg/adapters"
	Models "gocloudcamp/pkg/models"
	"gocloudcamp/pkg/service"
	"log"
	"net/http"
	"time"
)

func StartNewHTTPServer(s interface{}, httpPort int) error {
	svc := s.(service.ConfigService)

	r := http.NewServeMux()
	r.Handle("/config", configHandler{service: svc})

	ch := make(chan error)
	go func() {
		ch <- http.ListenAndServe(fmt.Sprintf(":%d", httpPort), r)
	}()

	var e error
	select {
	case e = <-ch:
		return e
	case <-time.After(time.Second * 1):
	}
	log.Printf("HTTP server listening at %v", httpPort)
	return nil
}

type configHandler struct {
	service service.ConfigService
}

func (h configHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	svc := h.service

	switch r.Method {
	case http.MethodPost:
		req, err := adapters.DecodeSetRequest(context.TODO(), r)
		if err != nil {
			returnErrorResponse(err, w)
			return
		}
		resp, err := svc.SetConfig(context.TODO(), req)
		if err != nil {
			returnErrorResponse(err, w)
		} else {
			returnSetResponse(resp, w)
		}

	case http.MethodGet:
		req, err := adapters.DecodeGetRequest(context.TODO(), r)
		if err != nil {
			returnErrorResponse(err, w)
			return
		}
		resp, err := svc.GetConfig(context.TODO(), req)
		if err != nil {
			returnErrorResponse(err, w)
		} else {
			returnGetResponse(resp, w)
		}

	case http.MethodPut:
		req, err := adapters.DecodeGetRequest(context.TODO(), r)
		if err != nil {
			returnErrorResponse(err, w)
			return
		}
		resp, err := svc.UpdConfig(context.TODO(), req)
		if err != nil {
			returnErrorResponse(err, w)
		} else {
			resp.Version = 0
			returnSetResponse(resp, w)
		}

	case http.MethodDelete:
		req, err := adapters.DecodeGetRequest(context.TODO(), r)
		if err != nil {
			returnErrorResponse(err, w)
			return
		}
		resp, err := svc.DelConfig(context.TODO(), req)
		if err != nil {
			returnErrorResponse(err, w)
		} else {
			resp.Version = 0
			returnSetResponse(resp, w)
		}

	default:
		w.Header().Set("Allow", "GET, POST, PUT, DELETE")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func returnErrorResponse(e interface{}, w http.ResponseWriter) {
	re := e.(Models.ResponseError)
	status := http.StatusInternalServerError
	if re.Status > 0 {
		status = re.Status
	}
	w.Header().Set("Content-Type", "application/json")
	respStruct := &jsonResponse{Success: false, Message: re.ErrorDescr}
	resp, _ := json.Marshal(respStruct)
	http.Error(w, string(resp), status)
}

func returnSetResponse(e interface{}, w http.ResponseWriter) {
	re := e.(*Models.ConfigRequest)
	w.Header().Set("Content-Type", "application/json")
	respStruct := &jsonResponse{Success: true, Version: re.Version}
	resp, _ := json.Marshal(respStruct)
	fmt.Fprintln(w, string(resp))
}

func returnGetResponse(e interface{}, w http.ResponseWriter) {
	re := e.(*Models.ConfigRequest)
	w.Header().Set("Content-Type", "application/json")
	if re.Extended {
		resp, _ := json.Marshal(re)
		fmt.Fprintln(w, string(resp))
	} else {
		resp, _ := json.Marshal(re.Data)
		fmt.Fprintln(w, string(resp))
	}
}

type jsonResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Version int    `json:"version,omitempty"`
}
