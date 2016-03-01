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

func addDefaultHeaders(fn http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		if origin := req.Header.Get("Origin"); origin != "" {
			rw.Header().Set("Access-Control-Allow-Origin", origin)
		}
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
		rw.Header().Set("Access-Control-Allow-Credentials", "true")
		fn(rw, req)
	}
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
	http.HandleFunc("/all/", addDefaultHeaders(allHandler))
	fmt.Println("listening in port 8080")
	http.ListenAndServe(":8080", nil)

}
