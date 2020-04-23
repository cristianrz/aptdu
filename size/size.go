package size

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// Human gets a line from dpkg and extract the size in human readable form
func Human(line string) (string, error) {
	var (
		c      = 0
		number = strings.Split(line,",")[0]
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

	switch c {
	case 0:
		units = "K"
	case 1:
		units = "M"
	}

	s := fmt.Sprintf("%.0f%s", f, units)
	return s, nil
}
