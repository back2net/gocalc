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

	doneAdd := make(chan struct{})
	doneSub := make(chan struct{})
	donePow := make(chan struct{})
	doneDiv := make(chan struct{})

	for _, s := range tmpArr {
		op, x, y := parseExpression(s)
		switch {
		case op == '+':
			go add(x, y, doneAdd)
		case op == '-':
			go sub(x, y, doneSub)
		case op == '*':
			go pow(x, y, donePow)
		case op == '/':
			go div(x, y, doneDiv)
		}

	}
	// goroutine sync
	<-doneAdd
	<-doneSub
	<-donePow
	<-doneDiv
}

func parseExpression(s string) (op rune, x, y int) {

	for idx, char := range s {
		switch {
		case char == '+':

			x, err := strconv.Atoi(s[:idx])
			if err != nil {
				log.Fatal(err)
			}

			y, err := strconv.Atoi(s[idx+1:])
			if err != nil {
				log.Fatal(err)
			}
			return char, x, y

		case char == '-':

			x, err := strconv.Atoi(s[:idx])
			if err != nil {
				log.Fatal(err)
			}

			y, err := strconv.Atoi(s[idx+1:])
			if err != nil {
				log.Fatal(err)
			}
			return char, x, y

		case char == '*':

			x, err := strconv.Atoi(s[:idx])
			if err != nil {
				log.Fatal(err)
			}

			y, err := strconv.Atoi(s[idx+1:])
			if err != nil {
				log.Fatal(err)
			}
			return char, x, y

		case char == '/':

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

func add(x, y int, doneAdd chan struct{}) {
	defer close(doneAdd)
	result := x + y
	fmt.Fprintf(os.Stdout, "%d + %d = %d\n", x, y, result)

}

func sub(x, y int, doneSub chan struct{}) {
	defer close(doneSub)
	result := x - y
	fmt.Fprintf(os.Stdout, "%d - %d = %d\n", x, y, result)
}

func pow(x, y int, donePow chan struct{}) {
	defer close(donePow)

	result := x * y
	fmt.Fprintf(os.Stdout, "%d * %d = %d\n", x, y, result)

}

func div(x, y int, doneDiv chan struct{}) {
	defer close(doneDiv)
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
