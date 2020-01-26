package subsonic

import (
    "log"
    "net/http"
)

type handler struct{}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(`{"message": "Hello world"}`))
}

func Init() {
    h := &handler{}
    http.Handle("/", h)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
