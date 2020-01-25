package main

import (
    "fmt"
    "log"

    "github.com/davidlouie/mpgo/server"
)

func main() {
    /*for {
        fmt.Print("Type the name of the song to play: ")
        scanner := bufio.NewScanner(os.Stdin)
        for scanner.Scan() {
            server.Play(scanner.Text())
        }
        if scanner.Err() != nil {
            log.Fatal(scanner.Err())
        }
    }*/
    // Init()
    filepaths, err := server.GetFiles("/home/david/Music/")
    if err != nil {
        log.Fatal(err)
    }
    for _, filepath := range filepaths {
        fmt.Println(filepath)
    }

    server.Play()
}
