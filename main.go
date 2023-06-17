package main

import(
    "net/http"
    "github.com/gorilla/mux"
    "encoding/json"
    "fmt"
    "log"
    "strconv"
)

type User struct{
    Id int
    Name string
    Email string
    Password string
  }

  var users = []User{ // the slice of users coming from the database
    {1, "John", "john@gmail.com", "123456"},
    {2, "Bob", "bob@gmail.com", "123456"},
    {3, "Sharon", "sharon@gmail.com", "123456"},
  }

  func index(w http.ResponseWriter, r *http.Request ){
    w.WriteHeader(http.StatusOK)
    fmt.Println(users)
    json.NewEncoder(w).Encode(users)
  }



func getUserByName(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)
  for _, user := range users {
    if user.Name == params["name"] {
      w.WriteHeader(http.StatusOK)
      json.NewEncoder(w).Encode(user)
      return
    }
  }
  w.WriteHeader(http.StatusNotFound)
  json.NewEncoder(w).Encode("User not found")
}


func getUserById(w http.ResponseWriter, r *http.Request ){
  params := mux.Vars(r)
  fmt.Println(params)
  userId, err := strconv.Atoi(params["id"])
  fmt.Println(userId)
  if err != nil {
    log.Fatal(err)
  }
  for _, user := range users {
    if user.Id == userId {
      w.WriteHeader(http.StatusOK)
      json.NewEncoder(w).Encode(user)
      return
    }
  }
  w.WriteHeader(http.StatusNotFound)
  json.NewEncoder(w).Encode("User not found")

}


func createUser(w http.ResponseWriter, r *http.Request ){
  var user User // the User struct coming from the request
  _ = json.NewDecoder(r.Body).Decode(&user)
  users = append(users, user)
  w.WriteHeader(http.StatusCreated)
  json.NewEncoder(w).Encode(user)
}

func updateUserById(w http.ResponseWriter, r *http.Request){
  // ensure you convert the id to an int
  params := mux.Vars(r)
  id, _ := strconv.Atoi(params["id"])
  for index, user := range users {
    if user.Id == id {
      users = append(users[:index], users[index+1:]...)
      var user User
      _ = json.NewDecoder(r.Body).Decode(&user)
      user.Id = id
      users = append(users, user)
      w.WriteHeader(http.StatusCreated)
      json.NewEncoder(w).Encode(user)
      return
    }
  }
  w.WriteHeader(http.StatusNotFound)
  json.NewEncoder(w).Encode("User not found")
}

func deleteUserById(w http.ResponseWriter, r *http.Request){
  params := mux.Vars(r)
  id, _ := strconv.Atoi(params["id"])
  for index, user := range users {
    if user.Id == id {
      users = append(users[:index], users[index+1:]...)
      w.WriteHeader(http.StatusNoContent)
      return
    }
  }
  w.WriteHeader(http.StatusNotFound)
  json.NewEncoder(w).Encode("User not found")
}

func main(){
  router := mux.NewRouter().StrictSlash(true)
  router.HandleFunc("/", index).Methods("GET")
  router.HandleFunc("/users/{id}", getUserById).Methods("GET")
  router.HandleFunc("/users/{name}", getUserByName).Methods("GET")
  router.HandleFunc("/users", createUser).Methods("POST")
  router.HandleFunc("/users/{id}", updateUserById).Methods("PUT")
  router.HandleFunc("/users/{id}", deleteUserById).Methods("DELETE")
  log.Fatal(http.ListenAndServe(":8081", router))
  }
