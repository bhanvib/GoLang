package main

import (
	"encoding/json"
	"fmt"
	"io"
	
	"net/http"

	"github.com/gorilla/mux"
	"project.com/DAL"
)

func GetProperty(w http.ResponseWriter, r *http.Request) {
	//io.WriteString(w,"get property")
	properties := DAL.GetProperty()
	json.NewEncoder(w).Encode(properties)
}
func GetPropertyByPropId(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	propid := params["propertyId"]
	prop := DAL.GetPropertyByPropId(propid)
	//convert golang object to json and sending to Response object
	json.NewEncoder(w).Encode(prop)
}
func PutProperty(w http.ResponseWriter, r *http.Request) {
	var prop DAL.Property
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&prop)
    fmt.Println(err)
    fmt.Println(prop)
    DAL.InsertProperty(prop)
}


func PostProperty(w http.ResponseWriter, r *http.Request) {
	// Accesing data send in Request body of Request object
	//r.Body  : JSOn
	var prop DAL.Property
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&prop)
	fmt.Println(err)
	fmt.Println(prop)
	DAL.InsertProperty(prop)
}
func UpdateProperty(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "update property")
}
func DeletePropertyByPropId(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	propid := params["propertyid"]
	DAL.GetPropertyByPropId(propid)
}
func FilterProperty(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	city := params["city"]
	bedroom := params["bedroom"]

	prop := DAL.GetPropertyByPropId(city)
	props := DAL.GetPropertyByPropId(bedroom)
	
	json.NewEncoder(w).Encode(prop)
	json.NewEncoder(w).Encode(props)
}
func main() {
	
	Router := mux.NewRouter()
	Router.HandleFunc("/property", GetProperty).Methods("GET")
	Router.HandleFunc("/property/{propid}", GetPropertyByPropId).Methods("GET")
	Router.HandleFunc("/property", PostProperty).Methods("POST")
	Router.HandleFunc("/property/{propid}", PutProperty).Methods("PUT")
	Router.HandleFunc("/property", UpdateProperty).Methods("UPDATE")
	Router.HandleFunc("/property/{propid}", DeletePropertyByPropId).Methods("DELETE")
	Router.HandleFunc("/property/{city}{bedroom}", FilterProperty).Methods("FILTER")
	http.ListenAndServe(":5050", Router)
}
