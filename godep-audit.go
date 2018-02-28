package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const (
	packageLine = "name ="
)

type Whitelist struct {
	Packages []Package `json:"whitelisted_packages"`
}

type Package struct {
	Name           string `json:"name"`
	UpstreamCommit string `json:"upstream_commit"`
}

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

func processDepOutput(b bytes.Buffer) map[string]string {
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

func getWhitelistedPackages(whitelistFile string) Whitelist {
	var whitelist Whitelist
	jsonFile, err := os.Open(whitelistFile)
	defer jsonFile.Close()
	if err != nil {
		json.Unmarshal([]byte("{}"), &whitelist)
		return whitelist
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		json.Unmarshal([]byte("{}"), &whitelist)
		return whitelist
	}

	err = json.Unmarshal(byteValue, &whitelist)
	if err != nil {
		json.Unmarshal([]byte("{}"), &whitelist)
	}
	return whitelist
}

func filterWhitelistedPkgs(pkgs map[string]string, whitelist Whitelist) {
	for _, whitelistPkg := range whitelist.Packages {
		if _, ok := pkgs[whitelistPkg.Name]; ok &&
			getRHSValue(pkgs[whitelistPkg.Name], "->") == whitelistPkg.UpstreamCommit {
			delete(pkgs, whitelistPkg.Name)
		}
	}
}

func main() {
	whitelistFile := flag.String("whitelist", "whitelist.json", "Path to whitelist file.")
	flag.Parse()

	pkgs := processDepOutput(runGoDepCommand())
	whitelistedPkgs := getWhitelistedPackages(*whitelistFile)
	filterWhitelistedPkgs(pkgs, whitelistedPkgs)

	xmlString := generateXMLString(pkgs)
	fmt.Println("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" + xmlString)
}
