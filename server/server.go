package server

import (
    "fmt"
    "log"
    "os"
    "time"

    "github.com/faiface/beep"
    "github.com/faiface/beep/mp3"
    "github.com/faiface/beep/speaker"
)

func Play(song string) {
    f, err := os.Open(song)
    if err != nil {
        log.Fatal(err)
    }

    streamer, format, err := mp3.Decode(f)
    if err != nil {
        log.Fatal(err)
    }
    defer streamer.Close()


    sr := format.SampleRate
    speaker.Init(sr, sr.N(time.Second / 10))

    ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamer), Paused: false}
    speaker.Play(ctrl)

    for {
        fmt.Print("Press [ENTER] to pause/resume.")
        fmt.Scanln()

        speaker.Lock()
        ctrl.Paused = !ctrl.Paused
        speaker.Unlock()
    }
}
