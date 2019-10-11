package main

import (
	"encoding/json"
	"fmt"
	"time"

	"log"
	"net/http"
	"os"

	"github.com/frozentech/go-tools/filesystem"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	hostEnv = "APP_HOST"
	portEnv = "APP_PORT"
)

// Log ...
type Log struct {
	entries []Entry
}

// Entry ...
type Entry struct {
	Time    time.Time `json:"timestamp"`
	Level   string    `json:"level"`
	Message string    `json:"message"`
}

// Entries ...
type Entries []Entry

// addEntry ...
func (l *Log) addEntry(level string, v ...interface{}) {
	l.entries = append(
		l.entries,
		Entry{
			Level:   level,
			Message: fmt.Sprint(v...),
			Time:    time.Now(),
		},
	)
}

// Dump ...
func (l *Log) Dump() {
	logrs := logrus.New()

	logrs.SetFormatter(&logrus.JSONFormatter{})
	logrs.SetOutput(os.Stdout)

	len := len(l.entries)

	var (
		entries = Entries{}
		lines   string
	)
	for i := 0; i < len; i++ {
		entries = append(entries, Entry{
			Time:    l.entries[i].Time,
			Level:   l.entries[i].Level,
			Message: l.entries[i].Message,
		})
		continue
	}

	jsonLogs, _ := json.Marshal(entries)
	lines = string(jsonLogs)
	logrs.Info(lines)

	go func() {
		// create folder for program if not exist
		p := "log/server"
		if exist, _ := filesystem.Exist(p); !exist {
			filesystem.Mkdir(p)
		}

		file, _ := os.Create(p + "/go-k8s.out.log")
		defer file.Close()
		l.entries = make([]Entry, 0)
		file.WriteString(lines)
	}()
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func check(w http.ResponseWriter, r *http.Request) {
	var (
		response = make(map[string]interface{})
	)

	response["timestamp"] = time.Now()
	response["message"] = "App is healthy!"

	l := Log{}
	l.addEntry("DEBUG", response)
	l.Dump()

	respondWithJSON(w, http.StatusOK, response)
	return
}

func getUser(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	if name == "" {
		name = "Guest"
	}

	var (
		response = make(map[string]interface{})
	)

	response["timestamp"] = time.Now()
	response["name"] = name

	l := Log{}
	l.addEntry("DEBUG", response)
	l.Dump()

	respondWithJSON(w, http.StatusOK, response)
	return
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", index)
	router.HandleFunc("/user", getUser)
	router.HandleFunc("/healthcheck", check)

	host := os.Getenv(hostEnv)
	port := os.Getenv(portEnv)

	// setting default value for host var
	if host == "" {
		host = "127.0.0.1"
	}

	// setting default value for port var
	if port == "" {
		port = "8080"
	}

	// start the listener
	fullURL := host + ":" + port
	log.Printf("Listening on `%s` ...", fullURL)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		panic("ListenAndServe error: " + err.Error())
	}

}
