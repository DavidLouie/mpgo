package scanner

import (
    "database/sql"
    "log"
    "os"
    "path/filepath"
    
    "github.com/davidlouie/mpgo/database"
    "github.com/dhowden/tag"
    "github.com/faiface/beep"
    "github.com/faiface/beep/flac"
    "github.com/faiface/beep/mp3"
    "github.com/faiface/beep/vorbis"
    _ "github.com/mattn/go-sqlite3"
)

const folder = "/home/david/Music/"
var audioExts = map[string]struct{}{
    ".flac": struct{}{},
    ".mp3": struct{}{},
    ".ogg": struct{}{},
}

// Scans for music files starting at folder
// Assumes folder structure is: {folder}/{Artist}/{Album}/{Song}
func Scan(db *sql.DB) {
    err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        if info.IsDir() {
            return nil 
        }
        return scanFile(db, path, info, err)
    })
    if err != nil {
        log.Fatal(err)
    }
}

func scanFile(db *sql.DB, path string, info os.FileInfo, err error) error {
    fileExt := filepath.Ext(path)
    if _, ok := audioExts[fileExt]; ok {
        f, err := os.Open(path)
        if err != nil {
            log.Fatal(err)
        }

        m, err := tag.ReadFrom(f)
        if err != nil {
            log.Fatal(err)
        }
        
        artistId := database.AddArtist(db, m.Artist())
        albumId := database.AddAlbum(db, m.Artist(), m.Genre(), m.Year(), artistId)

        trackNo, _ := m.Track()
        duration := getDuration(f, fileExt)
        size :=  info.Size()
        bitrate := size / int64(duration)
        database.AddSong(
            db,
            m.Title(),
            duration,
            size,
            trackNo,
            path,
            bitrate,
            fileExt,
            albumId,
        )

    }
    return nil
}

func getDuration(f *os.File, ext string) int {
    var streamer beep.StreamSeeker
    var err error
    switch ext {
    case ".flac":
        streamer, _, err = flac.Decode(f)
    case ".mp3":
        streamer, _, err = mp3.Decode(f)
    case ".ogg":
        streamer, _, err = vorbis.Decode(f)
    default:
        log.Fatalf("Wrong filetype %s\n", ext)
    }
    if err != nil {
        log.Fatal(err)
    }
    return streamer.Len()
}

