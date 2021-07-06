package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	usage()
	fmt.Fprint(os.Stdout, "$go-calc: ")
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
		return
	}
	text = strings.ReplaceAll(text, "\n", "")
	text = strings.ReplaceAll(text, " ", "")
	tmpArr := strings.Split(text, ",")

	done := make(chan int)

	for _, s := range tmpArr {
		op, x, y := parseExpression(s)
		switch {
		case op == '+':
			go add(x, y, done)
		case op == '-':
			go sub(x, y, done)
		case op == '*':
			go pow(x, y, done)
		case op == '/':
			go div(x, y, done)
		}

	}

	// goroutine sync
	func(chan int) {
		grtSpawned := len(tmpArr)
		var grtCount int
		defer close(done)
		for grtCount < grtSpawned {
			grtCount += <-done
		}
	}(done)

}

func parseExpression(s string) (op rune, x, y int) {

	for idx, char := range s {
		switch char {
		case '+':

			x, err := strconv.Atoi(s[:idx])
			if err != nil {
				log.Fatal(err)
			}

			y, err := strconv.Atoi(s[idx+1:])
			if err != nil {
				log.Fatal(err)
			}
			return char, x, y

		case '-':

			x, err := strconv.Atoi(s[:idx])
			if err != nil {
				log.Fatal(err)
			}

			y, err := strconv.Atoi(s[idx+1:])
			if err != nil {
				log.Fatal(err)
			}
			return char, x, y

		case '*':

			x, err := strconv.Atoi(s[:idx])
			if err != nil {
				log.Fatal(err)
			}

			y, err := strconv.Atoi(s[idx+1:])
			if err != nil {
				log.Fatal(err)
			}
			return char, x, y

		case '/':

			x, err := strconv.Atoi(s[:idx])
			if err != nil {
				log.Fatal(err)
			}

			y, err := strconv.Atoi(s[idx+1:])
			if err != nil {
				log.Fatal(err)
			}
			return char, x, y
		}
	}
	return
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
		log.Fatal("divizion by zero!")
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
	Order and spaces doesn't matters.`
	fmt.Fprintf(os.Stdout, "%s\n\n", text)
}
