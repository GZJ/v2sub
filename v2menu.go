package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	dir := flag.String("dir", ".", "Directory to search for JSON files")
	stopCmd := flag.String("stop-cmd", "pkill", "Command to stop the process")
	stopArgs := flag.String("stop-args", "v2ray", "Arguments for the stop command")
	runCmd := flag.String("run-cmd", "v2ray", "Command to run")
	runArgsTemplate := flag.String("run-args", "--config=%s", "Arguments template for the run command")
	flag.Parse()

	app := tview.NewApplication()

	files, err := getJSONFiles(*dir)
	if err != nil {
		log.Fatal(err)
	}

	list := tview.NewList().ShowSecondaryText(false)
	for _, file := range files {
		list.AddItem(file, "", 0, nil)
	}

	outputText := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetMaxLines(1000).
		ScrollToEnd()

	outputText.SetChangedFunc(func() {
		app.Draw()
	})

	flex := tview.NewFlex().
		AddItem(list, 0, 1, true).
		AddItem(outputText, 0, 2, false)

	var cmdMutex sync.Mutex

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case 'j':
				list.SetCurrentItem(list.GetCurrentItem() + 1)
			case 'k':
				list.SetCurrentItem(list.GetCurrentItem() - 1)
			case 'q':
				app.Stop()
			}
		} else if event.Key() == tcell.KeyEnter {
			selectedItem := list.GetCurrentItem()
			if selectedItem >= 0 && selectedItem < len(files) {
				fileName := files[selectedItem]
				cmdMutex.Lock()
				defer cmdMutex.Unlock()

				stopArgsSplit := strings.Split(*stopArgs, " ")
				cmd := exec.Command(*stopCmd, stopArgsSplit...)
				_, err := cmd.Output()
				if err != nil {
					fmt.Fprintf(outputText, "%s ", err)
				}

				runArgs := fmt.Sprintf(*runArgsTemplate, fileName)
				runArgsSplit := strings.Split(runArgs, " ")
				cmd = exec.Command(*runCmd, runArgsSplit...)
				cmd.Stdout = outputText
				cmd.Stderr = outputText
				err = cmd.Start()
				if err != nil {
					fmt.Fprintf(outputText, "%s ", err)
				}
				go func() {
					cmd.Wait()
					app.QueueUpdateDraw(func() {
						fmt.Fprintln(outputText, "------------------------------------------------------------")
					})
				}()
			}
		}
		return event
	})

	if err := app.SetRoot(flex, true).Run(); err != nil {
		log.Fatal(err)
	}
}

func getJSONFiles(dir string) ([]string, error) {
	var files []string
	fileInfos, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() && strings.HasSuffix(fileInfo.Name(), ".json") {
			files = append(files, fileInfo.Name())
		}
	}
	return files, nil
}
