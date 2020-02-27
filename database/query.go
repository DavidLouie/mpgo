package database

import (
    "log"
)

const ROOT = "/home/david/"
var FOLDERS = []string{"Music"}

func GetMusicFolders() []string {
    return FOLDERS
}

type counts struct {
    SongCount  int
    AlbumCount int
}

// Returns a map of {genre: {songCount, albumCount}}
func GetGenres() map[string]counts {
    genreMap := make(map[string]counts)
    rows, err := db.Query(`
        SELECT genre, COUNT(albumId) FROM Albums
        GROUP BY genre
    `)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    var (
        genre string
        count int
    )
    for rows.Next() {
        err := rows.Scan(
            &genre,
            &count)
        if err != nil {
            log.Fatal(err)
        }
        genreMap[genre] = counts{AlbumCount: count, SongCount: 0}
    }
    err = rows.Err()
    if err != nil {
        log.Fatal(err)
    }

    rows, err = db.Query(`
        SELECT genre, COUNT(songId)
        FROM Albums JOIN Songs ON Albums.albumId = Songs.albumId
        GROUP BY genre
    `)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    for rows.Next() {
        err := rows.Scan(
            &genre,
            &count)
        if err != nil {
            log.Fatal(err)
        }
        currCounts := genreMap[genre]
        genreMap[genre] = counts{AlbumCount: currCounts.AlbumCount, SongCount: count}
    }
    err = rows.Err()
    if err != nil {
        log.Fatal(err)
    }
    return genreMap
}
