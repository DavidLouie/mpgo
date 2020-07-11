package database

import (
	"database/sql"
	"log"
)

var folders = []string{"Music"}

// GetMusicFolders returns the music folders added to the database.
func GetMusicFolders() []string {
	return folders
}

// GenreCounts is a map from genre to the number of albums and songs in the genre
type GenreCounts map[string]counts

type counts struct {
	SongCount  int
	AlbumCount int
}

// GetGenres returns a map of {genre: {songCount, albumCount}}.
func GetGenres() GenreCounts {
	rows, err := db.Query(`
        SELECT genre, COUNT(DISTINCT A.albumID), COUNT(S.songID)
        FROM Albums A, Songs S
        WHERE A.albumID = S.albumID
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
