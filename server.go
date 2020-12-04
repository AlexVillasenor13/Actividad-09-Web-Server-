package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Grade struct {
	Student string
	Subject string
	Grade   float64
}

type Server struct {
	Subjects map[string]map[string]float64
	Students map[string]map[string]float64
}

var myServer Server

func (this *Server) AddGrade(grade *Grade, reply *string) error {

	for student, _ := range this.Subjects[grade.Subject] {
		if student == grade.Student {
			return errors.New("Existing grade")
		}
	}

	if this.findSubject(grade.Subject) {
		this.Subjects[grade.Subject][grade.Student] = grade.Grade
	} else {
		student := make(map[string]float64)
		student[grade.Student] = grade.Grade
		this.Subjects[grade.Subject] = student
	}
	if this.findStudent(grade.Student) {
		this.Students[grade.Student][grade.Subject] = grade.Grade
	} else {
		subject := make(map[string]float64)
		subject[grade.Subject] = grade.Grade
		this.Students[grade.Student] = subject
	}
	s := fmt.Sprintf("%.2f", grade.Grade)
	*reply = "Added Grade: " + grade.Subject + ": " + grade.Student + ": " + s
	return nil
}

func (this *Server) findStudent(name string) bool {
	for student, _ := range this.Students {
		if student == name {
			return true
		}
	}
	return false
}

func (this *Server) findSubject(name string) bool {
	for subject, _ := range this.Subjects {
		if subject == name {
			return true
		}
	}
	return false
}

func (this *Server) StudentProm(name string, reply *string) error {
	total_grades := 0.0
	total_subjects := 0.0
	for _, grade := range this.Students[name] {
		total_grades += grade
		total_subjects += 1
	}
	if total_subjects > 0 {
		aux := fmt.Sprintf("%.2f", total_grades/total_subjects)
		*reply = "<tr>" +
			"<td>" + name + "</td>" +
			"<td>" + aux + "</td>" +
			"</tr>"
		return nil
	} else {
		return errors.New("Wrong student")
	}
}

func (this *Server) SubjectProm(name string, reply *string) error {
	total_grades := 0.0
	total_students := 0.0
	for _, grade := range this.Subjects[name] {
		total_grades += grade
		total_students += 1
	}
	if total_students > 0 {
		aux := fmt.Sprintf("%.2f", total_grades/total_students)
		*reply = "<tr>" +
			"<td>" + name + "</td>" +
			"<td>" + aux + "</td>" +
			"</tr>"
		return nil
	} else {
		return errors.New("Wrong subject")
	}
}

func (this *Server) StudentPromFloat(name string, reply *float64) error {
	total_grades := 0.0
	total_subjects := 0.0
	for _, grade := range this.Students[name] {
		total_grades += grade
		total_subjects += 1
	}
	if total_subjects > 0 {
		*reply = total_grades / total_subjects
		return nil
	} else {
		return errors.New("Wrong student")
	}
}

func (this *Server) GralProm(reply_value string, reply *string) error {
	total_grades_students := 0.0
	total_students := 0.0
	for student, _ := range this.Students {
		prom := 0.0
		err := this.StudentPromFloat(student, &prom)
		if err != nil {
			return err
		}
		total_grades_students += prom
		total_students += 1
	}
	if total_students > 0 {
		aux := fmt.Sprintf("%.2f", total_grades_students/total_students)
		*reply = "<tr>" +
			"<td>" + aux + "</td>" +
			"</tr>"
		return nil
	} else {
		return errors.New("No hay estudiantes")
	}
}

func cargarHtml(a string) string {
	html, _ := ioutil.ReadFile(a)

	return string(html)
}

func form_student(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		cargarHtml("form-prom-student.html"),
	)
}

func form_subject(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		cargarHtml("form-prom-subject.html"),
	)
}

func form_new(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		cargarHtml("add-new-form.html"),
	)
}

func grl_prom(res http.ResponseWriter, req *http.Request) {
	var result string
	err := myServer.GralProm(result, &result)
	if err != nil {
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			cargarHtml("error.html"),
			err,
		)
	} else {
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			cargarHtml("grl-prom.html"),
			result,
		)
	}
}

func add_new(res http.ResponseWriter, req *http.Request) {
	var result string
	fmt.Println(req.Method)
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		fmt.Println(req.PostForm)
		grade_value, _ := strconv.ParseFloat(req.FormValue("grade"), 8)
		grade := Grade{Student: req.FormValue("student"),
			Subject: req.FormValue("subject"),
			Grade:   grade_value}
		err := myServer.AddGrade(&grade, &result)
		if err != nil {
			res.Header().Set(
				"Content-Type",
				"text/html",
			)
			fmt.Fprintf(
				res,
				cargarHtml("error.html"),
				err,
			)
		} else {
			fmt.Println(myServer.Students)
			res.Header().Set(
				"Content-Type",
				"text/html",
			)
			fmt.Fprintf(
				res,
				cargarHtml("answer.html"),
				result,
			)
		}

	}
}

func prom_student(res http.ResponseWriter, req *http.Request) {
	var result string
	fmt.Println(req.Method)
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		fmt.Println(req.PostForm)
		student := req.FormValue("student")
		err := myServer.StudentProm(student, &result)
		if err != nil {
			res.Header().Set(
				"Content-Type",
				"text/html",
			)
			fmt.Fprintf(
				res,
				cargarHtml("error.html"),
				err,
			)
		} else {
			fmt.Println(myServer.Students)
			res.Header().Set(
				"Content-Type",
				"text/html",
			)
			fmt.Fprintf(
				res,
				cargarHtml("prom-student.html"),
				result,
			)
		}
	}
}

func prom_subject(res http.ResponseWriter, req *http.Request) {
	var result string
	fmt.Println(req.Method)
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		fmt.Println(req.PostForm)
		subject := req.FormValue("subject")
		err := myServer.SubjectProm(subject, &result)
		if err != nil {
			res.Header().Set(
				"Content-Type",
				"text/html",
			)
			fmt.Fprintf(
				res,
				cargarHtml("error.html"),
				err,
			)
		} else {
			fmt.Println(myServer.Students)
			res.Header().Set(
				"Content-Type",
				"text/html",
			)
			fmt.Fprintf(
				res,
				cargarHtml("prom-subject.html"),
				result,
			)
		}
	}
}

func main() {

	myServer.Students = make(map[string]map[string]float64)
	myServer.Subjects = make(map[string]map[string]float64)

	http.HandleFunc("/new", form_new)
	http.HandleFunc("/new-result", add_new)

	http.HandleFunc("/student", form_student)
	http.HandleFunc("/subject", form_subject)

	http.HandleFunc("/prom-student", prom_student)
	http.HandleFunc("/prom-subject", prom_subject)
	http.HandleFunc("/gral-prom", grl_prom)
	fmt.Println("Corriendo servidor de Calificaciones...")
	http.ListenAndServe(":1306", nil)

}
