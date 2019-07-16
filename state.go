package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const stateFilePath = ".config/fpush/"
const stateFileName = "state"

var nowfunc = time.Now

func getLastPushTime() (time.Time, error) {
	home := getHome()
	return readStateFile(home + stateFilePath + stateFileName)
}

func readStateFile(path string) (time.Time, error) {
	def := nowfunc().Add(-60 * time.Minute)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return def, nil
		}
		return def, fmt.Errorf("unable to read state file: %v", err)
	}
	s := strings.TrimSpace(string(b))
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return def, fmt.Errorf("failed to parse state date: %v", err)
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
