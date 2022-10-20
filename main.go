package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

func handleRequests() {
	route := mux.NewRouter()

	// router path folder untuk public
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	// routing
	route.HandleFunc("/", Home).Methods("GET")
	route.HandleFunc("/contact", Contact).Methods("GET")
	route.HandleFunc("/formProject", formProject).Methods("GET")
	route.HandleFunc("/detailProject/{index}", DetailProject).Methods("GET")
	route.HandleFunc("/addProject", addProject).Methods("POST")
	route.HandleFunc("/deleteProject/{index}", deteleProject).Methods("GET")
	route.HandleFunc("/formEditProject/{index}", formEditProject).Methods("GET")
	route.HandleFunc("/editProject/{id}", editProject).Methods("POST")


	fmt.Println("Go Running on Port 5000")
	http.ListenAndServe(":5000", route)
}

type Project struct {
	Name 			string
	StartDate string
	EndDate 	string
	Duration	string
	Desc 			string
	Id				int
	Tech			[]string
}

var data = []Project{}


func Home(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Contect-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	
	card := map[string]interface{}{
		"Add": data,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, card)
}

func Contact(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Contect-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/contact.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func formProject(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Contect-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/addProject.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func DetailProject(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Contect-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/detailProject.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	
	var Detail = Project{}

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	for i, data := range data {
		if index == i {
			Detail = Project{
				Name: data.Name,
				StartDate: data.StartDate,
				EndDate: data.EndDate,
				Desc: data.Desc,
				Duration: data.Duration,
			}
		}
	}

	data := map[string]interface{}{
		"Details": Detail,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}


func addProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var name = r.PostForm.Get("inputTitle")
	var startdate = r.PostForm.Get("inputStartDate")
	var enddate = r.PostForm.Get("inputEndDate")
	var desc = r.PostForm.Get("inputDesc")
	var tech []string
	tech = r.Form["inputTech"]

	layout := "2006-01-02"
	dateStart, _ := time.Parse(layout, startdate)
	dateEnd, _ := time.Parse(layout, enddate)
	
	hours := dateEnd.Sub(dateStart).Hours()
	daysInHours := hours / 24
	monthInDay := daysInHours / 30
	yearInMonth := monthInDay / 12

	var duration string
	var month, _ float64 = math.Modf(monthInDay)
	var year, _ float64 = math.Modf(yearInMonth)

	if year > 0 {
		duration = strconv.FormatFloat(year, 'f', 0, 64) + "Years"
	} else if month > 0 {
		duration = strconv.FormatFloat(month, 'f', 0, 64) + "Months"
	}	else if daysInHours > 0 {
		duration = strconv.FormatFloat(daysInHours, 'f', 0, 64) + "Days"
	} else if hours > 0 {
		duration = strconv.FormatFloat(hours, 'f', 0, 64) + "Hours"
	} else {
		duration = "0 Days"
	}

	var newData = Project {
		Name:				name,
		StartDate: 	startdate,
		EndDate: 		enddate,
		Duration:		duration,
		Desc: 			desc,
		Id:					len(data),
		Tech:				tech,
	}

	data = append(data, newData)

	http.Redirect(w,r, "/", http.StatusMovedPermanently)
}

func deteleProject(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(mux.Vars(r)["index"])
	data = append(data[:index], data[index+1:]...)
	http.Redirect(w, r, "/", http.StatusFound)
}

func editProject(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	Name := r.PostForm.Get("Name")
	StartDate := r.PostForm.Get("StartDate")
	EndDate := r.PostForm.Get("EndDate")
	Desc := r.PostForm.Get("Desc")

	editData := Project{
		Name: 			Name,
		StartDate: 	StartDate,
		EndDate: 		EndDate,
		Desc: 			Desc,
		Id: 				id,	
	}

	data[id] = editData

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func formEditProject(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/editMyProject.html")

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	ProjectEdit := Project{}

	for i, data := range data {
		if index == i {
			ProjectEdit = Project{
				Name: 			data.Name,
				StartDate: 	data.StartDate,
				EndDate: 		data.EndDate,
				Desc: 			data.Desc,
				Id: 				data.Id,
			}
		}
	}

	response := map[string]interface{}{
		"Project": ProjectEdit,
	}

		if err == nil {
		tmpl.Execute(w, response)
	} else {
		panic(err)
	}
}

func main() {
	handleRequests() 
}