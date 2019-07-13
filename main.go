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
	verboseLog("last push %v", lastPush)

	now := time.Now()

	verboseLog("reading %s", dirname)
	dirs, err := ioutil.ReadDir(dirname)
	if err != nil {
		fmt.Printf("failed to read %s: %v", dirname, err)
		os.Exit(1)
	}

	filenames := make([]string, 0)

	for _, d := range dirs {
		if !strings.HasSuffix(d.Name(), ".jpg") {
			continue
		}

		age := now.Sub(d.ModTime())
		diff := d.ModTime().Sub(lastPush)
		verboseLog("checking %s (age=%v, since last push=%v)", d.Name(), age, diff)
		if age.Seconds() < 360 && diff.Seconds() > 3600 {
			log.Printf("ok %s", d.Name())
			filename := dirname + d.Name()
			filenames = append(filenames, filename)
		}
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
