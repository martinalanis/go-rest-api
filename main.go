package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Keys in Uppercase first letter so it can be exported(accesible in other packages)
// backthiks indicates struc tag, is meta data added and used to map outputs(?)
// example in response json data will be with keys ID, name, content with firs lower instead of upper as in struct
type Task struct {
	ID uint `json:"ID"`
	Name string `json:"name"`
	Content string `json:"content"`
}

type Tasks []Task

var taskList = Tasks {
	{
		ID: 1,
		Name: "Task 1",
		Content: "Content 1",
	},
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(taskList)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "Insert a valid task")
	}

	json.Unmarshal(reqBody, &newTask)

	newTask.ID = uint(len(taskList) + 1) // should be casted to uin since ID is type uin and operation returns int
	taskList = append(taskList, newTask)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newTask)
}

func showTask(w http.ResponseWriter, r *http.Request) {
	queryParams := mux.Vars(r)

	taskID, err := strconv.Atoi(queryParams["id"])

	if err != nil {
		fmt.Fprintf(w, "Invalid id")
		return
	}

	for _, task := range taskList {
		if task.ID == uint(taskID) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
			return
		}
	}

	fmt.Fprintf(w, "Task not found")
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	queryParams := mux.Vars(r)

	taskID, err := strconv.Atoi(queryParams["id"])

	if err != nil {
		fmt.Fprintf(w, "Invalid id")
		return
	}

	for index, task := range taskList {
		if task.ID == uint(taskID) {
			// concat slice from 0 to index and slice from index + 1 to the end
			taskList = append(taskList[:index], taskList[index + 1:]...)
			fmt.Fprintf(w, "The task with ID %v has been deleted", taskID)
			return
		}
	}

	fmt.Fprintf(w, "Task not found")
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	requestVars := mux.Vars(r)
	taskID, err := strconv.Atoi(requestVars["id"])

	if err != nil {
		fmt.Fprintf(w, "Invalid id")
		return
	}

	for index, task := range taskList {
		if task.ID == uint(taskID) {

			var newTask Task
			reqBody, err := ioutil.ReadAll(r.Body)

			if err != nil {
				fmt.Fprintf(w, "Insert a valid task")
				return
			}

			json.Unmarshal(reqBody, &newTask)
			newTask.ID = task.ID

			// create temp slice with task from 0 to index and new task
			tempSlice := append(taskList[:index], newTask)
			// concact temp slice with slice of task from index + 1 to the end
			taskList = append(tempSlice, taskList[index + 1:]...)
			fmt.Fprintf(w, "The task with ID %v has been updated", taskID)
			return
		}
	}

	fmt.Fprintf(w, "Task not found")
}


// w received as value and r is received as pointer
// call to this method should be: indexRoute(valueAsValue, &valueAsRefence)
// valueAsReference is passed with & at the beggining and uses * in the param definition
func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to GO-API")
}

func main() {
	// create router and restrict to not enter slash at the end like /some-url/ but instead just accept this /some-url
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks", createTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", showTask).Methods("GET")
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")
	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3000", router))

}

// Run api:
// bash: CompileDaemon
// bash: CompileDaemon -command="./go-api.exe