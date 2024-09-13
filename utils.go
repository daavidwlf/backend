package main

import (
	"bufio"
	"errors"
	"io"
)

func parseLogs(logs io.Reader) ([]string, error) {
	var logsArray []string

	scanner := bufio.NewScanner(logs)

	for scanner.Scan() {
		logsArray = append(logsArray, scanner.Text())
	}

	err := scanner.Err()

	if err != nil {
		return nil, errors.New("uable to scan logs of io Stream: " + err.Error())
	}

	return logsArray, nil
}
