package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Env environment
type Env struct {
	db *sql.DB
}

// User data struct
type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
	BirthDate int    `json:"birth_date"`
}

// Location data struct
type Location struct {
	ID       int    `json:"id"`
	Place    string `json:"place"`
	Country  string `json:"country"`
	City     string `json:"city"`
	Distance int    `json:"distance"`
}

// Visit data struct
type Visit struct {
	ID        int `json:"id"`
	Location  int `json:"location"`
	User      int `json:"user"`
	VisitedAt int `json:"visited_at"`
	Mark      int `json:"mark"`
}

func (u *User) insert(env *Env) error {
	_, err := env.db.Exec("insert into users values (?,?,?,?,?,?)", u.ID, u.Email, u.FirstName, u.LastName, u.Gender, u.BirthDate)
	return err
}

func (v *Visit) insert(env *Env) error {
	_, err := env.db.Exec("insert into visits values (?,?,?,?,?)", v.ID, v.User, v.Location, v.VisitedAt, v.Mark)
	return err
}

func (l *Location) insert(env *Env) error {
	_, err := env.db.Exec("insert into locations values (?,?,?,?,?)", l.ID, l.Place, l.Country, l.City, l.Distance)
	return err
}

func (u *User) get(id int, env *Env) error {
	err := env.db.QueryRow("select * from users where user_id = ?", id).Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Gender, &u.BirthDate)
	return err
}

func (v *Visit) get(id int, env *Env) error {
	err := env.db.QueryRow("select * from visits where visit_id = ?", id).Scan(&v.ID, &v.User, &v.Location, &v.VisitedAt, &v.Mark)
	return err
}

func (l *Location) get(id int, env *Env) error {
	err := env.db.QueryRow("select * from locations where loc_id = ?", id).Scan(&l.ID, &l.Place, &l.Country, &l.City, &l.Distance)
	return err
}

// Entity just test
type Entity interface {
	insert(*Env) error
	//	update(Env) error
	get(int, *Env) error
}

// MakeEntity create entity from parametr
func MakeEntity(name string) (Entity, error) {
	var err error
	if name == "users" {
		return &User{}, err
	} else if name == "locations" {
		return &Location{}, err
	} else if name == "visits" {
		return &Visit{}, err
	}
	return nil, err
}

func (env *Env) createEntity(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	entity, err := MakeEntity(params["entity"])
	if err != nil {
		fmt.Println(err)
	}

	_ = json.NewDecoder(r.Body).Decode(&entity)
	err = entity.insert(env)
	if err != nil {
		fmt.Println(err)
	}

	json.NewEncoder(w).Encode(entity)

	/*
			check entity type
			create entity
			call methods

		if params["entity"] == "users" {
			var user User
			_ = json.NewDecoder(r.Body).Decode(&user)
			_, err := env.db.Exec("insert into users values (?,?,?,?,?,?)", user.ID, user.Email, user.FirstName, user.LastName, user.Gender, user.BirthDate)
			if err != nil {
				fmt.Println(err)
			}
		} else if params["entity"] == "locations" {
			var location Location
			_ = json.NewDecoder(r.Body).Decode(&location)
			_, err := env.db.Exec("insert into locations values (?,?,?,?,?)", location.ID, location.Place, location.Country, location.City, location.Distance)
			if err != nil {
				fmt.Println(err)
			}
			json.NewEncoder(w).Encode(location)
		} else if params["entity"] == "visits" {
			var visit Visit
			_ = json.NewDecoder(r.Body).Decode(&visit)
			_, err := env.db.Exec("insert into visits values (?,?,?,?,?)", visit.ID, visit.User, visit.Location, visit.VisitedAt, visit.Mark)
			if err != nil {
				fmt.Println(err)
			}
			json.NewEncoder(w).Encode(visit)
		}
	*/
}

// GetEntity just test
func (env *Env) GetEntity(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	entity, err := MakeEntity(params["entity"])
	if err != nil {
		fmt.Println(err)
	}

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Println(err)
	}

	err = entity.get(id, env)
	if err != nil {
		fmt.Println(err)
	}

	json.NewEncoder(w).Encode(entity)

	/*
		if params["entity"] == "users" {
			var user User
			err := env.db.QueryRow("select * from users where user_id = ?", params["id"]).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Gender, &user.BirthDate)
			if err != nil {
				fmt.Println(err)
			}
			json.NewEncoder(w).Encode(user)
		} else if params["entity"] == "locations" {
			var location Location
			err := env.db.QueryRow("select * from locations where loc_id = ?", params["id"]).Scan(&location.ID, &location.Place, &location.Country, &location.City, &location.Distance)
			if err != nil {
				fmt.Println(err)
			}
			json.NewEncoder(w).Encode(location)
		} else if params["entity"] == "visits" {
			var visit Visit
			err := env.db.QueryRow("select * from visits where visit_id = ?", params["id"]).Scan(&visit.ID, &visit.User, &visit.Location, &visit.VisitedAt, &visit.Mark)
			if err != nil {
				fmt.Println(err)
			}
			json.NewEncoder(w).Encode(visit)
		}
	*/
	return
}

func main() {
	db, err := sql.Open("mysql", "kot:pass@/travel_service")
	if err != nil {
		fmt.Println(err)
	}
	env := &Env{db: db}
	r := mux.NewRouter()
	r.HandleFunc("/{entity}/{id}", env.GetEntity).Methods("GET")
	r.HandleFunc("/{entity}/new", env.createEntity).Methods("POST")
	http.ListenAndServe(":8080", r)
}
