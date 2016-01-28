package main

import (
    "encoding/json"
    "fmt"
    "github.com/gorilla/mux"
    // "github.com/gorilla/sessions"
    "github.com/jinzhu/gorm"
    _ "github.com/lib/pq"
    "html/template"
    "io/ioutil"
    "log"
    "net/http"
    "os"
)

func main() {
    port := os.Getenv("PORT")

    if port == "" {
        log.Fatal("$PORT must be set")
    }

    db := initDB()
    me := loadMe()
    db.Create(&me)

    static := http.FileServer(http.Dir("static"))
    http.Handle("/", static)

    root := mux.NewRouter()

    root.HandleFunc("/api/profile/{name}", func(res http.ResponseWriter, req *http.Request) {
        username := mux.Vars(req)["name"]
        fmt.Println("GET /api/profile/:name")
        user := User{}
        db.Where("username = ?", username).First(&user)
        userJson, err := json.Marshal(user)
        if err != nil {
            http.Error(res, err.Error(), http.StatusInternalServerError)
            return
        }
        res.Write(userJson)
    }).Methods("GET")

    root.HandleFunc("/api/profile", func(res http.ResponseWriter, req *http.Request) {
        fmt.Println("POST /api/profile")
        decoder := json.NewDecoder(req.Body)
        var user User
        err := decoder.Decode(&user)
        if err != nil {
            http.Error(res, err.Error(), http.StatusInternalServerError)
        }
        fmt.Println(user)
        res.WriteHeader(200)
    }).Methods("POST")

    root.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
        fmt.Println("GET /")
        tmpl, _ := template.ParseGlob("pages/*.tmpl.html")
        tmpl.ExecuteTemplate(res, "index", nil)
    }).Methods("GET")

    http.ListenAndServe(":"+port, root)
}

func initDB() (db gorm.DB) {
    //export DATABASE_URL="dbname=learnaboutme sslmode=disable"
    dbUrl := os.Getenv("DATABASE_URL")
    if dbUrl == "" {
        log.Fatal("$DATABASE_URL must be set")
    }
    fmt.Println("Initializing DB...")
    db, err := gorm.Open("postgres", dbUrl)
    if err != nil {
        log.Fatal(err)
    }
    db.DB()
    //Remove this when going production
    fmt.Println("Dropping existing tables...")
    db.DropTableIfExists(&User{}, &Contact{}, &Skill{}, &School{}, &Job{})
    fmt.Println("Recreating tables...")
    db.AutoMigrate(&User{}, &Contact{}, &Skill{}, &School{}, &Job{})
    return
}

func loadMe() (me User) {
    data, err := ioutil.ReadFile("peternagy.json")
    if err != nil {
        fmt.Println(err)
        return
    }
    me = User{}
    json.Unmarshal([]byte(data), &me)
    return
}
