package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"

    pusher "github.com/pusher/pusher-http-go"
)

// credentials of my test pusher app, usually we should not put these directly in the code
var client = pusher.Client{
    AppId:   "661942",
    Key:     "780d41210d8a237fc9a8",
    Secret:  "df10df3c620d0746608d",
    Cluster: "ap1",
    Secure:  true,
}

type user struct {
    Name  string `json:"name" xml:"name" form:"name" query:"name"`
    Email string `json:"email" xml:"email" form:"email" query:"email"`
}

func registerUser(rw http.ResponseWriter, req *http.Request) {
    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
        panic(err)
    }

    var newUser user

    err = json.Unmarshal(body, &newUser)
    if err != nil {
        panic(err)
    }

    // trigger Pusher event
    client.Trigger("update", "new-user", newUser)

    json.NewEncoder(rw).Encode(newUser)
}

func pusherAuth(res http.ResponseWriter, req *http.Request) {
    params, _ := ioutil.ReadAll(req.Body)
    response, err := client.AuthenticatePrivateChannel(params)
    if err != nil {
        panic(err)
    }

    fmt.Fprintf(res, string(response))
}

func main() {
    // my chat app view
    http.Handle("/", http.FileServer(http.Dir("./client")))

    // register new user
    http.HandleFunc("/user/register", registerUser)

    // authorize client to use pusher channels
    http.HandleFunc("/pusher/auth", pusherAuth)

    fmt.Printf("Server is running...")

    log.Fatal(http.ListenAndServe(":8090", nil))
}