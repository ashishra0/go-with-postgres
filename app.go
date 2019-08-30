package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// App struct exposes references to the router and the databases
// that the application uses.
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// Initialize method will take in the details required to connect to the database.
// It will create a database connection and wire up the routes
// to respond according to the requirements.
func (a *App) Initialize(user, password, dbname, sslmode string) {
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", user, password, dbname, sslmode)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

// Run will simply start the application
func (a *App) Run() {
	log.Fatal(http.ListenAndServe(":8000", a.Router))
}

// handler to fetch a single meal
func (a *App) getMeal(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	m := meal{ID: id}
	if err := m.getMeal(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Product not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, m)
}

// A handler to show all meals
func (a *App) getMeals(w http.ResponseWriter, r *http.Request) {
	meals, err := getMeals(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, meals)
}

// A handler to create a meal
func (a *App) createMeal(w http.ResponseWriter, r *http.Request) {
	var m meal
	decoder := json.NewDecoder(r.Body)
	// we need to decode because the body is in JSON
	if err := decoder.Decode(&m); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := m.createMeal(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, m)
}

// A handler to update a meal
func (a *App) updateMeal(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid meal ID")
		return
	}
	var m meal
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
	}
	m.ID = id
	if err := m.updateMeal(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, m)
}

// A handler to delete the meal
func (a *App) deleteMeal(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}

	p := meal{ID: id}
	if err := p.deleteMeal(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/meals", a.getMeals).Methods("GET")
	a.Router.HandleFunc("/meal", a.createMeal).Methods("POST")
	a.Router.HandleFunc("/meal/{id:[0-9]+}", a.getMeal).Methods("GET")
	a.Router.HandleFunc("/meal/{id:[0-9]+}", a.updateMeal).Methods("PUT")
	a.Router.HandleFunc("/meal/{id:[0-9]+}", a.deleteMeal).Methods("DELETE")
}
