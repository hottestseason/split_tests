package main

import (
	"io"
	"os"
	"path"
	"strconv"

	"github.com/antchfx/xmlquery"
	"github.com/bmatcuk/doublestar"
)

func loadJUnitXML(reader io.Reader) map[string]float64 {
	doc, err := xmlquery.Parse(reader)
	if err != nil {
		fatalMsg("failed to parse junit xml: %v\n", err)
	}

	testCases := make(map[string]float64)

	xmlquery.FindEach(doc, "//testsuite", func(_ int, node *xmlquery.Node) {
		timeStr := node.SelectAttr("time")
		fileTime, err := strconv.ParseFloat(timeStr, 64)
		if err != nil {
			fatalMsg("failed to parse time: %s %v\n", timeStr, err)
		}
		testCases[node.SelectAttr("name")] += fileTime
	})

	return testCases
}

func addFileTimesFromIOReader(fileTimes map[string]float64, reader io.Reader) {
	testCases := loadJUnitXML(reader)
	for file, fileTime := range testCases {
		filePath := path.Clean(file)
		fileTimes[filePath] += fileTime
	}
}

func getFileTimesFromJUnitXML(fileTimes map[string]float64) {
	if junitXMLPath != "" {
		filenames, err := doublestar.Glob(junitXMLPath)
		if err != nil {
			fatalMsg("failed to match jUnit filename pattern: %v", err)
		}
		for _, junitFilename := range filenames {
			file, err := os.Open(junitFilename)
			if err != nil {
				fatalMsg("failed to open junit xml: %v\n", err)
			}
			printMsg("using test times from JUnit report %s\n", junitFilename)
			addFileTimesFromIOReader(fileTimes, file)
			file.Close()
		}
	} else {
		printMsg("using test times from JUnit report at stdin\n")
		addFileTimesFromIOReader(fileTimes, os.Stdin)
	}
}
