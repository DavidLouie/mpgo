package database

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/mattn/go-sqlite3"
)

const dbPath = "./database/mpgo.db"

func Init() *sql.DB {
    os.Remove(dbPath) // TODO Remove
    _, err := os.Stat(dbPath)
    dbNotExists := os.IsNotExist(err)
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        log.Fatal(err)
    }

    // If database doesn't exist yet, create tables
    if dbNotExists {
        fmt.Println("creating database")

        sqlStmt := `CREATE TABLE Artists (
            artistId   INTEGER PRIMARY KEY,
            artistName VARCHAR(32) UNIQUE
        )`
        _, err = db.Exec(sqlStmt)
        if err != nil {
            log.Fatal(err)
        }

        sqlStmt = `CREATE TABLE Albums (
            albumId    INTEGER PRIMARY KEY,
            albumTitle VARCHAR(32),
            genre      VARCHAR(32),
            date       INTEGER,
            artistId   INTEGER NOT NULL,
            FOREIGN KEY (artistId) REFERENCES Artists
        )`
        _, err = db.Exec(sqlStmt)
        if err != nil {
            log.Fatal(err)
        }

        sqlStmt = `CREATE TABLE Songs (
            songId     INTEGER PRIMARY KEY,
            songTitle  VARCHAR(32),
            duration   INTEGER,
            size       INTEGER,
            number     INTEGER,
            path       VARCHAR(32),
            bitrate    INTEGER,
            ext        INTEGER,
            albumId    INTEGER NOT NULL,
            FOREIGN KEY (albumId) REFERENCES Albums
        )`
        _, err = db.Exec(sqlStmt)
        if err != nil {
            log.Fatal(err)
        }
    }

    AddArtist = initAddArtist(db)
    AddAlbum = initAddAlbum(db)
    AddSong = initAddSong(db)
    return db
}

// Creates the AddArtist func, preparing statements in the closure
func initAddArtist(db *sql.DB) func(string) int {
    // insert if artist doesn't already exist
    insStmt, err := db.Prepare(`
        INSERT INTO Artists(artistId, artistName)
        SELECT NULL, ?
        WHERE NOT EXISTS (SELECT * FROM Artists
                          WHERE artistName = ?)
    `)
    if err != nil {
        log.Fatal(err)
    }

    // get the artistId of the artist we just inserted
    qStmt, err := db.Prepare("SELECT artistId FROM Artists WHERE artistName = ?")
    if err != nil {
        log.Fatal(err)
    }
    return func(artistName string) int {
        _, err = insStmt.Exec(artistName, artistName)
        if err != nil {
            log.Fatal(err)
        }

        row := qStmt.QueryRow(artistName)
        var artistId int
        err = row.Scan(&artistId)
        if err != nil {
            log.Fatal(err)
        }
        return artistId
    }
}

// Adds the artist with given name to database, returning the ID
var AddArtist func (artistName string) int

// Creates the AddAlbum func, preparing statements in the closure
func initAddAlbum(db *sql.DB) func(string, string, int, int) int {
    // insert if album doesn't already exist
    insStmt, err := db.Prepare(`
        INSERT INTO Albums(albumId, albumTitle, genre, date, artistId)
        SELECT NULL, ?, ?, ?, ?
        WHERE NOT EXISTS (SELECT * FROM Albums
                          WHERE albumTitle = ?
                          AND artistId = ?)
    `)
    if err != nil {
        log.Fatal(err)
    }

    // get albumId of the album we just inserted
    qStmt, err := db.Prepare("SELECT albumId FROM Albums WHERE albumTitle = ? AND artistId = ?")
    if err != nil {
        log.Fatal(err)
    }
    return func(albumTitle string, genre string, date int, artistId int) int {
        _, err = insStmt.Exec(albumTitle, genre, date, artistId, albumTitle, artistId)
        if err != nil {
            log.Fatal(err)
        }

        row := qStmt.QueryRow(albumTitle, artistId)
        var albumId int
        err = row.Scan(&albumId)
        if err != nil {
            log.Fatal(err)
        }
        return albumId
    }
}

// Adds the album with given metadata to database, returning the ID
var AddAlbum func (albumTitle string, genre string, date int, artistId int) int

// Creates the AddSong func, preparing statements in the closure
func initAddSong(db *sql.DB) func(string, int, int64, int, string, int64, string, int) {
    stmt, err := db.Prepare(`
        INSERT INTO Songs(
            songId,
            songTitle,
            duration,
            size,
            number,
            path,
            bitrate,
            ext,
            albumId)
        VALUES (NULL, ?, ?, ?, ?, ?, ?, ?, ?)
    `)
    if err != nil {
        log.Fatal(err)
    }

    return func(
        songTitle string,
        duration int,
        size int64,
        trackNo int,
        path string,
        bitrate int64,
        ext string,
        albumId int,
    ) {
        _, err = stmt.Exec(
            songTitle,
            duration,
            size,
            trackNo,
            path,
            bitrate,
            ext,
            albumId,
        )
        if err != nil {
            log.Fatal(err)
        }
    }
}

var AddSong func(
        songTitle string,
        duration  int,
        size      int64,
        trackNo   int,
        path      string,
        bitrate   int64,
        ext       string,
        albumId   int,
)

func PrintArtists(db *sql.DB) {
    rows, err := db.Query("SELECT * FROM Artists")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var (
        artistId int
        artistName string
    )
    for rows.Next() {
        err := rows.Scan(&artistId, &artistName)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Println(artistId, artistName)
    }
    err = rows.Err()
    if err != nil {
        log.Fatal(err)
    }
}

func PrintAlbums(db *sql.DB) {
    rows, err := db.Query("SELECT albumTitle, artistId FROM Albums")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var (
        albumTitle string
        artistId int
    )
    for rows.Next() {
        err := rows.Scan(&albumTitle, &artistId)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Println(albumTitle, artistId)
    }
    err = rows.Err()
    if err != nil {
        log.Fatal(err)
    }
}

func PrintSongs(db *sql.DB) {
    rows, err := db.Query("SELECT songTitle, albumId FROM Songs")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var (
        songTitle string
        albumId int
    )
    for rows.Next() {
        err := rows.Scan(
            &songTitle,
            &albumId)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Println(songTitle, albumId)
    }
    err = rows.Err()
    if err != nil {
        log.Fatal(err)
    }
}
