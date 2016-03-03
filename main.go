package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"lastName"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

var users []User

func loadUsers() ([]User, error) {
	data, err := ioutil.ReadFile("users.json")
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &users)
	return users, nil
}

func addDefaultHeaders(fn func(http.ResponseWriter, *http.Request, httprouter.Params)) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
		if origin := req.Header.Get("Origin"); origin != "" {
			rw.Header().Set("Access-Control-Allow-Origin", origin)
		}
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
		rw.Header().Set("Access-Control-Allow-Credentials", "true")
		fn(rw, req, params)
	}
}

func allHandler(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	jsonResponse, err := json.Marshal(users)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(jsonResponse)
}

func editUserHandler(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	id, err := strconv.Atoi(params.ByName("id"))
	currentUser := User{}
	var jsonResponse []byte
	for _, user := range users {
		if id == user.Id {
			currentUser = user
		}
	}
	jsonResponse, err = json.Marshal(currentUser)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(jsonResponse)
}

func dummyHandler(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
}

func updateUserHandler(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	newUser := loadFromJSON(rw, req)
	for i, user := range users {
		if newUser.Id == user.Id {
			users[i] = newUser
		}
	}
	saveToJSON()
}

func createUserHandler(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	newUser := loadFromJSON(rw, req)
	users = append(users, newUser)
	saveToJSON()
}

func loadFromJSON(rw http.ResponseWriter, req *http.Request) User {
	decoder := json.NewDecoder(req.Body)

	var newUser User
	err := decoder.Decode(&newUser)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	return newUser
}

func saveToJSON() {
	data, _ := json.Marshal(users)
	ioutil.WriteFile("users.json", data, 0644)
}

func main() {
	_, err := loadUsers()
	if err != nil {
		fmt.Println("Error loading the DB")
		os.Exit(1)
		return
	}
	router := httprouter.New()
	router.GET("/all/", addDefaultHeaders(allHandler))
	router.GET("/user/:id/edit/", addDefaultHeaders(editUserHandler))

	router.OPTIONS("/users/update/", addDefaultHeaders(dummyHandler))
	router.PUT("/users/update/", addDefaultHeaders(updateUserHandler))

	router.OPTIONS("/users/create/", addDefaultHeaders(dummyHandler))
	router.POST("/users/create/", addDefaultHeaders(createUserHandler))

	fmt.Println("listening in port 8080")
	http.ListenAndServe(":8080", router)
}
