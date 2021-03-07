package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/github/hub/v2/github"
)

const defaultEditor = "vim"

const helpTemplate = `
# ---------------------------------
# keep| k - keep release and tag
# del| d - delete both release and tag
# rel| r - delete release but keep the tag
`

var releasesLimit int

func init() {
	flag.IntVar(&releasesLimit, "limit", 50, "number of releases to show")
	flag.Parse()
}

func getReleases() ([]github.Release, error) {
	localRepo, err := github.LocalRepo()
	if err != nil {
		return nil, err
	}

	project, err := localRepo.MainProject()
	if err != nil {
		return nil, err
	}

	gh := github.NewClient(project.Host)

	approveAllF := func(*github.Release) bool { return true }
	return gh.FetchReleases(project, releasesLimit, approveAllF)
}

func createTemplate(releases []github.Release) string {
	output := ""
	for _, rel := range releases {
		draftString := ""
		if rel.Draft == true {
			draftString = "[DRAFT] "
		}

		preRelString := ""
		if rel.Prerelease == true {
			preRelString = "[PRERELEASE] "
		}
		output = fmt.Sprintf("%skeep - %s%s(%s) %s\n", output, draftString, preRelString, rel.TagName, rel.Name)
	}

	return fmt.Sprintf("%s%s", output, helpTemplate)
}

func writeTemplate(str string) (string, error) {

	file, err := ioutil.TempFile(os.TempDir(), "prerel_cleaner")
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(file.Name(), []byte(str), 0644)
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}

func openFileInEditor(fileName string) error {

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = defaultEditor
	}

	executable, err := exec.LookPath(editor)
	if err != nil {
		return err
	}

	cmd := exec.Command(executable, fileName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func parseTMPFile(fileName string) ([]string, []string, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	releasesToDelete := []string{}
	tagsToDelete := []string{}

	re := regexp.MustCompile(`(keep|k|del|d|rel|r)\s+\-\s+(\[PRERELEASE\]|)(\s+|)(\[DRAFT\]|)(\s+|)\((.*?)\).*$`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		matches := re.FindStringSubmatch(line)

		if len(matches) < 7 {
			return nil, nil, fmt.Errorf("could not parse line '%s'", line)
		}

		action := matches[1]
		if action == "keep" || action == "k" {
			continue
		}

		tagName := matches[6]
		if action == "rel" || action == "r" {
			releasesToDelete = append(releasesToDelete, tagName)
			continue
		}

		if action == "del" || action == "d" {
			releasesToDelete = append(releasesToDelete, tagName)
			tagsToDelete = append(tagsToDelete, tagName)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return releasesToDelete, tagsToDelete, nil
}

func askForConfirmation(s string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			return false, err
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true, nil
		}
		if response == "n" || response == "no" {
			return false, nil
		}
	}
}

func deleteRelease(tagName string) error {

	localRepo, err := github.LocalRepo()
	if err != nil {
		return err
	}

	project, err := localRepo.MainProject()
	if err != nil {
		return err
	}

	gh := github.NewClient(project.Host)

	release, err := gh.FetchRelease(project, tagName)
	if err != nil {
		return err
	}

	return gh.DeleteRelease(release)
}

func deleteTag(tagName string) error {
	executable, err := exec.LookPath("git")
	if err != nil {
		return fmt.Errorf("could not find git: %w", err)
	}

	// @todo - hardcoding "origin" here seems very wrong
	cmd := exec.Command(executable, "push", "--delete", "origin", tagName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func handleError(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}

func handleSucces(msg string) {
	fmt.Println(msg)
	os.Exit(0)
}

func main() {

	releases, err := getReleases()
	if err != nil {
		handleError(fmt.Errorf("could not fetch releases %w", err))
	}

	tpl := createTemplate(releases)

	fileName, err := writeTemplate(tpl)
	if err != nil {
		handleError(fmt.Errorf("cannot create template: %w", err))
	}

	err = openFileInEditor(fileName)
	if err != nil {
		handleError(fmt.Errorf("could not open system editor: %w", err))
	}

	relToDelete, tagsToDelete, err := parseTMPFile(fileName)
	if err != nil {
		handleError(fmt.Errorf("error parsing releases file: %w", err))
	}

	if len(relToDelete) == 0 && len(tagsToDelete) == 0 {
		handleSucces("Nothing to do")
	}

	if len(relToDelete) > 0 {
		fmt.Println("I am going to delete the following releases:")
		for _, tagName := range relToDelete {
			fmt.Printf("\t %s \n", tagName)
		}
	}

	if len(tagsToDelete) > 0 {
		fmt.Println("I am going to delete the following tags:")
		for _, tagName := range tagsToDelete {
			fmt.Printf("\t %s \n", tagName)
		}
	}

	res, err := askForConfirmation("Does this look good to you ?")
	if err != nil {
		handleError(fmt.Errorf("error reading user input: %w", err))
	}

	if res == false {
		handleSucces("Aborting ...")
	}
	for _, tagName := range relToDelete {
		fmt.Printf("Deleting release: %s \n", tagName)
		err = deleteRelease(tagName)
		if err != nil {
			handleError(err)
		}
	}

	for _, tagName := range tagsToDelete {
		fmt.Printf("Deleting tag: %s \n", tagName)
		err = deleteTag(tagName)
		if err != nil {
			handleError(err)
		}
	}

	// @todo - clean the tmp file
}
