package packages

import (
	"bufio"
	"github.com/cristianrz/aptlist/size"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
)

// Sizes gets a list of packages sorted by size
func Sizes(human bool) ([]string, error) {
	var lines []string

	cmd := exec.Command("dpkg-query", "-Wf", "${Installed-Size},${Package}\n")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err = cmd.Start(); err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()

		log.Println("scanning line", line)

		number := strings.Split(line, ",")[0]
		l := len(number)

		// adds leading zeroes to ease sorting
		for i := 0; i < 10-l; i++ {
			line = "0" + line
		}

		if human {
			h, err := size.Human(line)
			if err != nil {
				return nil, err
			}

			line = line + "," + h
		}

		lines = append(lines, line)

	}

	sort.Strings(lines)

	if err = cmd.Wait(); err != nil {
		return nil, err
	}

	return lines, nil
}

// List gets a list of installed packages
func List() error {
	cmd := exec.Command("dpkg-query", "-Wf", "${Package}\n")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
