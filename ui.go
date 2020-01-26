package main

import (
    "io/ioutil"
    "os"
    "path/filepath"

    "github.com/davidlouie/mpgo/server"
    "github.com/gdamore/tcell"
    "github.com/rivo/tview"
)

func Init() {
    rootDir := "."
    root := tview.NewTreeNode(rootDir).
            SetColor(tcell.ColorRed)
    tree := tview.NewTreeView().
            SetRoot(root).
            SetCurrentNode(root)

    getFiles(root, rootDir)
    tree.SetSelectedFunc(openDir)
    if err := tview.NewApplication().SetRoot(tree, true).Run(); err != nil {
        panic(err)
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

// expanding directory or file, adding files if required
func openDir(node *tview.TreeNode) {
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
                server.Add(file.Name())
        }
    } else {
        // collapse if visible, expand if collapsed
        node.SetExpanded(!node.IsExpanded())
    }
}
