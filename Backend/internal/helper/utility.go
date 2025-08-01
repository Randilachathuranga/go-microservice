package helper

import (
	"crypto/rand"
	"strconv"
)

func Randomnumbers(length int) (int, error) {

	const numbers = "1234567890"
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return 0, err
	}

	numLenght := len(buffer)

	for i := 0; i < length; i++ {
		buffer[i] = numbers[int(buffer[i])%numLenght]
	}

	return strconv.Atoi(string(buffer))
}
