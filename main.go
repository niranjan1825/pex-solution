package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
)

var c = cache.New(5*time.Minute, 10*time.Minute)

func init() {
	c.Set("current", 0, cache.DefaultExpiration)
}
func main() {
	router := mux.NewRouter()
	router.Use(HandleServerCrash)

	router.HandleFunc("/current", GetCurrentValue).Methods(http.MethodGet)
	router.HandleFunc("/previous", GetPreviousValue).Methods(http.MethodGet)
	router.HandleFunc("/next", GetNextValue).Methods(http.MethodGet)
	err := http.ListenAndServe(":8443", router)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}

func HandleServerCrash(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Println("crash/runtime error occured and we are recovering from panic -------->>", string(debug.Stack()))
				sendResponse(w, "we are recovering from panics/crash", http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	})
}

func GetCurrentValue(w http.ResponseWriter, r *http.Request) {
	val, ok := c.Get("current")
	if !ok {
		log.Println("no current value found in cache")
		sendResponse(w, "no current value found in cache", http.StatusInternalServerError)
		return
	}
	sendResponse(w, val, http.StatusOK)
}

func GetPreviousValue(w http.ResponseWriter, r *http.Request) {
	previous, ok := c.Get("previous")
	if !ok {
		log.Println("no previous value found in cache")
		sendResponse(w, "no previous value found since current value is 0", http.StatusOK)
		return
	}
	sendResponse(w, previous, http.StatusOK)
}

func GetNextValue(w http.ResponseWriter, r *http.Request) {
	var next int
	var previous interface{}
	current, ok := c.Get("current")
	if !ok {
		log.Println("no current value found in cache")
		sendResponse(w, "no current value found", http.StatusInternalServerError)
		return
	}
	if current.(int) == 0 {
		next = 1
		current, previous = next, current.(int)
		c.Set("current", current, cache.DefaultExpiration)
		c.Set("previous", previous, cache.DefaultExpiration)
	} else {
		previous, ok := c.Get("previous")
		if !ok {
			log.Println("no previous value found in cache")
			sendResponse(w, "no previous value found in cache to calculate next value", http.StatusInternalServerError)
			return
		}
		next = current.(int) + previous.(int)
		current, previous = next, current
		c.Set("current", current, cache.DefaultExpiration)
		c.Set("previous", previous, cache.DefaultExpiration)
	}
	fmt.Println("next:", next)
	sendResponse(w, next, http.StatusOK)

}

func sendResponse(w http.ResponseWriter, data interface{}, status int) error {
	result, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("error occured in marshalling response")
	}
	w.WriteHeader(status)
	w.Write(result)
	return nil
}
