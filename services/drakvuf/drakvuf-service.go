package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type RespStatus struct {
	Degraded  bool
	Error     string
	FreeSlots int
}

type RespNewTask struct {
	Error  string
	TaskID string
}

type RespCheckTask struct {
	Error string
	Done  bool
}

type RespTaskResults struct {
	Error   string
	Results interface{}
}

func main() {

	HTTPBinding := ":8090"

	// routing
	r := http.NewServeMux()

	r.HandleFunc("/status/", HTTPStatus)
	r.HandleFunc("/feed/", HTTPFeed)
	r.HandleFunc("/check/", HTTPCheck)
	r.HandleFunc("/results/", HTTPResults)

	srv := &http.Server{
		Handler:      r,
		Addr:         HTTPBinding,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}

func HTTPStatus(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("*** HTTPStatus ***\n")

	resp := &RespStatus{
		Degraded:  false,
		Error:     "",
		FreeSlots: 0,
	}

	fmt.Printf("Request: %v\n", r)
	fmt.Printf("Method: %s\n", r.Method)

	json.NewEncoder(w).Encode(resp)
}

func HTTPFeed(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("*** HTTPFeed ***\n")

	resp := &RespNewTask{
		Error:  "",
		TaskID: "",
	}

	fmt.Printf("Request: %v\n", r)
	fmt.Printf("Method: %s\n", r.Method)

	sample := r.URL.Query().Get("obj")
	if sample == "" {
		resp.Error = "No sample given"
		HTTP500(w, r, resp)
		return
	}
	fmt.Printf("Sample: %s\n", sample)

	resp.TaskID = "SAMPLE_ID_XYZ123"

	json.NewEncoder(w).Encode(resp)
}

func HTTPCheck(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("*** HTTPCheck ***\n")

	resp := &RespCheckTask{
		Error: "",
		Done:  false,
	}

	fmt.Printf("Request: %v\n", r)
	fmt.Printf("Method: %s\n", r.Method)

	taskIDstr := r.URL.Query().Get("taskid")
	if taskIDstr == "" {
		resp.Error = "No taskID given"
		HTTP500(w, r, resp)
		return
	}
	fmt.Printf("taskID: %s\n", taskIDstr)

	json.NewEncoder(w).Encode(resp)
}

func HTTPResults(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("*** HTTPResults ***\n")

	resp := &RespTaskResults{
		Error: "",
	}

	fmt.Printf("Request: %v\n", r)
	fmt.Printf("Method: %s\n", r.Method)

	taskIDstr := r.URL.Query().Get("taskid")
	if taskIDstr == "" {
		resp.Error = "No taskID given"
		HTTP500(w, r, resp)
		return
	}
	fmt.Printf("taskID: %s\n", taskIDstr)

	json.NewEncoder(w).Encode(resp)
}

func HTTP500(w http.ResponseWriter, r *http.Request, response interface{}) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(response)
	return
}
