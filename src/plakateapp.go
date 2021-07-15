/*
 * The contents of this file are subject to the terms of the
 * Common Development and Distribution License, Version 1.0 only
 * (the "License").  You may not use this file except in compliance
 * with the License.
 *
 * See the file LICENSE in this distribution for details.
 * A copy of the CDDL is also available via the Internet at
 * http://www.opensource.org/licenses/cddl1.txt
 *
 * When distributing Covered Code, include this CDDL HEADER in each
 * file and include the contents of the LICENSE file from this
 * distribution.
 */

package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	// Bibliotheken aus GitHub:
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Plakat struct {
	// Mapping DB <-> Go-Datentypen:
	ID        int     `db:"id"`
	Latitude  float32 `db:"lat"`
	Longitude float32 `db:"lon"`
	Location  string  `db:"location"`
}

// Reverse-Geocoding via Nominatim:

type ReverseGeoCode struct {
	// Bildet <reversegeocode> ab.
	XMLName     xml.Name    `xml:"reversegeocode"`
	AdressParts AdressParts `xml:"addressparts"`
}

type AdressParts struct {
	// Bildet <reversegeocode><adressparts> ab.
	XMLName     xml.Name `xml:"addressparts"`
	HouseNumber string   `xml:"house_number"`
	Road        string   `xml:"road"`
	Suburb      string   `xml:"suburb"`
	District    string   `xml:"city_district"`
	City        string   `xml:"city"`
	State       string   `xml:"state"`
	Postcode    string   `xml:"postcode"`
	Country     string   `xml:"country"`
	CountryCode string   `xml:"country_code"`
}

// ----------------------------------------

func CheckError(err error) {
	// Bei Fehlern schreiend im Kreis rennen:
	if err != nil {
		panic(err)
	}
}

func GetXML(url string) ([]byte, error) {
	// XML herunterladen:
	resp, err := http.Get(url)
	CheckError(err)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("Statusfehler: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	CheckError(err)

	return data, nil
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
	idint, err := strconv.Atoi(id)
	CheckError(err)
	_, err = stmt.Exec(idint)
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

	// Geocoding:
	url := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?lat=%s&lon=%s", lat, lon)
	xmlBytes, err := GetXML(url)
	CheckError(err)

	var xmlFile ReverseGeoCode
	xml.Unmarshal(xmlBytes, &xmlFile)

	adressData := xmlFile.AdressParts

	p1 := ""
	p2 := ""
	p3 := ""

	if adressData.Road != "" {
		if adressData.HouseNumber != "" {
			p1 = fmt.Sprintf("%s %s, ", adressData.Road, adressData.HouseNumber)
		} else {
			// Manche Straßen werden ohne konkrete Hausnummer gefunden.
			p1 = fmt.Sprintf("%s, ", adressData.Road)
		}
	}

	if adressData.District != "" {
		p2 = fmt.Sprintf("%s, ", adressData.District)
	}

	if adressData.Postcode != "" {
		p3 = fmt.Sprintf("%s %s", adressData.Postcode, adressData.City)
	}

	location := fmt.Sprintf("%s%s%s", p1, p2, p3)

	// Datenbank aufrufen:
	db, err := sqlx.Open("sqlite3", "./plakate.db")
	CheckError(err)
	defer db.Close()

	stmt, err := db.Prepare("insert into plakate (lat, lon, location) values (?, ?, ?)")
	_, err = stmt.Exec(lat, lon, location)
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

func serveHTTP() {
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

func main() {
	startServer()
}
