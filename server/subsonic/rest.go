package subsonic

import (
    "log"
    "net/http"
)

func Init() {
    http.HandleFunc("/rest/getMusicFolders", GetMusicFolders)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
