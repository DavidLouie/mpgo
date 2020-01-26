package main

import (
    "io"
    "io/ioutil"
    "os"
    "path/filepath"

    "github.com/davidlouie/mpgo/server"
    "github.com/davidlouie/mpgo/server/subsonic"
    "github.com/gdamore/tcell"
    "github.com/rivo/tview"
)

var PAGE_MAP = map[rune]string{
    '1': "browsing",
    '2': "queue",
}

func Init() {
    app := tview.NewApplication()
    pages := tview.NewPages()
    var c chan string = make(chan string)
    initBrowsingPage(pages, c)
    initQueuePage(pages, c)
    go subsonic.Init()
    if err := app.SetRoot(pages, true).SetFocus(pages).Run(); err != nil {
        panic(err)
    }
}

// creates the browsing page and adds it to the Pages
func initBrowsingPage(pages *tview.Pages, c chan<- string) {
    rootDir := "."
    root := tview.NewTreeNode(rootDir).
            SetColor(tcell.ColorRed)
    tree := tview.NewTreeView().
            SetRoot(root).
            SetCurrentNode(root)

    getFiles(root, rootDir)
    setTreeCallback(tree, c)
    tree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        return swapPage(pages, event)
    })
    pages.AddPage(PAGE_MAP['1'], tree, false, true)
}

// on selected item, either expand the directory or add song to queue
func setTreeCallback(tree *tview.TreeView, c chan<- string) {
    tree.SetSelectedFunc(func(node *tview.TreeNode) {
        reference := node.GetReference()

        // selecting the root node does nothing
        if reference == nil {
            return
        }
        children := node.GetChildren()
        if len(children) == 0 {
            // load and show files in this directory
            path := reference.(string)
            file, err := os.Stat(path)
            switch {
            case err != nil:
                panic(err)
            case file.IsDir():
                getFiles(node, path)
            default:
                name := file.Name()
                server.Add(name)
                c <- name
            }
        } else {
            // collapse if visible, expand if collapsed
            node.SetExpanded(!node.IsExpanded())
        }
    })
}

// changes the page to the page specified in PAGE_MAP
func swapPage(pages *tview.Pages, event *tcell.EventKey) *tcell.EventKey {
    if val, ok := PAGE_MAP[event.Rune()]; ok {
        pages.SwitchToPage(val)
        return nil
    }
    return event
}

// creates the queue page and adds it to the Pages
func initQueuePage(pages *tview.Pages, c <-chan string) {
    textView := tview.NewTextView().
                SetDynamicColors(true)
    textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        return swapPage(pages, event)
    })
    pages.AddPage(PAGE_MAP['2'], textView, false, false)
    go queueUpdater(textView, c)
}

// loop routine that updates the queue page when a new song is queued
func queueUpdater(textView *tview.TextView, c <-chan string) {
    for {
        queued := <- c
        io.WriteString(textView, queued + "\n")
    }
}

// builds the file treeview, only showing directories and mp3s
func getFiles(target *tview.TreeNode, path string) {
    files, err := ioutil.ReadDir(path)
    if err != nil {
        panic(err)
    }

    for _, file := range files {
        name := file.Name()
        node := tview.NewTreeNode(name).
                SetReference(filepath.Join(path, name)).
                SetSelectable(true)
        if file.IsDir() {
            node.SetColor(tcell.ColorGreen)
        } else {
            ext := filepath.Ext(name)
            if ext != ".mp3" {
                continue
            }
        }

        target.AddChild(node)
    }
}
