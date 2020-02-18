package subsonic

import (
    "encoding/xml"
    "fmt"
    "net/http"
    "strconv"

    "github.com/davidlouie/mpgo/query"
)

func GetMusicFolders(w http.ResponseWriter, r *http.Request) {
    params, err := parseParams(w, r)
    if err != nil {
        return
    }
    err = authenticate(w, params)
    if err != nil {
        return
    }

    f := query.GetMusicFolders()
    response := &subResp{Status: "ok", Version: apiVersion}
    foldersArr := make([]musicFolder, len(f))
    for i, folder := range f{
        foldersArr[i].Id = strconv.Itoa(i)
        foldersArr[i].Name = folder
    }

    folders := &musicFolders{Folders: foldersArr}
    response.MusicFolders = folders
    encoded, err := xml.MarshalIndent(response, "  ", "    ")
    if err != nil {
        fmt.Println(err)
    }

    w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-Type", "application/xml")
    w.Write(encoded)
}
