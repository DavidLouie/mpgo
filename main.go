package main

import (
    "fmt"
    //"log"

    "github.com/davidlouie/mpgo/database"
    "github.com/davidlouie/mpgo/server"
)

func main() {
    /*filepaths, err := server.GetFiles("/home/david/Music/")
    if err != nil {
        log.Fatal(err)
    }
    for _, filepath := range filepaths {
        fmt.Println(filepath)
    }*/

    fmt.Println("about to call db.Init()")
    database.Init()
    database.Scan()
    server.Play()
    Init()
}
