package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

func TasksFlatHandle(w http.ResponseWriter, r *http.Request) {
	HandleFunc(w, r, GetTasksFlat, nil)
}

func TasksHierarchyHandle(w http.ResponseWriter, r *http.Request) {
	HandleFunc(w, r, GetTasksHierarchy, nil)
}

func SolutionHandle(w http.ResponseWriter, r *http.Request) {
	HandleFunc(w, r, nil, PostSolution)
}

func RegistrationHandle(w http.ResponseWriter, r *http.Request) {
	HandleFunc(w, r, GetRegistration, PostRegistration)
}

func VerifyHandle(w http.ResponseWriter, r *http.Request) {
	HandleFunc(w, r, GetVerify, nil)
}

func LoginHandle(w http.ResponseWriter, r *http.Request) {
	HandleFunc(w, r, GetLogin, nil)
}

func TemplateHandle(w http.ResponseWriter, r *http.Request) {
	HandleFunc(w, r, GetTemplate, nil)
}

func LanguagesHandle(w http.ResponseWriter, r *http.Request) {
	HandleFunc(w, r, GetLanguages, nil)
}

func RestoreHandle(w http.ResponseWriter, r *http.Request) {
	HandleFunc(w, r, GetRestore, PostRestore)
}

type WebMethodFunc func(r *http.Request, resp *map[string]interface{}) WebError

func HandleFunc(w http.ResponseWriter, r *http.Request, get WebMethodFunc, post WebMethodFunc) {
	var web_err WebError
	web_err = MethodNotSupported
	resp := map[string]interface{}{}
	defer RecoverRequest(w)
	switch r.Method {
	case "GET":
		if get == nil {
			web_err = MethodNotSupported
		} else {
			web_err = get(r, &resp)
		}
	case "POST":
		if post == nil {
			web_err = MethodNotSupported
		} else {
			web_err = post(r, &resp)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if web_err != NoError {
		log.Printf("Failed user request, error code: %d", web_err)
		resp["error"] = web_err
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	jsonResp, err := json.Marshal(resp)
	Err(err)
	w.Write(jsonResp)
}

func RecoverRequest(w http.ResponseWriter) {
	r := recover()
	if r != nil {
		debug.PrintStack()
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		log.Printf("[INTERNAL ERROR]: %s", r)
		response := fmt.Sprintf("{\"error\": \"%d\"}", Internal)
		w.Write([]byte(response))
	}
}
