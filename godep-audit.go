package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

const (
	packageLine = "name ="
)

func runGoDepCommand() bytes.Buffer {
	var out bytes.Buffer
	cmd := exec.Command("dep", "ensure", "-update", "-dry-run")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return out
}

func isAPackageLine(s string) bool {
	matched, err := regexp.MatchString(packageLine, s)
	if err != nil {
		log.Fatal(err)
	}
	return matched
}

func getRHSValue(s string, delim string) string {
	return strings.Trim(strings.Split(s, delim)[1], " \"")
}

func processOutput(b bytes.Buffer) map[string]string {
	pkgs := make(map[string]string)
	lines := strings.Split(b.String(), "\n")

	for i, line := range lines {
		if isAPackageLine(line) {
			name := getRHSValue(line, "=")
			version := getRHSValue(lines[i+1], "=")
			pkgs[name] = version
		}
	}

	return pkgs
}

func outputToJUnitXML(pkgs map[string]string) {
	for name, version := range pkgs {
		fmt.Println(name + " " + version)
	}
}

func main() {
	pkgs := processOutput(runGoDepCommand())
	outputToJUnitXML(pkgs)
}
