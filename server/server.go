package server

import (
    "fmt"
    "path/filepath"
    "os"
    "time"

    "github.com/davidlouie/mpgo/server/subsonic"
    "github.com/faiface/beep"
    "github.com/faiface/beep/effects"
    "github.com/faiface/beep/mp3"
    "github.com/faiface/beep/speaker"
)

const sampleRate = 44100
var queue Queue

func GetFiles(root string) ([]string, error) {
    ext := "*.mp3"
    var matches []string
    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            return nil
        }
        if matched, err := filepath.Match(ext, filepath.Base(path)); err != nil {
            return err
        } else if matched {
            matches = append(matches, path)
        }
        return nil
    })
    if err != nil {
        return nil, err
    }
    return matches, nil
}

func Play() {
    subsonic.Init()
    sr := beep.SampleRate(sampleRate)
    speaker.Init(sr, sr.N(time.Second / 10))
    speaker.Play(&queue)
}

func Add(song string) {
    f, err := os.Open(song)
    if err != nil {
        fmt.Println(err)
    }

    streamer, format, err := mp3.Decode(f)
    if err != nil {
        fmt.Println(err)
    }

    // set the volume of the streamer
    volume := &effects.Volume{
        Streamer:   streamer,
        Base:       2,
        Volume:     -4,
        Silent:     false,
    }

    // we fixed speaker's sample rate,
    // so need to resample file in case it doesn't match
    sr := beep.SampleRate(sampleRate)
    resampled := beep.Resample(4, format.SampleRate, sr, volume)

    speaker.Lock()
    queue.Add(resampled)
    speaker.Unlock()
}
