package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
)

func main() {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}
	if err := s.Init(); err != nil {
		log.Fatal(err)
	}

	s.SetStyle(tcell.StyleDefault)
	s.Clear()

	// Create log file
	logFile, _ := os.Create("keylog.txt")
	defer logFile.Close()

	quit := make(chan struct{})

	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				logMsg := fmt.Sprintf("Key: %v, Rune: %q (%d), Modifiers: %v\n",
					ev.Key(), ev.Rune(), ev.Rune(), ev.Modifiers())
				logFile.WriteString(logMsg)
				logFile.Sync()

				s.Clear()
				s.ShowCursor(0, 0)
				printStr(s, 0, 0, "Press keys to see their codes (Ctrl+C to quit)")
				printStr(s, 0, 2, fmt.Sprintf("Last Key: %v", ev.Key()))
				printStr(s, 0, 3, fmt.Sprintf("Rune: %q (%d)", ev.Rune(), ev.Rune()))
				printStr(s, 0, 4, fmt.Sprintf("Modifiers: %v", ev.Modifiers()))
				printStr(s, 0, 6, "Check keylog.txt for full log")
				s.Show()

				if ev.Key() == tcell.KeyCtrlC {
					close(quit)
					return
				}
			}
		}
	}()

	s.Clear()
	printStr(s, 0, 0, "Press keys to see their codes (Ctrl+C to quit)")
	s.Show()

	<-quit
	s.Fini()
	fmt.Println("Check keylog.txt for the full key log")
}

func printStr(s tcell.Screen, x, y int, str string) {
	for i, c := range str {
		s.SetContent(x+i, y, c, nil, tcell.StyleDefault)
	}
}
