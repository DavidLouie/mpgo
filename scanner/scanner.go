package scanner

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/davidlouie/mpgo/database"
	"github.com/dhowden/tag"
)

const folder = "/home/david/Documents/Projects/"

// const folder = "/shared/david/Music/"
var audioExts = map[string]struct{}{
	".flac": {},
	".mp3":  {},
	".ogg":  {},
}

// Scan recursively searches for music files starting at folder.
// Assumes folder structure is: {folder}/{Artist}/{Album}/{Song}.
func Scan() {
	database.SetLastScannedTime()
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fmt.Println(path)

		if info.IsDir() {
			database.AddDirectory(path)
			return nil
		}
		return scanFile(path, info, err)
	})
	if err != nil {
		log.Fatal(err)
	}
}

// ScanNewFiles searches for new music files created or modified since lastScanned
func ScanNewFiles() {
	lastScanned := database.GetLastScannedTime()
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if info.ModTime().After(lastScanned) {
			return scanFile(path, info, err)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

// Scans given file and adds it into the library database based on tags
func scanFile(path string, info os.FileInfo, err error) error {
	fileExt := filepath.Ext(path)
	// fmt.Println(filepath.Base(path))
	if _, ok := audioExts[fileExt]; ok {
		f, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}

		m, err := tag.ReadFrom(f)
		if err != nil {
			return nil
		}

		duration, err := getDuration(f, fileExt)
		if err != nil {
			return nil
		}

		dir := filepath.Dir(path)
		artistID := database.AddArtist(m.Artist(), dir)
		albumID := database.AddAlbum(m.Album(), m.Genre(), m.Year(), artistID, dir)

		trackNo, _ := m.Track()
		size := info.Size()
		bitrate := size / int64(duration)
		database.AddSong(
			m.Title(),
			duration,
			size,
			trackNo,
			path,
			bitrate,
			fileExt,
			albumID,
			dir,
		)
	}
	return nil
}

// Returns the duration of the given flac/mp3/ogg audio file
func getDuration(f *os.File, ext string) (int, error) {
	// TODO
	return 300, nil
}
