package packages

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// Sizes gets a list of packages sorted by size
func Sizes(filter string, size bool, human bool) ([]string, error) {
	var (
		lines []string
		cmd   *exec.Cmd
	)

	if size {
		cmd = exec.Command("dpkg-query", "-Wf", "${Installed-Size},${Package}\n")
	} else {
		cmd = exec.Command("dpkg-query", "-Wf", "${Package}\n")
	}

	//cmd := exec.Command("testlist", "-l", "a")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err = cmd.Start(); err != nil {
		return nil, err
	}

	for scanner := bufio.NewScanner(stdout); scanner.Scan(); {
		line := scanner.Text()

		if !strings.Contains(line, filter) {
			continue
		}

		if !size {
			lines = append(lines, line)
			continue
		}

		log.Println("scanning line", line)

		line = addLeadingZeroes(line)

		if human {
			line, err = addHuman(line)
			if err != nil {
				return nil, err
			}
		}

		lines = appendSorted(lines, line)
	}

	if err = cmd.Wait(); err != nil {
		return nil, err
	}

	log.Println("finished with packages.Sizes")

	return lines, nil
}

func appendSorted(lines []string, line string) []string {
	index := sort.Search(len(lines), func(n int) bool { return lines[n] >= line })

	//log.Printf("line '%v' will fit into position %v\n", line, index)
	lines = append(lines, "")
	copy(lines[index+1:], lines[index:])
	lines[index] = line

	return lines
}

func addLeadingZeroes(line string) string {
	number := strings.Split(line, ",")[0]
	l := len(number)
	//adds leading zeroes to ease sorting
	for i := 0; i < 10-l; i++ {
		line = "0" + line
	}

	return line
}

func addHuman(line string) (string, error) {
	var (
		c      = 0
		number = strings.Split(line, ",")[0]
		units  = "G"
	)

	f, err := strconv.ParseFloat(number, 32)
	if err != nil {
		e := fmt.Sprintf("can't convert %s to double", number)
		return "", errors.New(e)
	}

	log.Println("converted", line, "into", f)

	for f > 1024 && c < 3 {
		f = f / 1024
		c++
	}

	if c == 0 {
		units = "K"
	} else if c == 1 {
		units = "M"
	}

	s := fmt.Sprintf("%v,%.0f%s", line, f, units)
	return s, nil
}
