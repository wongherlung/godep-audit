package main

import (
	"encoding/xml"
	"log"
)

type XMLReport struct {
	XMLName    xml.Name  `xml:"testsuites"`
	Testsuites Testsuite `xml:"testsuite"`
}

type Testsuite struct {
	XMLName   xml.Name   `xml:"testsuite"`
	Name      string     `xml:"name,attr"`
	Tests     int        `xml:"tests,attr"`
	Testcases []Testcase `xml:"testcase"`
}

type Testcase struct {
	XMLName xml.Name `xml:"testcase"`
	Name    string   `xml:"name,attr"`
	Failure Failure  `xml:"failure"`
}

type Failure struct {
	XMLName xml.Name `xml:"failure"`
	Message string   `xml:"message,attr"`
	Text    string   `xml:",innerxml"`
}

func generateXMLString(pkgs map[string]string) string {
	var xmlReport = XMLReport{
		Testsuites: Testsuite{
			Name:      "godep-audit",
			Tests:     len(pkgs),
			Testcases: []Testcase{},
		},
	}
	for name, version := range pkgs {
		var testcase = Testcase{
			Name: name + " has new updates.",
			Failure: Failure{
				Message: "This package has updates.",
				Text:    version,
			},
		}
		xmlReport.Testsuites.Testcases = append(xmlReport.Testsuites.Testcases, testcase)
	}

	xmlString, err := xml.MarshalIndent(xmlReport, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	return string(xmlString)
}
