package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/saimoon/Holmes-Totem-Dynamic/services/drakvuf/drakvuf"
)

const (
	STORAGE_SAMPLE_DIR string = "/tmp/"
)

type Config struct {
	HTTPBinding   string
	VerifySSL     bool
	IncomingDir   string
	ProcessingDir string
	FinishedDir   string
	MaxPending    int
	MaxAPICalls   int
	LogFile       string
	LogLevel      string
}

type Ctx struct {
	Config  *Config
	Drakvuf *drakvuf.Drakvuf
}

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
	Error  string
	Result interface{}
}

var (
	ctx *Ctx
)

func main() {
	// prepare context
	ctx = &Ctx{
		Config: &Config{},
	}

	// parse configuration file
	cFile, err := os.Open("./service.conf")
	if err != nil {
		panic(err.Error())
	}

	decoder := json.NewDecoder(cFile)
	err = decoder.Decode(ctx.Config)
	if err != nil {
		panic(err.Error())
	}

	// create Drakvuf object
	drakvuf, err := drakvuf.New(ctx.Config.IncomingDir, ctx.Config.ProcessingDir, ctx.Config.FinishedDir)
	if err != nil {
		panic(err.Error())
	}
	ctx.Drakvuf = drakvuf

	// prepare routing
	r := http.NewServeMux()
	r.HandleFunc("/status/", HTTPStatus)
	r.HandleFunc("/feed/", HTTPFeed)
	r.HandleFunc("/check/", HTTPCheck)
	r.HandleFunc("/results/", HTTPResults)

	// listening server
	srv := &http.Server{
		Handler:      r,
		Addr:         ctx.Config.HTTPBinding,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func HTTPStatus(w http.ResponseWriter, r *http.Request) {
	resp := &RespStatus{
		Degraded:  false,
		Error:     "",
		FreeSlots: 0,
	}

	s, err := ctx.Drakvuf.GetStatus()
	if err != nil {
		resp.Degraded = true
		resp.Error = err.Error()
		HTTP500(w, r, resp)
		return
	}

	resp.FreeSlots = ctx.Config.MaxPending - s.PendingTaskNum

	log.Println("HTTPStatus FreeSlots: ", resp.FreeSlots)

	json.NewEncoder(w).Encode(resp)
}

func HTTPFeed(w http.ResponseWriter, r *http.Request) {
	resp := &RespNewTask{
		Error:  "",
		TaskID: "",
	}

	sample := r.URL.Query().Get("obj")
	if sample == "" {
		resp.Error = "No sample given"
		HTTP500(w, r, resp)
		return
	}

	sampleBytes, err := ioutil.ReadFile(STORAGE_SAMPLE_DIR + sample)
	if err != nil {
		resp.Error = err.Error()
		HTTP500(w, r, resp)
		return
	}

	log.Println("DEBUG: Drakvuf.NewTask() sample=", sample)
	taskID, err := ctx.Drakvuf.NewTask(sampleBytes, sample)
	if err != nil {
		resp.Error = err.Error()
		HTTP500(w, r, resp)
		return
	}
	log.Println("HTTPFeed TaskID: ", taskID)

	resp.TaskID = taskID

	json.NewEncoder(w).Encode(resp)
}

func HTTPCheck(w http.ResponseWriter, r *http.Request) {
	resp := &RespCheckTask{
		Error: "",
		Done:  false,
	}

	taskID := r.URL.Query().Get("taskid")
	if taskID == "" {
		resp.Error = "No taskID given"
		HTTP500(w, r, resp)
		return
	}

	s, err := ctx.Drakvuf.TaskStatus(taskID)
	if err != nil {
		resp.Error = err.Error()
		HTTP500(w, r, resp)
		return
	}

	resp.Done = (s == 1)

	log.Println("HTTPCheck Done: ", s)

	json.NewEncoder(w).Encode(resp)
}

func HTTPResults(w http.ResponseWriter, r *http.Request) {
	resp := &RespTaskResults{
		Error: "",
	}

	taskID := r.URL.Query().Get("taskid")
	if taskID == "" {
		resp.Error = "No taskID given"
		HTTP500(w, r, resp)
		return
	}

	report, err := ctx.Drakvuf.TaskReport(taskID)
	if err != nil {
		resp.Error = err.Error()
		HTTP500(w, r, resp)
		return
	}

	resp.Result = report.Result
	log.Println("DEBUG: HTTPResults Result = ", report.Result)

	if err = ctx.Drakvuf.DeleteTask(taskID); err != nil {
		log.Println("Cleaning drakvuf up failed for task", taskID, err.Error())
	}

	log.Println("HTTPResults Result = ", resp.Result)

	json.NewEncoder(w).Encode(resp)
}

func HTTP500(w http.ResponseWriter, r *http.Request, response interface{}) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(response)
	return
}
