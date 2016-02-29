package main

import (
	"encoding/json"
	"io/ioutil"
    "os"
)


var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

type User struct{
    id int 'json:"id"'
    name string 'json:"name"'
    last_name string 'json:"last_name"'
    email string 'json:"email"'
    username string 'json:"username"'
}

func loadUsers() ([]User, error) {
    data, err := ioutil.ReadFile("users.json")
    if err != nil{
        return nil, err
    }
    var users []User
    json.Unmarshal(data, &users)
    return users, nil    
}

func saveHandler(rw http.ResponseWriter, req *http.Request) {

}

func editHandler(rw http.ResponseWriter, req *http.Request) {

}
func main() {
	/*http.HandleFunc("/view/", makeHandler(viewHandler))
	  http.HandleFunc("/edit/", makeHandler(editHandler))*/
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(":8080", nil)

}
