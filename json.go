package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"lastName"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func loadUsers() ([]User, error) {
	data, err := ioutil.ReadFile("users.json")
	if err != nil {
		return nil, err
	}
	var users []User
	json.Unmarshal(data, &users)
	return users, nil
}

func allHandler(rw http.ResponseWriter, req *http.Request) {
	users, err := loadUsers()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	var jsonResponse []byte
	jsonResponse, err = json.Marshal(users)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(jsonResponse)
}

func main() {
	http.HandleFunc("/all/", allHandler)
	fmt.Println("listening in port 8000")
	http.ListenAndServe(":8000", nil)

}
