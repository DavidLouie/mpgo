package server

import (
    "github.com/faiface/beep"
)

type Queue struct {
    streamers []beep.Streamer
}

func (q *Queue) Add(streamers ...beep.Streamer) {
    q.streamers = append(q.streamers, streamers...)
}

func (q *Queue) Stream(samples [][2]float64) (n int, ok bool) {
    // continue playing samples until we've filled all of them
    filled := 0
    for filled < len(samples) {
        if len(q.streamers) == 0 {
            for i := range samples[filled:] {
                samples[i][0] = 0
                samples[i][1] = 0
            }
            break
        }

        // stream from next streamer in queue
        n, ok := q.streamers[0].Stream(samples[filled:])

        // streamer drained, pop from queue and continue
        if !ok {
            q.streamers = q.streamers[1:]
        }
        filled += n
    }
    return len(samples), true
}

func (q *Queue) Err() error {
    return nil
}
