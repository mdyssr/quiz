package helpers

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func OpenFile(fileName string) (*os.File, error) {
	f, err := os.Open(fileName)

	if err != nil {
		err := errors.New(fmt.Sprintf("Couldn't open %s", fileName))
		return nil, err
	}

	return f, nil
}

func CleanText(text string) string {
	return strings.ToLower(strings.Trim(text, "\n\t"))
}

func Shuffle(elements [][]string) {
	rand.Seed(time.Now().UnixNano())

	rand.Shuffle(len(elements), func(i, j int) {
		elements[i], elements[j] = elements[j], elements[i]
	})
}
