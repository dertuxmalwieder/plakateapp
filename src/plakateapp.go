package main

import (
    "encoding/json"
    "fmt"
    "html/template"
    "log"
    "net/http"

    // Bibliotheken aus GitHub:
    _ "github.com/mattn/go-sqlite3"
    "github.com/jmoiron/sqlx"
    "github.com/gorilla/mux"
)

type Plakat struct {
    // Mapping DB <-> Go-Datentypen:
    ID          int     `db:"id"`
    Latitude    float32 `db:"lat"`
    Longitude   float32 `db:"lon"`
}

func CheckError(err error) {
    if err != nil {
        panic(err)
    }
}

func FetchPlakate() []Plakat {
    db, err := sqlx.Open("sqlite3", "./plakate.db")
    CheckError(err)
    defer db.Close()
    
    // Liste erzeugen:
    plakate := []Plakat{}
    db.Select(&plakate, "SELECT * FROM plakate")
    
    return plakate
}
    
// ----------------------------------------
   
func HomeHandler(w http.ResponseWriter, r *http.Request) {
    // Startseite aufrufen:
    tmpl := template.Must(template.ParseFiles("templates/index.htm"))
    tmpl.Execute(w, "")
}

func ManagePlakateHandler(w http.ResponseWriter, r *http.Request) {
    // Datenbank aufrufen:
    plakate := FetchPlakate()
    
    // delete.htm mit Plakate-struct aufrufen:
    tmpl := template.Must(template.ParseFiles("templates/delete.htm"))
    tmpl.Execute(w, plakate)
}

func ListPlakateHandler(w http.ResponseWriter, r *http.Request) {
    // Datenbank aufrufen:
    plakate := FetchPlakate()
    
    // JSON-Objekt ausgeben:
    jsonobj, err := json.Marshal(plakate)
    CheckError(err)
    fmt.Fprintf(w, "%s", string(jsonobj))
}

func NeuesPlakatHandler(w http.ResponseWriter, r *http.Request) {
    // Plakat mit "lat" und "lon" erzeugen:
    lat := r.FormValue("lat")
    lon := r.FormValue("lon")
    
    // Datenbank aufrufen:
    db, err := sqlx.Open("sqlite3", "./plakate.db")
    CheckError(err)
    defer db.Close()
    
    stmt, err := db.Prepare("insert into plakate (lat, lon) values (?, ?)")
    CheckError(err)
    _, err = stmt.Exec(lat, lon)
    CheckError(err)
    
    fmt.Fprintf(w, "Plakat erfolgreich eingetragen!")
}

func DelHandler(w http.ResponseWriter, r *http.Request) {
    // Plakat mit vars["id"] löschen:
    vars := mux.Vars(r)

    // Datenbank aufrufen:
    db, err := sqlx.Open("sqlite3", "./plakate.db")
    CheckError(err)
    defer db.Close()
    
    // Löschen, falls möglich:
    stmt, err := db.Prepare("delete from plakate where id = ?")
    CheckError(err)
    _, err = stmt.Exec(vars["id"])
    CheckError(err)
    
    // Falls kein Fehler aufgetreten ist, umleiten auf /manageplakate:
    http.Redirect(w, r, "/manageplakate", http.StatusMovedPermanently)
}

func DelPostHandler(w http.ResponseWriter, r *http.Request) {
    data := r.FormValue("id")

    // Datenbank aufrufen:
    db, err := sqlx.Open("sqlite3", "./plakate.db")
    CheckError(err)
    defer db.Close()
    
    // Löschen, falls möglich:
    stmt, err := db.Prepare("delete from plakate where id = ?")
    CheckError(err)
    _, err = stmt.Exec(data)
    CheckError(err)
    
    // Falls kein Fehler aufgetreten ist, umleiten auf /manageplakate:
    http.Redirect(w, r, "/manageplakate", http.StatusMovedPermanently)
}

func main() {
    // Routing:
    r := mux.NewRouter()
    r.HandleFunc("/", HomeHandler)
    r.HandleFunc("/manageplakate", ManagePlakateHandler)
    r.HandleFunc("/listplakate", ListPlakateHandler).Methods("POST")
    r.HandleFunc("/neuesplakat", NeuesPlakatHandler).Methods("POST")
    r.HandleFunc("/del/{id:[0-9]+}", DelHandler)
    r.HandleFunc("/delpost", DelPostHandler).Methods("POST")
    
    // script.js auch aus /templates ausliefern:
    r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./templates"))))
    http.Handle("/", r)

    // Server starten:
    log.Fatal(http.ListenAndServe(":6090", nil))
}