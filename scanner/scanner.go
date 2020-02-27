package scanner

import (
    // "fmt"
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
// const folder = "/shared/david/Music/"
var audioExts = map[string]struct{}{
    ".flac": struct{}{},
    ".mp3": struct{}{},
    ".ogg": struct{}{},
}

// Scans for music files starting at folder
// Assumes folder structure is: {folder}/{Artist}/{Album}/{Song}
func Scan() {
    database.SetLastScannedTime()
    err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        if info.IsDir() {
            return nil 
        }
        return scanFile(path, info, err)
    })
    if err != nil {
        log.Fatal(err)
    }
}

// Scans for new music files created or modified since lastScanned
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

        artistId := database.AddArtist(m.Artist())
        albumId := database.AddAlbum(m.Album(), m.Genre(), m.Year(), artistId)

        trackNo, _ := m.Track()
        size :=  info.Size()
        bitrate := size / int64(duration)
        database.AddSong(
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

// Returns the duration of the given flac/mp3/ogg audio file
func getDuration(f *os.File, ext string) (int, error) {
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
        return 0, err
    }
    return streamer.Len(), nil
}

