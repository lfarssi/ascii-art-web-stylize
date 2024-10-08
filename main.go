package main

import (
	"ascii-art/functions"
	"fmt"
	"html/template"
	"net/http"

	"strings"
)


type Data struct {
	Str string
	Banner string
	Res string
	A	template.HTML
}

func processHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
        http.ServeFile(w, r, "400.html")
        return
    }
	temp, err := template.ParseFiles("home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	

	var data Data
	
	data.Str = r.FormValue("data")
	if data.Str == "" {
        http.ServeFile(w, r, "400.html")
		return
	}
	if len(data.Str) > 200 {
        http.ServeFile(w, r, "400.html")
		return
	}

	data.Banner = r.FormValue("banner")
	if !function.BannerExists(data.Banner) {
		http.Error(w, "Banner not found", http.StatusNotFound)
		return
	}

	data.Str = strings.ReplaceAll(data.Str, "\r\n", "\n")
	
	data.Res = function.TraitmentData(w, data.Banner, data.Str)
	if data.Res == "" { // If TraitmentData failed to generate the result
        http.ServeFile(w, r, "500.html")
		return
	}
	data.A = template.HTML(strings.ReplaceAll(data.Res, "\n", "<br>"))

	if err := temp.Execute(w, data); err != nil {
        http.ServeFile(w, r, "500.html")
		return
	}
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
        http.ServeFile(w, r, "404.html")
		return
	}
	if r.Method != "GET" {
        http.ServeFile(w, r, "400.html")
        return
    }
	t, err := template.ParseFiles("home.html")
	if err != nil {
        http.ServeFile(w, r, "404.html")
		return
	}

	if err := t.Execute(w, nil); err != nil {
        http.ServeFile(w, r, "500.html")
		return
	}
}
func cssHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "style.css")
}

func main() {

	http.HandleFunc("/", pageHandler)
	http.HandleFunc("/ascii-art", processHandler)
	http.HandleFunc("/style.css", cssHandler)
	fmt.Println("Server is running at http://localhost:8088")
	err := http.ListenAndServe(":8088", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}