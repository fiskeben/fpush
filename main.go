package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gregdel/pushover"
)

const stateFilePath = ".config/fpush/"
const stateFileName = "state"

var verbose bool

func main() {
	var dirname string

	flag.BoolVar(&verbose, "verbose", false, "log extra information")
	flag.StringVar(&dirname, "path", ".", "path to check for changes")
	flag.Parse()

	log.Printf("args %v %v", verbose, dirname)

	if !strings.HasSuffix(dirname, "/") {
		dirname = dirname + "/"
	}

	lastPush, err := readStateFile()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	now := time.Now()

	verboseLog("reading %s", dirname)
	dirs, err := ioutil.ReadDir(dirname)
	if err != nil {
		fmt.Printf("failed to read %s: %v", dirname, err)
		os.Exit(1)
	}

	var filename string

	for _, d := range dirs {
		if !strings.HasSuffix(d.Name(), ".jpg") {
			continue
		}

		age := now.Sub(d.ModTime())
		diff := d.ModTime().Sub(lastPush)
		verboseLog("checking %s (age=%v)", d.Name(), age)
		if age.Seconds() < 360 && diff > 3600 {
			filename = dirname + d.Name()
			lastPush = now
			break
		}
	}

	if filename != "" {
		verboseLog("found file to alert")
		if err = sendPushNotification(filename); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	verboseLog("writing state file")
	if err = writeStateFile(lastPush); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func verboseLog(msg string, args ...interface{}) {
	if verbose {
		fmt.Printf(msg+"\n", args...)
	}
}

func readStateFile() (time.Time, error) {
	home := getHome()
	b, err := ioutil.ReadFile(home + stateFilePath + stateFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return time.Now(), nil
		}
		return time.Now(), fmt.Errorf("unable to read state file: %v", err)
	}
	t, err := time.Parse(time.RFC3339, string(b))
	if err != nil {
		return time.Now(), fmt.Errorf("failed to parse state date: %v", err)
	}
	return t, nil
}

func writeStateFile(t time.Time) error {
	home := getHome()
	d := t.Format(time.RFC3339)
	err := ioutil.WriteFile(home+stateFilePath+stateFileName, []byte(d), os.FileMode(0600))
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to write state file: %v", err)
		}
		if err = os.MkdirAll(home+stateFilePath, os.FileMode(0770)); err != nil {
			return fmt.Errorf("failed to create config folder: %v", err)
		}
		f, err := os.Create(home + stateFilePath + stateFileName)
		if err != nil {
			return fmt.Errorf("failed to create state file: %v", err)
		}
		_, err = f.WriteString(d)
		if err != nil {
			return fmt.Errorf("failed to write new state file: %v", err)
		}
	}
	return nil
}

func sendPushNotification(filename string) error {
	p := pushover.New(os.Getenv("PUSHOVER_KEY"))
	r := pushover.NewRecipient(os.Getenv("PUSHOVER_RECIPIENT_KEY"))
	m := pushover.NewMessageWithTitle("Movement!", "Movement has been detected.")

	if err := addAttachment(m, filename); err != nil {
		log.Println(err.Error())
	}

	res, err := p.SendMessage(m, r)
	if err != nil {
		return fmt.Errorf("failed to send push notification: %v", err)
	}
	log.Println(res.String())
	return nil
}

func addAttachment(m *pushover.Message, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file '%s': %v", filename, err)
	}
	if err := m.AddAttachment(file); err != nil {
		return fmt.Errorf("failed to add attachment '%s', skipping: %v", filename, err)
	}
	return nil
}

func getHome() string {
	h := os.Getenv("HOME")
	if h != "" {
		if !strings.HasSuffix(h, "/") {
			h = h + "/"
		}
		return h
	}
	panic("unable to locate home dir!") // todo
}
