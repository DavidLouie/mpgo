package database

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "time"

    _ "github.com/mattn/go-sqlite3"
)

const dbPath = "./database/mpgo.db"
const datetimeId = 0

var db *sql.DB

func Init() *sql.DB {
    os.Remove(dbPath) // TODO Remove
    _, err := os.Stat(dbPath)
    dbNotExists := os.IsNotExist(err)
    db, err = sql.Open("sqlite3", dbPath)
    if err != nil {
        log.Fatal(err)
    }

    // If database doesn't exist yet, create tables
    if dbNotExists {
        fmt.Println("creating database")
        createArtistsTable()
        createAlbumsTable()
        createSongsTable()
        createDatetimeTable()
        createDirectoryTable()
    }

    AddArtist = initAddArtist()
    AddAlbum = initAddAlbum()
    AddSong = initAddSong()
    AddDirectory = initAddDirectory()
    GetDirectoryId = initGetDirectoryId()
    return db
}

func createArtistsTable() {
    sqlStmt := `CREATE TABLE Artists (
        artistId   INTEGER PRIMARY KEY,
        artistName TEXT UNIQUE
    )`
    _, err := db.Exec(sqlStmt)
    if err != nil {
        log.Fatal(err)
    }
}

func createAlbumsTable() {
    sqlStmt := `CREATE TABLE Albums (
        albumId     INTEGER PRIMARY KEY,
        albumTitle  TEXT,
        genre       TEXT,
        date        INTEGER,
        artistId    INTEGER NOT NULL,
        directoryId INTEGER NOT NULL,
        FOREIGN KEY (artistId) REFERENCES Artists,
        FOREIGN KEY (directoryId) REFERENCES Directory
    )`
    _, err := db.Exec(sqlStmt)
    if err != nil {
        log.Fatal(err)
    }
}

func createSongsTable() {
    sqlStmt := `CREATE TABLE Songs (
        songId      INTEGER PRIMARY KEY,
        songTitle   TEXT,
        duration    INTEGER,
        size        INTEGER,
        number      INTEGER,
        path        TEXT,
        bitrate     INTEGER,
        ext         INTEGER,
        albumId     INTEGER NOT NULL,
        directoryId INTEGER NOT NULL,
        FOREIGN KEY (albumId) REFERENCES Albums,
        FOREIGN KEY (directoryId) REFERENCES Directory
    )`
    _, err := db.Exec(sqlStmt)
    if err != nil {
        log.Fatal(err)
    }
}

func createDatetimeTable() {
    sqlStmt := `CREATE TABLE Datetime (
        id INTEGER PRIMARY KEY CHECK (id = 0),
        dt TEXT
    )`
    _, err := db.Exec(sqlStmt)
    if err != nil {
        log.Fatal(err)
    }
}

func createDirectoryTable() {
    sqlStmt := `CREATE TABLE Directory (
        directoryId   INTEGER PRIMARY KEY,
        directoryPath TEXT
    )`
    _, err := db.Exec(sqlStmt)
    if err != nil {
        log.Fatal(err)
    }
}

// Inserts or updates the current time into Datetime
// format is YYYY-MM-DD HH:MM:SS
func SetLastScannedTime() {
    sqlStmt := `INSERT INTO Datetime (id, dt) VALUES(?, datetime('now'))
        ON CONFLICT (id) DO UPDATE SET
        dt = excluded.dt`
    _, err := db.Exec(sqlStmt, datetimeId)
    if err != nil {
        log.Fatal(err)
    }
}

func GetLastScannedTime() time.Time {
    sqlStmt := "SELECT * FROM Datetime"
    row := db.QueryRow(sqlStmt, datetimeId)
    var dtStr string
    var id int
    err := row.Scan(&id, &dtStr)
    if err != nil {
        log.Fatal(err)
    }

    // Assumes dtStr is given by a timestamp with format YYYY-MM-DD HH:MM:SS
    fmt.Println(dtStr)
    dt, err := time.Parse("2006-01-02 15:04:05", dtStr)
    if err != nil {
        log.Fatal(err)
    }
    return dt
}

// Creates the AddArtist func, preparing statements in the closure
func initAddArtist() func(string) int {
    insStmt, err := db.Prepare(`
        INSERT INTO Artists(artistId, artistName)
        SELECT NULL, ?
        WHERE NOT EXISTS (SELECT * FROM Artists
                          WHERE artistName = ?)
    `)
    if err != nil {
        log.Fatal(err)
    }

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

var AddDirectory func(directoryPath string) int

func initAddDirectory() func(string) int {
    insStmt, err := db.Prepare(`
        INSERT INTO Directory (directoryId, directoryPath)
        SELECT NULL, ?
        WHERE NOT EXISTS (SELECT * FROM Directory
                          WHERE directoryPath = ?)
    `)
    if err != nil {
        log.Fatal(err)
    }

    qStmt, err := db.Prepare("SELECT directoryId FROM Directory WHERE directoryPath = ?")
    if err != nil {
        log.Fatal(err)
    }
    return func(directoryPath string) int {
        _, err = insStmt.Exec(directoryPath, directoryPath)
        if err != nil {
            log.Fatal(err)
        }

        row := qStmt.QueryRow(directoryPath)
        var directoryId int
        err = row.Scan(&directoryId)
        if err != nil {
            log.Fatal(err)
        }
        return directoryId
    }
}

var GetDirectoryId func (directoryPath string) int

func initGetDirectoryId() func(string) int {
    stmt, err := db.Prepare(`
        SELECT directoryId FROM Directory
        WHERE directoryPath = ?
    `)
    if err != nil {
        log.Fatal(err)
    }

    return func(directoryPath string) int {
        row := stmt.QueryRow(directoryPath)
        var directoryId int
        err = row.Scan(&directoryId)
        if err != nil {
            log.Fatal(err)
        }
        return directoryId
    }
}

// Adds the artist with given name to database, returning the ID
var AddArtist func (artistName string) int

// Creates the AddAlbum func, preparing statements in the closure
func initAddAlbum() func(string, string, int, int, string) int {
    insStmt, err := db.Prepare(`
        INSERT INTO Albums(albumId, albumTitle, genre, date, artistId, directoryId)
        SELECT NULL, ?, ?, ?, ?, ?
        WHERE NOT EXISTS (SELECT * FROM Albums
                          WHERE albumTitle = ?
                          AND artistId = ?)
    `)
    if err != nil {
        log.Fatal(err)
    }

    qStmt, err := db.Prepare("SELECT albumId FROM Albums WHERE albumTitle = ? AND artistId = ?")
    if err != nil {
        log.Fatal(err)
    }
    return func(albumTitle string, genre string, date int, artistId int, directoryPath string) int {
        directoryId := GetDirectoryId(directoryPath)
        _, err = insStmt.Exec(albumTitle, genre, date, artistId, directoryId, albumTitle, artistId)
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
var AddAlbum func (albumTitle string, genre string, date int, artistId int, directoryPath string) int

// Creates the AddSong func, preparing statements in the closure
func initAddSong() func(string, int, int64, int, string, int64, string, int, string) {
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
            albumId,
            directoryId)
        VALUES (NULL, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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
        directoryPath string,
    ) {
        directoryId := GetDirectoryId(directoryPath)
        _, err = stmt.Exec(
            songTitle,
            duration,
            size,
            trackNo,
            path,
            bitrate,
            ext,
            albumId,
            directoryId,
        )
        if err != nil {
            log.Fatal(err)
        }
    }
}

var AddSong func(
        songTitle     string,
        duration      int,
        size          int64,
        trackNo       int,
        path          string,
        bitrate       int64,
        ext           string,
        albumId       int,
        directoryPath string,
)

func PrintArtists() {
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

func PrintAlbums() {
    rows, err := db.Query("SELECT albumTitle, artistId, genre FROM Albums")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var (
        albumTitle string
        artistId int
        genre string
    )
    for rows.Next() {
        err := rows.Scan(&albumTitle, &artistId, &genre)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Println(albumTitle, artistId, genre)
    }
    err = rows.Err()
    if err != nil {
        log.Fatal(err)
    }
}

func PrintSongs() {
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
