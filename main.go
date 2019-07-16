package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"
	"time"
)

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

	lastPush, err := getLastPushTime()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	verboseLog("last push %v", lastPush)

	now := time.Now()

	verboseLog("reading %s", dirname)

	filenames, err := findFiles(dirname, now, lastPush)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(filenames) > 0 {
		verboseLog("found files to alert")
		filenames = limitFiles(filenames)
		filename, err := concatenateFiles(filenames)
		if err != nil {
			log.Println(err)
			filename = filenames[0]
		}
		if err = sendPushNotification(filename); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		log.Println(filename)
		lastPush = now
	}

	verboseLog("writing state file")
	if err = writeStateFile(lastPush); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func findFiles(dirname string, now time.Time, lastPush time.Time) ([]string, error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir '%s': %v", dirname, err)
	}

	filenames := make([]string, 0)

	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".jpg") {
			continue
		}

		if checkFile(f, now, lastPush) {
			filename := dirname + f.Name()
			filenames = append(filenames, filename)
		}
	}

	return filenames, nil
}

type fileInfo interface {
	ModTime() time.Time
}

func checkFile(f fileInfo, now, lastPush time.Time) bool {
	age := now.Sub(f.ModTime())
	diff := now.Sub(lastPush)
	return age.Seconds() <= (5*time.Minute).Seconds() && diff.Seconds() >= (60*time.Minute).Seconds()
}

func limitFiles(files []string) []string {
	if len(files) <= 4 {
		return files
	}

	res := make([]string, 4)
	res[0] = files[0]
	res[3] = files[len(files)-1]
	files = files[1 : len(files)-1]
	halves := math.Floor(float64(len(files) / 2))
	first := int(math.Floor(float64(halves / 2)))
	second := int(halves + math.Floor(float64(halves/2)))
	res[1] = files[first]
	res[2] = files[second]

	return res
}

func verboseLog(msg string, args ...interface{}) {
	if verbose {
		fmt.Printf(msg+"\n", args...)
	}
}
