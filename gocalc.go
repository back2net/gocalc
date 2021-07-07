package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	usage()
	fmt.Fprint(os.Stdout, "$go-calc: ")
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	text = strings.ReplaceAll(text, "\n", "")
	text = strings.ReplaceAll(text, " ", "")

	if len(text) < 3 {
		fmt.Fprintln(os.Stderr, "Error. Enter at least one expression.")
		return
	}

	tmpArr := strings.Split(text, ",")
	tmpArrSize := len(tmpArr)
	done := make(chan int, tmpArrSize)

	for _, s := range tmpArr {
		parseExpression(s, done)
	}

	// goroutine sync
	func(chan int) {
		grtSpawned := tmpArrSize
		var grtCount int
		defer close(done)

		for grtCount < grtSpawned {
			grtCount += <-done
		}
	}(done)

}

func isMathOperator(r rune) bool {
	return r == '+' || r == '-' || r == '*' || r == '/'
}

func parseExpression(s string, done chan int) {

	for idx, char := range s {
		if unicode.IsDigit(char) || isMathOperator(char) {

			switch char {

			case '+':
				x, err := strconv.Atoi(s[:idx])
				if err != nil {
					break
				}
				y, err := strconv.Atoi(s[idx+1:])
				if err != nil {
					break
				}
				go add(x, y, done)
				return

			case '-':
				x, err := strconv.Atoi(s[:idx])
				if err != nil {
					break
				}
				y, err := strconv.Atoi(s[idx+1:])
				if err != nil {
					break
				}
				go sub(x, y, done)
				return

			case '*':
				x, err := strconv.Atoi(s[:idx])
				if err != nil {
					break
				}
				y, err := strconv.Atoi(s[idx+1:])
				if err != nil {
					break
				}
				go pow(x, y, done)
				return

			case '/':
				x, err := strconv.Atoi(s[:idx])
				if err != nil {
					break
				}
				y, err := strconv.Atoi(s[idx+1:])
				if err != nil {
					break
				}
				go div(x, y, done)
				return

			default:
				continue
			}
		}
	}
	// we are here if string contain non digit or none of math operators'
	fmt.Fprintf(os.Stderr, "%s: is not valid math expression.\n", s)
	done <- 1
}

func add(x, y int, done chan int) {
	defer func() { done <- 1 }()
	result := x + y
	fmt.Fprintf(os.Stdout, "%d + %d = %d\n", x, y, result)
}

func sub(x, y int, done chan int) {
	defer func() { done <- 1 }()
	result := x - y
	fmt.Fprintf(os.Stdout, "%d - %d = %d\n", x, y, result)
}

func pow(x, y int, done chan int) {
	defer func() { done <- 1 }()
	result := x * y
	fmt.Fprintf(os.Stdout, "%d * %d = %d\n", x, y, result)
}

func div(x, y int, done chan int) {
	defer func() { done <- 1 }()
	if y == 0 {
		fmt.Fprintf(os.Stderr, "%d / %d: divizion by zero!\n", x, y)
		return
	}
	result := float64(x / y)
	fmt.Fprintf(os.Stdout, "%d / %d = %f\n", x, y, result)
}

func usage() {
	text := `This is simple goroutine calculator.
	Usage:
	Enter comma separated math expressions.
	For example:
	2 + 5, 4 - 6, 8* 10, 14/2
	Quantity, order and spaces doesn't matters.`
	fmt.Fprintf(os.Stdout, "%s\n\n", text)
}
