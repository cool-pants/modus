package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	// Give botton-left padding of 2 from the edge
	if v, err := g.SetView("hello", 0, 0, maxX-2, maxY-5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Hello world!")
	}
	if tbline, err := g.SetView("tabline", 0, maxY-4, maxX-2, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(tbline, "NORMAL")
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
