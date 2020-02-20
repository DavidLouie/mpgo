package main

import (
    "fmt"

    "github.com/davidlouie/mpgo/database"
    "github.com/davidlouie/mpgo/scanner"
    // "github.com/davidlouie/mpgo/server"
)

func main() {
    db := database.Init()
    defer db.Close()
    fmt.Println("Scanning files")
    scanner.Scan(db)
    fmt.Println("Scanning done")
    // server.Play()
    // Init()
}
