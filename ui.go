package main

import (
    "github.com/rivo/tview"
)

func Init() {
    box := tview.NewBox().SetBorder(true).SetTitle("Hello world!")
    if err := tview.NewApplication().SetRoot(box, true).Run(); err != nil {
        panic(err)
    }
}
