package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

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

	fmt.Printf("Request: %v\n", r)

	fmt.Fprintf(w, "OK")
}

func HTTPFeed(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("*** HTTPFeed ***\n")

	fmt.Printf("Request: %v\n", r)

	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			HTTP500(w, r, "Error reading request body")
			return
		}

		fmt.Printf("Body: %s\n", string(body))
	}

	sample := r.URL.Query().Get("obj")
	if sample == "" {
		HTTP500(w, r, "No sample given")
		return
	}
	fmt.Printf("Sample: %s\n", sample)

	fmt.Fprintf(w, "SAMPLE_ID_XYZ123")
}

func HTTPCheck(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("*** HTTPCheck ***\n")

	fmt.Printf("Request: %v\n", r)

	fmt.Fprintf(w, "OK")
}

func HTTPResults(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("*** HTTPResults ***\n")

	fmt.Printf("Request: %v\n", r)

	fmt.Fprintf(w, "OK")
}

func HTTP500(w http.ResponseWriter, r *http.Request, response string) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, response)
	return
}
