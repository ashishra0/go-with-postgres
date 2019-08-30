package main

import (
	"database/sql"
)

type meal struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cuisine  string `json:"cuisine"`
	Category string `json:"category"`
}

func (m *meal) getMeal(db *sql.DB) error {
	return db.QueryRow("SELECT name, cuisine, category FROM dishes WHERE id=$1", m.ID).Scan(
		&m.Name, &m.Cuisine, &m.Category)
}

func (m *meal) updateMeal(db *sql.DB) error {
	_, err := db.Exec("UPDATE dishes SET name=$1, cuisine=$2, category=$3 WHERE id=$4", m.Name, m.Cuisine, m.Category)
	return err
}

func (m *meal) deleteMeal(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM dishes WHERE id=$1", m.ID)
	return err
}

func (m *meal) createMeal(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO dishes(name, cuisine, category) VALUES ($1, $2, $3) RETURNING id", m.Name, m.Cuisine, m.Category).Scan(&m.ID)
	if err != nil {
		return err
	}

	return nil
}

// This function fetches records from the products table.
func getMeals(db *sql.DB) ([]meal, error) {
	rows, err := db.Query(
		"SELECT id, name, cuisine, category FROM dishes")
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	meals := []meal{}
	for rows.Next() {
		var m meal
		if err := rows.Scan(&m.ID, &m.Name, &m.Cuisine, &m.Category); err != nil {
			return nil, err
		}
		meals = append(meals, m)
	}
	return meals, nil
}
