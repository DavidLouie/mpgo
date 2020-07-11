package database

import (
	"database/sql"
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
	rows, err := db.Query(`
        SELECT genre, COUNT(DISTINCT albumId), COUNT(songId)
        FROM Albums, Songs
        WHERE Albums.albumId = Songs.albumId
        GROUP BY genre
    `)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	return getGenreMapFromRows(rows)
}

func getGenreMapFromRows(rows *sql.Rows) map[string]counts {
	genreMap := make(map[string]counts)
	var (
		genre      string
		albumCount int
		songCount  int
	)
	for rows.Next() {
		err := rows.Scan(
			&genre,
			&albumCount,
			&songCount,
		)
		if err != nil {
			log.Fatal(err)
		}
		genreMap[genre] = counts{AlbumCount: albumCount, SongCount: songCount}
	}

	err := rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return genreMap
}
