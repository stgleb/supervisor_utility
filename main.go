package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/prometheus/common/log"
	"os"
	"text/template"
	"strconv"
)

var (
	programName  string
	configPath   string
	workdir      string
	templateName string
	outputFile   string
	cpuCount     int
	port int
)

type Payload struct {
	ProgramName string
	ConfigPath  string
	Workdir     string
	Mask        string
	Port		int
}

func main() {
	flag.StringVar(&programName, "programName", "redirector", "name of process in supervisor")
	flag.StringVar(&configPath, "configPath", "config.toml", "path to config file")
	flag.StringVar(&workdir, "workdir", "", "working directory")
	flag.StringVar(&templateName, "templateName", "template.conf", "template file full name")
	flag.StringVar(&outputFile, "outputFile", "redirector.conf", "output file name")
	flag.IntVar(&cpuCount, "cpuCount", 4, "cpu count")
	flag.IntVar(&port, "port", 9001, "port")
	flag.Parse()

	buffer := &bytes.Buffer{}

	for i := 0; i < cpuCount; i += 1 {
		payload := Payload{
			ProgramName: fmt.Sprintf("%s_%d", programName, i),
			ConfigPath:  configPath,
			Workdir:     workdir,
			Mask:        strconv.Itoa(1 << uint(i)),
			Port: port + i,
		}

		if t, err := template.ParseFiles(templateName); err == nil {
			err = t.ExecuteTemplate(buffer, templateName, payload)

			if err != nil {
				log.Fatal(err)
			}

			buffer.Write([]byte{13, 10, 13, 10})
		} else {
			log.Fatal(err)
		}
	}

	if f, err := os.Create(outputFile); err == nil {
		f.Write(buffer.Bytes())
	} else {
		log.Fatal(err)
	}
}
