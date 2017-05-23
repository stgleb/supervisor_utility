package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/prometheus/common/log"
	"os"
	"text/template"
)

var (
	programName  string
	configPath   string
	workdir      string
	templateName string
	outputFile   string
	cpuCount     int
)

type Payload struct {
	ProgramName string
	ConfigPath  string
	Workdir     string
	Mask        string
}

func GetMask(procCount, procNum int) string {
	mask := make([]byte, procCount)

	for i := 0; i < procCount; i++ {
		if i == procNum {
			mask[i] = byte('1')
		} else {
			mask[i] = byte('0')
		}
	}

	return string(mask)
}

func main() {
	flag.StringVar(&programName, "programName", "redirector", "name of process in supervisor")
	flag.StringVar(&configPath, "configPath", "", "path to config file")
	flag.StringVar(&workdir, "workdir", "", "working directory")
	flag.StringVar(&templateName, "templateName", "template.conf", "template file full name")
	flag.StringVar(&outputFile, "outputFile", "redirector.conf", "output file name")
	flag.IntVar(&cpuCount, "cpuCount", 4, "cpu count")
	flag.Parse()

	buffer := &bytes.Buffer{}

	for i := 0; i < cpuCount; i += 1 {
		payload := Payload{
			ProgramName: fmt.Sprintf("%s_%d", programName, i),
			ConfigPath:  configPath,
			Workdir:     workdir,
			Mask:        GetMask(cpuCount, i),
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
