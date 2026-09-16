package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var scanner = bufio.NewScanner(os.Stdin)

func String(label string) (string, error) {
	fmt.Printf("%s: ", label)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", err
		}
		return "", fmt.Errorf("no input")
	}
	return strings.TrimSpace(scanner.Text()), nil
}

func StringRequired(label string) (string, error) {
	for {
		val, err := String(label)
		if err != nil {
			return "", err
		}
		if val != "" {
			return val, nil
		}
		fmt.Println("  (required)")
	}
}

func StringWithDefault(label, defaultVal string) (string, error) {
	fmt.Printf("%s [%s]: ", label, defaultVal)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", err
		}
		return defaultVal, nil
	}
	val := strings.TrimSpace(scanner.Text())
	if val == "" {
		return defaultVal, nil
	}
	return val, nil
}

func Confirm(label string) bool {
	fmt.Printf("%s [y/N]: ", label)
	if !scanner.Scan() {
		return false
	}
	return strings.HasPrefix(strings.ToLower(strings.TrimSpace(scanner.Text())), "y")
}
