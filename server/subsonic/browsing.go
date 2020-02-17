package subsonic

import (
    "encoding/xml"
    "fmt"
    "net/http"
    "strconv"

    "github.com/davidlouie/mpgo/query"
)

func GetMusicFolders(w http.ResponseWriter, r *http.Request) {
    f := query.GetMusicFolders()
    response := &SubResp{Status: "ok", Version: "1.16.1"}
    foldersArr := make([]MusicFolder, len(f))
    for i, folder := range f{
        foldersArr[i].Id = strconv.Itoa(i)
        foldersArr[i].Name = folder
    }

    folders := &MusicFolders{Folders: foldersArr}
    response.MusicFolders = folders
    encoded, err := xml.MarshalIndent(response, "  ", "    ")
    if err != nil {
        fmt.Println(err)
    }

    w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-Type", "application/xml")
    w.Write(encoded)
}
