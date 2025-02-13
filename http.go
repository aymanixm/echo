package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    _ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USERNAME")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")
    dbSchema := os.Getenv("DB_SCHEMA")

    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s search_path=%s sslmode=disable",
        dbHost, dbPort, dbUser, dbPassword, dbName, dbSchema)

    var err error
    db, err = sql.Open("postgres", dsn)
    if err != nil {
        log.Fatalf("Failed to connect to the database: %v", err)
    }

    if err = db.Ping(); err != nil {
        log.Fatalf("Database is not reachable: %v", err)
    }
    fmt.Println("Connected to the database successfully!")
}

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    initDB()

    fmt.Printf("Listening on :%s\n", port)

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            queryParams := r.URL.Query()
            if len(queryParams) == 0 {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusOK)
                w.Write([]byte("{}"))
            } else {
                response := make(map[string]interface{})
                for key, values := range queryParams {
                    if len(values) > 1 {
                        response[key] = values
                    } else {
                        response[key] = values[0]
                    }
                }
                jsonResponse, _ := json.Marshal(response)
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusOK)
                w.Write(jsonResponse)
            }
        case http.MethodPost:
            body, err := io.ReadAll(r.Body)
            if err != nil {
                http.Error(w, "Cannot read body", http.StatusBadRequest)
                return
            }
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusOK)
            w.Write(body)
        default:
            http.NotFound(w, r)
        }
    })

    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            var healthCheck int
            err := db.QueryRow("SELECT 1").Scan(&healthCheck)
            if err != nil {
                w.WriteHeader(http.StatusInternalServerError)
                w.Write([]byte("DB connection failed"))
                return
            }
            w.WriteHeader(http.StatusOK)
            w.Write([]byte("OK"))
        } else {
            http.NotFound(w, r)
        }
    })

    log.Fatal(http.ListenAndServe(":"+port, nil))
}
