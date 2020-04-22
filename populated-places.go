package dataapi

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// This file creates a series of endpoints to return all possible names for
// populated places, associated with their county IDs from AHCB as of 1926-1928.

// PopPlacesCounty represents a county with ID and name
type PopPlacesCounty struct {
	CountyAHCB string `json:"county_ahcb"`
	County     string `json:"name"`
}

// PopPlacesPlace represents a place with ID and name
type PopPlacesPlace struct {
	PlaceID int    `json:"place_id"`
	Place   string `json:"place"`
}

// PopPlacesPlaceDetail represents details about a place
type PopPlacesPlaceDetail struct {
	PlaceID    int    `json:"place_id"`
	Place      string `json:"place"`
	County     string `json:"county"`
	CountyAHCB string `json:"county_ahcb"`
	State      string `json:"state"`
}

// PopPlacesCountiesInState returns a list of all the counties in a state, with
// IDs from AHCB.
func (s *Server) PopPlacesCountiesInState() http.HandlerFunc {

	query := `
		SELECT DISTINCT county_ahcb, county
		FROM popplaces_1926
		WHERE state = $1
		ORDER BY county;
		`

	stmt, err := s.Database.Prepare(query)
	if err != nil {
		log.Fatalln(err)
	}
	s.Statements["pop-places-counties-in-state"] = stmt // Will be closed at shutdown

	return func(w http.ResponseWriter, r *http.Request) {

		state := mux.Vars(r)["state"]
		state = strings.ToUpper(state)

		results := make([]PopPlacesCounty, 0)
		var row PopPlacesCounty

		rows, err := stmt.Query(state)
		if err != nil {
			log.Println(err)
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&row.CountyAHCB, &row.County)
			if err != nil {
				log.Println(err)
			}
			results = append(results, row)
		}
		err = rows.Err()
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		response, _ := json.Marshal(results)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(response))

	}
}

// PopPlacesPlacesInCounty returns a list of all the populated places in a county.
func (s *Server) PopPlacesPlacesInCounty() http.HandlerFunc {

	query := `
		SELECT place_id, place
		FROM popplaces_1926
		WHERE county_ahcb = $1
		ORDER BY place;
		`

	stmt, err := s.Database.Prepare(query)
	if err != nil {
		log.Fatalln(err)
	}
	s.Statements["pop-places-places-in-county"] = stmt // Will be closed at shutdown

	return func(w http.ResponseWriter, r *http.Request) {

		county := mux.Vars(r)["county"]
		county = strings.ToLower(county)

		results := make([]PopPlacesPlace, 0)
		var row PopPlacesPlace

		rows, err := stmt.Query(county)
		if err != nil {
			log.Println(err)
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&row.PlaceID, &row.Place)
			if err != nil {
				log.Println(err)
			}
			results = append(results, row)
		}
		err = rows.Err()
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		response, _ := json.Marshal(results)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(response))

	}
}

// PopPlacesPlacesInState returns a list of all the populated places in a state.
func (s *Server) PopPlacesPlacesInState() http.HandlerFunc {

	query := `
		SELECT place_id, place
		FROM popplaces_1926
		WHERE state = $1
		ORDER BY place;
		`

	stmt, err := s.Database.Prepare(query)
	if err != nil {
		log.Fatalln(err)
	}
	s.Statements["pop-places-places-in-state"] = stmt // Will be closed at shutdown

	return func(w http.ResponseWriter, r *http.Request) {

		state := mux.Vars(r)["state"]
		state = strings.ToUpper(state)

		results := make([]PopPlacesPlace, 0)
		var row PopPlacesPlace

		rows, err := stmt.Query(state)
		if err != nil {
			log.Println(err)
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&row.PlaceID, &row.Place)
			if err != nil {
				log.Println(err)
			}
			results = append(results, row)
		}
		err = rows.Err()
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		response, _ := json.Marshal(results)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(response))

	}
}

// PopPlacesPlace returns a list of all the populated places in a state.
func (s *Server) PopPlacesPlace() http.HandlerFunc {

	query := `
		SELECT place_id, place, county, county_ahcb, state
		FROM popplaces_1926
		WHERE place_id = $1
		`

	stmt, err := s.Database.Prepare(query)
	if err != nil {
		log.Fatalln(err)
	}
	s.Statements["pop-places-details"] = stmt // Will be closed at shutdown

	return func(w http.ResponseWriter, r *http.Request) {

		placeID, err := strconv.Atoi(mux.Vars(r)["place"])
		if err != nil {
			http.Error(w, "Bad request: place ID must be an integer", http.StatusBadRequest)
			return
		}

		var result PopPlacesPlaceDetail

		err = stmt.QueryRow(placeID).Scan(&result.PlaceID, &result.Place,
			&result.County, &result.CountyAHCB, &result.State)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, fmt.Sprintf("Not found: No place with id %v.", placeID), http.StatusNotFound)
				return
			}
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		response, _ := json.Marshal(result)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(response))

	}
}
