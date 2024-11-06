package tests

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

func compileFeatures(pathToDefinitions, pathToPatterns, pathToFeatures string) error {
	m, err := parsePatternFolder(pathToPatterns)
	if err != nil {
		return err
	}
	logrus.Debug(m)
	return replaceDefs(pathToDefinitions, pathToFeatures, m)
}

func replaceDefs(pathToDefinitions, pathToFeatures string, patterns map[string][]string) error {
	files, err := listFolder(pathToDefinitions)
	if err != nil {
		return err
	}

	err = os.RemoveAll(pathToFeatures)
	if err != nil {
		return err
	}

	err = os.CopyFS(pathToFeatures, os.DirFS(pathToDefinitions))
	if err != nil {
		return err
	}

	for _, file := range files {
		filePath := filepath.Join(pathToFeatures, file)
		f, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer f.Close()

		var lines []string
		isReplaced := false

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			indentation := strings.Repeat(" ", len(line)-len(strings.TrimLeft(line, " ")))
			for pattern, replacementLines := range patterns {
				if strings.Contains(line, pattern) {
					for _, replaceLine := range replacementLines {
						lines = append(lines, indentation+replaceLine)
					}
					isReplaced = true
					break
				}
			}

			if !isReplaced {
				lines = append(lines, line)
			}
			isReplaced = false
		}

		if scanner.Err() != nil {
			return scanner.Err()
		}

		f, err = os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		defer f.Close()

		writer := bufio.NewWriter(f)
		for _, line := range lines {
			fmt.Fprintln(writer, line)
		}

		writer.Flush()
	}

	return nil
}

func parsePatternFolder(pathToPatterns string) (map[string][]string, error) {
	files, err := listFolder(pathToPatterns)
	if err != nil {
		return nil, err
	}

	result := map[string][]string{}
	for _, filename := range files {
		f, err := os.Open(filepath.Join(pathToPatterns, filename))
		if err != nil {
			return nil, err
		}

		m, err := parseDefinitions(f)
		if err != nil {
			return nil, err
		}

		for k, v := range m {
			result[k] = append(result[k], v...)
		}
	}

	return result, nil
}

func listFolder(dir_name string) ([]string, error) {
	dir, err := os.ReadDir(dir_name)
	if err != nil {
		return nil, err
	}

	var filenames []string
	for _, file := range dir {
		filenames = append(filenames, file.Name())
	}

	return filenames, nil
}

func parseDefinitions(reader io.Reader) (map[string][]string, error) {
	templates := map[string][]string{}

	var currentKey string
	var currentLines []string
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		if strings.HasSuffix(line, ":") {
			if currentKey != "" {
				templates[currentKey] = currentLines
			}

			currentKey = strings.TrimSuffix(line, ":")
			currentLines = []string{}
		} else if currentKey != "" {
			currentLines = append(currentLines, line)
		}
	}

	if currentKey != "" {
		templates[currentKey] = currentLines
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err
	}

	return templates, nil
}
