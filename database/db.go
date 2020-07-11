package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/davidlouie/mpgo/util"
	// blank import to register driver with database/sql
	_ "github.com/mattn/go-sqlite3"
)

const dbPath = "./database/mpgo.db"
const datetimeID = 0

var db *sql.DB

// Init creates the database tables and prepares access sql statements.
// Returns the opened db.
func Init() *sql.DB {
	os.Remove(dbPath) // TODO Remove
	_, err := os.Stat(dbPath)
	dbNotExists := os.IsNotExist(err)
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

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

	GetDirectoryID = initGetDirectoryID()
	GetArtistsFromDirID = initGetArtistsFromDirID()
	GetAlbumsFromDirID = initGetAlbumsFromDirID()
	GetSongsFromDirID = initGetSongsFromDirID()
	return db
}

func createArtistsTable() {
	sqlStmt := `CREATE TABLE Artists (
        artistID    INTEGER PRIMARY KEY,
        artistName  TEXT    UNIQUE,
        directoryID INTEGER,
        FOREIGN KEY (directoryID) REFERENCES Directory
    )`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
}

func createAlbumsTable() {
	sqlStmt := `CREATE TABLE Albums (
        albumID     INTEGER PRIMARY KEY,
        albumTitle  TEXT,
        genre       TEXT,
        year        INTEGER,
        artistID    INTEGER NOT NULL,
        directoryID INTEGER NOT NULL,
        FOREIGN KEY (artistID) REFERENCES Artists,
        FOREIGN KEY (directoryID) REFERENCES Directory
    )`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
}

func createSongsTable() {
	sqlStmt := `CREATE TABLE Songs (
        songID      INTEGER PRIMARY KEY,
        songTitle   TEXT,
        duration    INTEGER,
        size        INTEGER,
        track       INTEGER,
        path        TEXT,
        bitrate     INTEGER,
        ext         INTEGER,
        albumID     INTEGER NOT NULL,
        directoryID INTEGER NOT NULL,
        FOREIGN KEY (albumID) REFERENCES Albums,
        FOREIGN KEY (directoryID) REFERENCES Directory
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
        directoryID   INTEGER PRIMARY KEY,
        directoryPath TEXT
    )`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
}

// SetLastScannedTime inserts or updates the current time into Datetime.
// The time format is YYYY-MM-DD HH:MM:SS.
func SetLastScannedTime() {
	sqlStmt := `INSERT INTO Datetime (id, dt) VALUES(?, datetime('now'))
        ON CONFLICT (id) DO UPDATE SET
        dt = excluded.dt`
	_, err := db.Exec(sqlStmt, datetimeID)
	if err != nil {
		log.Fatal(err)
	}
}

// GetLastScannedTime returns the time the directory was last scanned.
func GetLastScannedTime() time.Time {
	sqlStmt := "SELECT * FROM Datetime"
	row := db.QueryRow(sqlStmt, datetimeID)
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
func initAddArtist() func(string, string) int {
	// insert if artist doesn't already exist
	insStmt, err := db.Prepare(`
        INSERT INTO Artists(artistID, artistName, directoryID)
        SELECT NULL, ?, ?
        WHERE NOT EXISTS (SELECT * FROM Artists
                          WHERE artistName = ?)
    `)
	if err != nil {
		log.Fatal(err)
	}

	qStmt, err := db.Prepare("SELECT artistID FROM Artists WHERE artistName = ?")
	if err != nil {
		log.Fatal(err)
	}
	return func(artistName string, directoryPath string) int {
		directoryID := GetDirectoryID(directoryPath)
		_, err = insStmt.Exec(artistName, directoryID, artistName)
		if err != nil {
			log.Fatal(err)
		}
		row := qStmt.QueryRow(artistName)
		var artistID int
		err = row.Scan(&artistID)
		if err != nil {
			log.Fatal(err)
		}
		return artistID
	}
}

// AddDirectory adds the directory path to the database.
// Returns the corresponding database directoryID.
var AddDirectory func(directoryPath string) int

func initAddDirectory() func(string) int {
	insStmt, err := db.Prepare(`
        INSERT INTO Directory (directoryID, directoryPath)
        SELECT NULL, ?
        WHERE NOT EXISTS (SELECT * FROM Directory
                          WHERE directoryPath = ?)
    `)
	if err != nil {
		log.Fatal(err)
	}

	qStmt, err := db.Prepare("SELECT directoryID FROM Directory WHERE directoryPath = ?")
	if err != nil {
		log.Fatal(err)
	}
	return func(directoryPath string) int {
		_, err = insStmt.Exec(directoryPath, directoryPath)
		if err != nil {
			log.Fatal(err)
		}

		row := qStmt.QueryRow(directoryPath)
		var directoryID int
		err = row.Scan(&directoryID)
		if err != nil {
			log.Fatal(err)
		}
		return directoryID
	}
}

// GetDirectoryID returns the directoryID corresponding to the given path.
var GetDirectoryID func(directoryPath string) int

func initGetDirectoryID() func(string) int {
	stmt, err := db.Prepare(`
        SELECT directoryID FROM Directory
        WHERE directoryPath = ?
    `)
	if err != nil {
		log.Fatal(err)
	}

	return func(directoryPath string) int {
		row := stmt.QueryRow(directoryPath)
		var directoryID int
		err = row.Scan(&directoryID)
		if err != nil {
			log.Fatal(err)
		}
		return directoryID
	}
}

// AddArtist adds the artist with given name to database.
// Returns the created artistID.
var AddArtist func(artistName string, directoryPath string) int

func initAddAlbum() func(string, string, int, int, string) int {
	insStmt, err := db.Prepare(`
        INSERT INTO Albums(albumID, albumTitle, genre, year, artistID, directoryID)
        SELECT NULL, ?, ?, ?, ?, ?
        WHERE NOT EXISTS (SELECT * FROM Albums
                          WHERE albumTitle = ?
                          AND artistID = ?)
    `)
	if err != nil {
		log.Fatal(err)
	}

	qStmt, err := db.Prepare("SELECT albumID FROM Albums WHERE albumTitle = ? AND artistID = ?")
	if err != nil {
		log.Fatal(err)
	}
	return func(albumTitle string, genre string, year int, artistID int, directoryPath string) int {
		directoryID := GetDirectoryID(directoryPath)
		_, err = insStmt.Exec(albumTitle, genre, year, artistID, directoryID, albumTitle, artistID)
		if err != nil {
			log.Fatal(err)
		}

		row := qStmt.QueryRow(albumTitle, artistID)
		var albumID int
		err = row.Scan(&albumID)
		if err != nil {
			log.Fatal(err)
		}
		return albumID
	}
}

// AddAlbum adds the album with given metadata to database.
// Returns the created albumID.
var AddAlbum func(albumTitle string, genre string, year int, artistID int, directoryPath string) int

func initAddSong() func(string, int, int64, int, string, int64, string, int, string) {
	stmt, err := db.Prepare(`
        INSERT INTO Songs(
            songID,
            songTitle,
            duration,
            size,
            track,
            path,
            bitrate,
            ext,
            albumID,
            directoryID)
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
		albumID int,
		directoryPath string,
	) {
		directoryID := GetDirectoryID(directoryPath)
		_, err = stmt.Exec(
			songTitle,
			duration,
			size,
			trackNo,
			path,
			bitrate,
			ext,
			albumID,
			directoryID,
		)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// AddSong adds the song with the given metadata to the database.
var AddSong func(
	songTitle string,
	duration int,
	size int64,
	trackNo int,
	path string,
	bitrate int64,
	ext string,
	albumID int,
	directoryPath string,
)

// GetArtistsFromDirID returns array of artists in the given directory.
var GetArtistsFromDirID func(directoryID int) []util.Artist

func initGetArtistsFromDirID() func(int) []util.Artist {
	qStmt, err := db.Prepare(`
        SELECT artistID, artistName FROM Artists INNER JOIN Directory
        ON Artists.directoryID = Directory.directoryID
        WHERE Artists.directoryID = ?`)
	if err != nil {
		log.Fatal(err)
	}

	return func(directoryID int) []util.Artist {
		rows, err := qStmt.Query(directoryID)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var artists []util.Artist
		for rows.Next() {
			artist := util.NewArtist()
			err := rows.Scan(&artist.ID, &artist.Name)
			if err != nil {
				log.Fatal(err)
			}
			artists = append(artists, artist)
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
		return artists
	}
}

// GetAlbumsFromDirID returns array of albums in the given directory.
var GetAlbumsFromDirID func(directoryID int) []util.Album

func initGetAlbumsFromDirID() func(int) []util.Album {
	qStmt, err := db.Prepare(`
        SELECT albumID, albumTitle, genre, year, artistID
        FROM Albums INNER JOIN Directory
        ON Albums.directoryID = Directory.directoryID
        WHERE Albums.directoryID = ?`)
	if err != nil {
		log.Fatal(err)
	}

	return func(directoryID int) []util.Album {
		rows, err := qStmt.Query(directoryID)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var albums []util.Album
		for rows.Next() {
			album := util.NewAlbum()
			err := rows.Scan(
				&album.ID,
				&album.Title,
				&album.Genre,
				&album.Year,
				&album.ArtistID)
			if err != nil {
				log.Fatal(err)
			}
			albums = append(albums, album)
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
		return albums
	}
}

// GetSongsFromDirID returns array of songs in the given directory.
var GetSongsFromDirID func(directoryID int) []util.Song

func initGetSongsFromDirID() func(int) []util.Song {
	qStmt, err := db.Prepare(`
        SELECT songID, songTitle, duration, size, track, path, bitrate, ext, albumID
        FROM Songs INNER JOIN Directory
        ON Songs.directoryID = Directory.directoryID
        WHERE Songs.directoryID = ?`)
	if err != nil {
		log.Fatal(err)
	}

	return func(directoryID int) []util.Song {
		rows, err := qStmt.Query(directoryID)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var songs []util.Song
		for rows.Next() {
			song := util.NewSong()
			err := rows.Scan(
				&song.ID,
				&song.Title,
				&song.Duration,
				&song.Size,
				&song.Track,
				&song.Path,
				&song.BitRate,
				&song.Ext,
				&song.AlbumID)
			if err != nil {
				log.Fatal(err)
			}
			songs = append(songs, song)
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
		return songs
	}
}

// PrintArtists prints all the artists in the database.
func PrintArtists() {
	rows, err := db.Query("SELECT * FROM Artists")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var (
		artistID   int
		artistName string
	)
	for rows.Next() {
		err := rows.Scan(&artistID, &artistName)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(artistID, artistName)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

// PrintAlbums prints all the albums in the database.
func PrintAlbums() {
	rows, err := db.Query("SELECT albumTitle, artistID, genre FROM Albums")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var (
		albumTitle string
		artistID   int
		genre      string
	)
	for rows.Next() {
		err := rows.Scan(&albumTitle, &artistID, &genre)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(albumTitle, artistID, genre)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

// PrintSongs prints all the songs in the database.
func PrintSongs() {
	rows, err := db.Query("SELECT songTitle, albumID FROM Songs")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var (
		songTitle string
		albumID   int
	)
	for rows.Next() {
		err := rows.Scan(
			&songTitle,
			&albumID)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(songTitle, albumID)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
