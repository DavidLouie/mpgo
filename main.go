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
    fmt.Println("Scanning or updating files")
    scanner.Scan(db)
    //scanner.ScanNewFiles(db)
    // server.Play()
    // Init()
}
