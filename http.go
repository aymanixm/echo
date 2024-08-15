package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    fmt.Fprintf(os.Stdout, "Listening on :%s\n", port)

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            queryParams := r.URL.Query()
            if len(queryParams) == 0 {
                // If no query parameters, return 200 OK with an empty JSON object
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusOK)
                w.Write([]byte("{}"))
            } else {
                // Return the query parameters as a JSON response
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
            body, err := ioutil.ReadAll(r.Body)
            if err != nil {
                http.Error(w, "Cannot read body", http.StatusBadRequest)
                return
            }
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusOK)
            w.Write(body)

        default:
            // For other methods, return 404
            http.NotFound(w, r)
        }
    })

    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            w.WriteHeader(http.StatusOK)
            w.Write([]byte("OK"))
        } else {
            // For other methods, return 404
            http.NotFound(w, r)
        }
    })

    log.Fatal(http.ListenAndServe(":"+port, nil))
}
