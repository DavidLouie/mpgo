package database

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    "github.com/davidlouie/mpgo/scanner"
    _ "github.com/mattn/go-sqlite3"
)

const dbPath = "./database/mpgo.db"

func Init() {
    if _, err := os.Stat(dbPath); os.IsNotExist(err) {
        fmt.Println("creating database")
        db, err := sql.Open("sqlite3", dbPath)
        if err != nil {
            log.Fatal(err)
        }
        defer db.Close()

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
}

func Scan() {
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    scanner.Scan(db)
}

func AddArtist(db *sql.DB, artistName string) int64 {
    stmt, err := db.Prepare(`
        INSERT INTO Artists(artistId, artistName)
        SELECT NULL, ?
        WHERE NOT EXISTS (SELECT * FROM Artists
                          WHERE artistName = ?)
    `)
    if err != nil {
        log.Fatal(err)
    }

    res, err := stmt.Exec(artistName, artistName)
    if err != nil {
        log.Fatal(err)
    }

    lastId, err := res.LastInsertId()
    if err != nil {
        log.Fatal(err)
    }
    return lastId
}

func AddAlbum(db *sql.DB, albumTitle string, genre string, date int, artistId int64) int64 {
    stmt, err := db.Prepare(`
        INSERT INTO Albums(albumId, albumTitle, genre, date, artistId)
        SELECT NULL, ?, ?, ?, ?
        WHERE NOT EXISTS (SELECT * FROM Albums
                          WHERE albumName = ?
                          AND artistId = ?)
    `)
    if err != nil {
        log.Fatal(err)
    }

    res, err := stmt.Exec(albumTitle, genre, date, artistId, albumTitle, artistId)
    if err != nil {
        log.Fatal(err)
    }

    lastId, err := res.LastInsertId()
    if err != nil {
        log.Fatal(err)
    }
    return lastId
}

func AddSong(
        db        *sql.DB,
        songTitle string,
        duration  int,
        size      int64,
        trackNo   int,
        path      string,
        bitrate   int64,
        ext       string,
        albumId   int64,
) {
    stmt, err := db.Prepare(`
        INSERT INTO Songs(
            songId,
            songTitle,
            duration,
            size,
            number
            path,
            bitrate,
            ext,
            albumId)
        VALUES (NULL, ?, ?, ?, ?, ?, ?, ?)
    `)
    if err != nil {
        log.Fatal(err)
    }

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
