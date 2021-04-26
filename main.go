package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jasonlvhit/gocron"
)

type ResponeObject map[string]interface{}

func myTask() {
	layout := "2006-01-02 15:04:05"
	fmt.Println(time.Now().Format(layout), " : This task will run periodically")
}
func executeCronJob() {
	gocron.Every(10).Seconds().Do(myTask)
	<-gocron.Start()
}

func main() {

	go executeCronJob()

	r := mux.NewRouter()
	r.StrictSlash(true)

	login := r.PathPrefix("/api").Subrouter()
	login.Use(commonMiddleware)

	// protected.Use(auth.JwtVerify)
	login.HandleFunc("/helloWorld", HelloWorld).Methods("GET")
	log.Fatal(http.ListenAndServe("localhost:8088", commonMiddleware(r)))
	// createNewBMBUser()
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	uname := GetUrlStringParam(r, "uname")
	json.NewEncoder(w).Encode(ResponeObject{"err": 0, "msg": "Welcome to my system", "dt": uname})
}

func GetUrlStringParam(r *http.Request, paramKey string) string {
	var paramValue string

	keys, ok := r.URL.Query()[paramKey]

	if !ok || len(keys[0]) < 1 {
		paramValue = ""
	} else {
		paramValue = keys[0]
	}

	return paramValue
}

func commonMiddleware(h http.Handler) http.Handler {
	println("server started with port 8088")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		/*for k, v := range r.Header {
			fmt.Fprintf(w, "Header field %q, Value %q\n", k, v)
		}*/

		w.Header().Set("Access-Control-Allow-Origin", origin)
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST, PUT, DELETE, PATCH, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token, X-Auth-Token, Authorization")
			return
		} else {

			h.ServeHTTP(w, r)
		}
	})
}
