package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"personal-web/connection"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter()

	connection.DatabaseConnect()

	// Route untuk menginisialisasi folder public
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/project", project).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/project-detail/{id}", projectDetail).Methods("GET")
	route.HandleFunc("/add-project", addProjects).Methods("POST")
	route.HandleFunc("/delete-project/{index}", deleteProjects).Methods("GET")
	route.HandleFunc("/edit-project/{in}", editProject).Methods("GET")
	route.HandleFunc("/edit-project/{in}", formEditProject).Methods("POST")

	fmt.Println("Server sedang berjalan pada port 5000")
	http.ListenAndServe("localhost:5000", route)
}

type Project struct {
	Id              int
	Title           string
	Description     string
	StartDate       time.Time
	EndDate         time.Time
	FormatStartDate string
	FormatEndDate   string
	NodeJs          string
	ReactJs         string
	JavaScript      string
	TypeScript      string
}

func addProjects(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	title := r.PostForm.Get("title")
	description := r.PostForm.Get("description")
	startDate := r.PostForm.Get("startdate")
	endDate := r.PostForm.Get("enddate")

	_, errQuery := connection.Conn.Exec(context.Background(), "INSERT INTO public.tb_projects(title, start_date, end_date, description) VALUES($1, $2, $3, $4)", title, startDate, endDate, description)

	if errQuery != nil {
		fmt.Println("Message : " + errQuery.Error())
		return
	}

	// var newProject = Project{
	// 	Title:       title,
	// 	Description: description,
	// 	StartDate:   startDate,
	// 	EndDate:     endDate,
	// }

	// Untuk Push ke Array projects
	// projects = append(projects, newProject)
	fmt.Println(startDate)
	fmt.Println(endDate)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/index.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	dataProject, errQuery := connection.Conn.Query(context.Background(), "SELECT Id, title, start_date, end_date, description FROM public.tb_projects")

	if errQuery != nil {
		fmt.Println("Message2 : " + errQuery.Error())
		return
	}

	var result []Project

	for dataProject.Next() {
		var each = Project{}

		err := dataProject.Scan(&each.Id, &each.Title, &each.StartDate, &each.EndDate, &each.Description)
		if err != nil {
			fmt.Println("Message dataProject : " + err.Error())
			return
		}

		each.FormatStartDate = each.StartDate.Format("2 January 2006")
		each.FormatEndDate = each.EndDate.Format("2 January 2006")

		result = append(result, each)
	}

	resData := map[string]interface{}{
		"Projects": result,
	}

	// dataProject := map[string]interface{}{
	// 	"Projects": projects,
	// }

	tmpt.Execute(w, resData)
}

func projectDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/projectDetail.html")

	if err != nil {
		w.Write([]byte("Message1 : " + err.Error()))
		return
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// w.Write([]byte("Message : " + err.Error()))

	var ProjectDetail = Project{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, title, start_date, end_date, description FROM tb_projects WHERE id = $1", id).Scan(&ProjectDetail.Id, &ProjectDetail.Title, &ProjectDetail.StartDate, &ProjectDetail.EndDate, &ProjectDetail.Description)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message2: " + err.Error()))
	}

	ProjectDetail.FormatStartDate = ProjectDetail.StartDate.Format("2 january 2006")
	ProjectDetail.FormatEndDate = ProjectDetail.EndDate.Format("2 january 2006")

	dataDetail := map[string]interface{}{
		"Project": ProjectDetail,
	}

	tmpt.Execute(w, dataDetail)
}

// ---------------------

func project(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/addProject.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}
	tmpt.Execute(w, nil)
}

func editProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/editProject.html")
	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	in, _ := strconv.Atoi(mux.Vars(r)["in"])

	var EditProject = Project{}

	errQuery := connection.Conn.QueryRow(context.Background(), "SELECT id, title, start_date, end_date, description FROM public.tb_projects WHERE id = $1", in).Scan(&EditProject.Id, &EditProject.Title, &EditProject.StartDate, &EditProject.EndDate, &EditProject.Description)

	if errQuery != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	}

	dataEdit := map[string]interface{}{
		"Project": EditProject,
	}

	EditProject.FormatStartDate = EditProject.StartDate.Format("2 January 2006")
	EditProject.FormatEndDate = EditProject.EndDate.Format("2 January 2006")

	tmpt.Execute(w, dataEdit)
}
func formEditProject(w http.ResponseWriter, r *http.Request) {
	in, _ := strconv.Atoi(mux.Vars(r)["in"])
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	title := r.PostForm.Get("title")
	description := r.PostForm.Get("description")
	startDate := r.PostForm.Get("startdate")
	endDate := r.PostForm.Get("enddate")

	_, errQuery := connection.Conn.Exec(context.Background(), "UPDATE public.tb_projects SET title=$1, start_date=$2, end_date=$3, description=$4 WHERE id=$5;", title, startDate, endDate, description, in)

	if errQuery != nil {
		fmt.Println("Message : " + errQuery.Error())
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func deleteProjects(w http.ResponseWriter, r *http.Request) {

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	_, errQuery := connection.Conn.Exec(context.Background(), "DELETE FROM tb_projects WHERE id = $1", index)

	if errQuery != nil {
		fmt.Println("Message : " + errQuery.Error())
		return
	}

	// projects = append(projects[:index], projects[index+1:]...)

	http.Redirect(w, r, "/", http.StatusFound)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/contact.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}
	tmpt.Execute(w, nil)
}
