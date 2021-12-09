package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Course struct {
	CourseId    string  `json:"courseId"`
	CourseName  string  `json:"courseName"`
	CoursePrice string  `json:"coursePrice"`
	Author      *Author `json:"author"`
}

type Author struct {
	FullName string `json:"fullName"`
	WebSite  string `json:"webSite"`
}

var courses []Course

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", serveHome).Methods("GET")
	r.HandleFunc("/create", createOneCourse).Methods("POST")
	r.HandleFunc("/get", getAllCourses).Methods("GET")
	r.HandleFunc("/getone", getOneCourse).Methods("GET")
	r.HandleFunc("/update", updateOneCourse).Methods("POST")
	r.HandleFunc("/delete", deleteOneCourse).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":3000", r))
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello welcome to Golang"))
}

func (c *Course) isEmpty() bool {
	return c.CourseName == ""
}

func getAllCourses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func getOneCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//get id
	params := mux.Vars(r)

	//loop through the courses

	for _, course := range courses {
		if course.CourseId == params["id"] {
			json.NewEncoder(w).Encode(course)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode("No record found")
}

func createOneCourse(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("requestBody cannot be Empty")
		return
	}

	var course Course
	json.NewDecoder(r.Body).Decode(&course)
	if course.isEmpty() {
		json.NewEncoder(w).Encode("courseName cannot be Empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rand.Seed(time.Now().UnixNano())
	course.CourseId = strconv.Itoa(rand.Intn(100))
	courses = append(courses, course)
	json.NewEncoder(w).Encode(course)
}

func updateOneCourse(w http.ResponseWriter, r *http.Request) {
	//find the id
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, course := range courses {
		if params["id"] == course.CourseId {
			courses = append(courses[0:index], courses[index+1:]...)
			json.NewDecoder(r.Body).Decode(&course)
			course.CourseId = params["id"]
			courses = append(courses, course)
			json.NewEncoder(w).Encode(course)
			return
		}
	}
	json.NewEncoder(w).Encode("Not found")
}

func deleteOneCourse(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	for index, course := range courses {
		if params["id"] == course.CourseId {
			courses = append(courses[0:index], courses[index+1:]...)
			json.NewEncoder(w).Encode("delete completed")
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode("record Not found")
}
