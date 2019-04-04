package main

import (
    "encoding/json"
    "fmt"
    "html/template"
    "log"
    "net/http"
    "strconv"

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
    // Bei Fehlern schreiend im Kreis rennen:
    if err != nil {
        panic(err)
    }
}

func FetchPlakate() []Plakat {
    // Liste von Plakaten aus der DB in ein Plakat-Array schieben:
    db, err := sqlx.Open("sqlite3", "./plakate.db")
    CheckError(err)
    defer db.Close()
    
    // Liste erzeugen:
    plakate := []Plakat{}
    db.Select(&plakate, "SELECT * FROM plakate")
    
    return plakate
}

func DeletePlakat(id string) {
    // Datenbank aufrufen:
    db, err := sqlx.Open("sqlite3", "./plakate.db")
    CheckError(err)
    defer db.Close()
    
    // Löschen, falls möglich:
    stmt, err := db.Prepare("delete from plakate where id = ?")
    _, err = stmt.Exec(strconv.Atoi(id))
    CheckError(err)
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
    _, err = stmt.Exec(lat, lon)
    CheckError(err)
    
    fmt.Fprintf(w, "Plakat erfolgreich eingetragen!")
}

func DelHandler(w http.ResponseWriter, r *http.Request) {
    // Plakat mit URL-Parameter "id" löschen:
    vars := mux.Vars(r)
    DeletePlakat(vars["id"])
    
    // Falls kein Fehler aufgetreten ist, umleiten auf /manageplakate:
    http.Redirect(w, r, "/manageplakate", http.StatusMovedPermanently)
}

func DelPostHandler(w http.ResponseWriter, r *http.Request) {
    // Plakat mit POST-Parameter "id" löschen:
    data := r.FormValue("id")
    DeletePlakat(data)
    
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