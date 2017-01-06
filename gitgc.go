package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	configFileName = ".gitgc"
)

var (
	userHome = findUserHome()
)

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	user, repo := parse(args[0])

	if user == "" || repo == "" {
		flag.Usage()
		os.Exit(1)
	}

	base := findBase()
	if base == "" {
		configFilePath := findConfigFilePath()
		log.Printf("Make sure the config file at '%s' has a valid base path.\n", configFilePath)
		os.Exit(1)
	}

	userDir := filepath.Join(base, user)
	repoDir := filepath.Join(userDir, repo)
	repoURL := "https://github.com/" + user + "/" + repo
	fmt.Printf("Cloning repo 'github.com/%s/%s'...", user, repo)

	err := os.MkdirAll(userDir, os.ModeDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directory at '%s'\n", userDir)
		fmt.Println(err)
		return
	}

	gitCmd := exec.Command("git", "clone", "--progress", repoURL, repoDir)
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr
	err = gitCmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func findBase() string {
	configFilePath := findConfigFilePath()
	config := ensureConfigFile(configFilePath)
	config = strings.TrimSpace(config)

	return config
}

func findConfigFilePath() string {
	return filepath.Join(userHome, configFileName)
}

// ensureConfigFile ensures that the config file exists and returns its content.
func ensureConfigFile(p string) string {
	_, err := os.Stat(p)
	if os.IsNotExist(err) {
		createConfigFile(p)
	}

	byteContent, err := ioutil.ReadFile(p)
	content := string(byteContent)
	return content
}

func createConfigFile(p string) {
	f, err := os.Create(p)
	defer f.Close()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Creating config file at '%s'.\n", p)
	_, err = f.WriteString(getDefaultBase())
	if err != nil {
		log.Fatal(err)
	}
}

func getDefaultBase() string {
	return filepath.Join(userHome, "git")
}

func findUserHome() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

func parse(arg string) (user string, repo string) {
	regex, _ := regexp.Compile("(.+)/(.+)")
	subs := regex.FindStringSubmatch(arg)

	if subs == nil || len(subs) < 3 { // subs[0] contains the whole capture
		return
	}

	user = subs[1]
	repo = subs[2]
	return
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: gitgc [user]/[repo]")
	flag.PrintDefaults()
}
