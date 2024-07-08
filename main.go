package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// model for course -file
type Course struct {
	CourseID    string  `json:"courseid"`
	CourseName  string  `json:"coursename"`
	CoursePrice int     `json:"price"`
	Author      *Author `json:"author"`
}

// model for author -file
type Author struct {
	FullName string `json:"fullname"`
	Website  string `json:"website"`
}

// fake db
var courses []Course

// middleware/helper -file
func (c *Course) IsEmpty() bool {
	// return c.CourseID == "" && c.CourseName == ""
	return c.CourseName == ""
}

func main() {
	//seeding data into demo db
	courses = append(courses, Course{CourseID: "1", CourseName: "Go	Language", CoursePrice: 199, Author: &Author{FullName: "John Doe", Website: "https://www.johndoe.com"}})
	courses = append(courses, Course{CourseID: "2", CourseName: "Javascript", CoursePrice: 499, Author: &Author{FullName: "Subash Choudhary", Website: "https://www.chaiwithcode.com"}})

	fmt.Println("API - Learn Code Online")
	r := mux.NewRouter()
	r.HandleFunc("/", serveHome).Methods("GET")
	r.HandleFunc("/courses", getAllCourses).Methods("GET")
	r.HandleFunc("/course/{id}", getOneCourse).Methods("GET")
	r.HandleFunc("/course/create", createOneCourse).Methods("POST")
	r.HandleFunc("/course/update/{id}", updateOneCourse).Methods("PUT")
	r.HandleFunc("/course/delete/{id}", deleteOneCourse).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3500", r))
}

//controllers -file

// serve home route
func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Welcome to the home page of courses application!</h1>"))
}

func getAllCourses(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get all COurses")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func getOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get one course")
	w.Header().Set("Content-Type", "application/json")

	//grab id from request
	params := mux.Vars(r)
	fmt.Println(params)

	//loop through courses, find matching id and return the response
	for _, course := range courses {
		if course.CourseID == params["id"] {
			json.NewEncoder(w).Encode(course)
			return
		}
	}
	message := "No course found with id: " + params["id"]
	json.NewEncoder(w).Encode(message)
	return
}

func createOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create one course")
	w.Header().Set("Content-Type", "application/json")

	//what if body is empty
	if r.Body == nil {
		json.NewEncoder(w).Encode("please send some data")
		return
	}
	// what about this type of data - {}

	var course Course
	_ = json.NewDecoder(r.Body).Decode(&course)
	if course.IsEmpty() {
		json.NewEncoder(w).Encode("Json data is empty")
		return
	}

	//generate unique id, convert it to string
	//append this new course to courses
	rand.Seed(time.Now().UnixNano())
	course.CourseID = strconv.Itoa(rand.Intn(100))
	courses = append(courses, course)
	json.NewEncoder(w).Encode(course)
	return
}

func updateOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update one course")
	w.Header().Set("Content-Type", "application/json")

	//first - grab id from req
	params := mux.Vars(r)

	//loop,id, remove, add with my ID
	for index, item := range courses {
		if item.CourseID == params["id"] {
			courses = append(courses[:index], courses[index+1:]...)

			var course Course
			_ = json.NewDecoder(r.Body).Decode(&course)
			course.CourseID = params["id"]
			courses = append(courses, course)
			json.NewEncoder(w).Encode(course)
			return
		}
	}
	//send a response when id is not found
	message := "No course found with id: " + params["id"]
	json.NewEncoder(w).Encode(message)
	return
}

func deleteOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update one course")
	w.Header().Set("Content-Type", "application/json")

	//first - grab id from req
	params := mux.Vars(r)

	//loop,id, remove
	for index, item := range courses {
		if item.CourseID == params["id"] {
			courses = append(courses[:index], courses[index+1:]...)
			json.NewEncoder(w).Encode("deleted successfully")
			return
		}
	}
	message := "No course found with id: " + params["id"]
	json.NewEncoder(w).Encode(message)
	return
}
