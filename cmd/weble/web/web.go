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

type advert struct {
	Addr        string
	Connectable bool
	RSSI        int
	LocalName   string
	Company     string
}

func scanHandler(w http.ResponseWriter, r *http.Request) {
	dm, err := ble.Scan()
	if err != nil {
		http.Error(w, "Error scanning BLE devices: "+err.Error(), http.StatusInternalServerError)
		return
	}
	s := dm.String()
	fmt.Print(s)

	// must convert the goble.Advertisement structs which has accessor functions for each property into
	// a format with values that can be displayed by the template
	rows := make([]advert, 0)
	for _, a := range dm {
		row := advert{}
		row.Addr = a.Addr().String()
		row.Connectable = a.Connectable()
		row.RSSI = a.RSSI()
		row.LocalName = a.LocalName()
		row.Company = ble.GetCompanyName(a)
		rows = append(rows, row)
	}
	data := &map[string]any{
		"Rows": rows,
	}
	tmpl.ExecuteTemplate(w, "scan.html", data)
}

func Run() error {
	tmpl = template.Must(template.ParseGlob("web/templates/*.html"))

	assetsHandler := http.StripPrefix("/assets", http.FileServer(http.Dir("web/assets")))
	http.Handle("/assets/", assetsHandler)

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/scan", scanHandler)
	fmt.Printf("Starting server on :8080\n")
	return http.ListenAndServe(":8080", nil)
}
