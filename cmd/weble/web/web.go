package web

import (
	"fmt"
	"goweb/pkg/ble"
	"html/template"
	"net/http"
)

var tmpl *template.Template

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "index.html", nil)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "about.html", nil)
}

func scanHandler(w http.ResponseWriter, r *http.Request) {
	result, err := ble.Scan()
	if err != nil {
		http.Error(w, "Error scanning BLE devices: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf(result)
	tmpl.ExecuteTemplate(w, "scan.html", result)
}

func Run() error {
	tmpl = template.Must(template.ParseGlob("web/templates/*.html"))

	assetsHandler := http.StripPrefix("/assets", http.FileServer(http.Dir("web/assets")))
	http.Handle("/assets/", assetsHandler)

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/scan", scanHandler)
	return http.ListenAndServe(":8080", nil)
}
