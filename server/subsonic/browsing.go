package subsonic

import (
    "encoding/xml"
    "fmt"
    "net/http"
    "strconv"

    "github.com/davidlouie/mpgo/database"
)

// Returns all configured top-level music folders
func getMusicFolders(w http.ResponseWriter, r *http.Request) {
    err := parseAndAuth(w, r)
    if err != nil {
        return
    }

    f := database.GetMusicFolders()
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

// Returns an indexed structure of all artists
func getIndexes(w http.ResponseWriter, r *http.Request) {
    err := parseAndAuth(w, r)
    if err != nil {
        return
    }

    //response := &subResp{Status: "ok", Version: apiVersion}
    //genresArr := make([]genre, 10)
}

// Returns a listing of all files in a music directory
func getMusicDirectory(w http.ResponseWriter, r *http.Request) {
    err := parseAndAuth(w, r)
    if err != nil {
        return
    }

    // TODO: not implemented yet
    sendError(w, 0)
}

// Returns all genres
func getGenres(w http.ResponseWriter, r *http.Request) {
    err := parseAndAuth(w, r)
    if err != nil {
        return
    }

    genreMap := database.GetGenres()
    response := &subResp{Status: "ok", Version: apiVersion}
    var genresArr []genre
    for k, v := range genreMap {
        genresArr = append(genresArr,
            genre{
                SongCount: strconv.Itoa(v.SongCount),
                AlbumCount: strconv.Itoa(v.AlbumCount),
                Genre: k},
        )
    }

    genres := &genres{Genres: genresArr}
    response.Genres = genres
    encoded, err := xml.MarshalIndent(response, "  ", "    ")
    if err != nil {
        fmt.Println(err)
    }

    w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-Type", "application/xml")
    w.Write(encoded)
}

// Returns an indexed structure of all artists organized by ID3 tags
func getArtists(w http.ResponseWriter, r *http.Request) {
    err := parseAndAuth(w, r)
    if err != nil {
        return
    }

    // TODO: not implemented yet
    sendError(w, 0)
}

// Returns details of an artist, including a list of albums
func getArtist(w http.ResponseWriter, r *http.Request) {
    err := parseAndAuth(w, r)
    if err != nil {
        return
    }

    // TODO: not implemented yet
    sendError(w, 0)
}

// Returns details for an album, including a list of songs
func getAlbum(w http.ResponseWriter, r *http.Request) {
    err := parseAndAuth(w, r)
    if err != nil {
        return
    }

    // TODO: not implemented yet
    sendError(w, 0)
}

// Returns details for a song
func getSong(w http.ResponseWriter, r *http.Request) {
    err := parseAndAuth(w, r)
    if err != nil {
        return
    }

    // TODO: not implemented yet
    sendError(w, 0)
}

