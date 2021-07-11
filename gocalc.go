package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Operator rune
type Operators []Operator

func (op *Operator) isValid() (valid bool) {
	for _, check := range ValidOperators {
		if valid = check == *op; valid {
			return
		}
	}
	return
}

var ValidOperators = Operators{'+', '-', '*', '/'}

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

func parseExpression(s string, done chan int) {
	for idx, char := range s {
		if unicode.IsDigit(char) || char == '.' {
			continue
		}
		operator := Operator(char)
		if operator.isValid() {
			strX := s[:idx]
			strY := s[idx+1:]
			var err error
			var x, y interface{}

			x, err = strconv.Atoi(strX)
			if err != nil {
				x, err = strconv.ParseFloat(strX, 64)
				if err != nil {
					break
				}
			}
			y, err = (strconv.Atoi(strY))
			if err != nil {
				y, err = strconv.ParseFloat(strY, 64)
				if err != nil {
					break
				}
			}

			outX := makeFloat(x)
			outY := makeFloat(y)

			switch char {
			case '+':
				go add(outX, outY, done)
			case '-':
				go sub(outX, outY, done)
			case '*':
				go pow(outX, outY, done)
			case '/':
				go div(outX, outY, done)
			}
			return
		}
	}
	// we are here if string contain non digit or none of math operators'
	fmt.Fprintf(os.Stderr, "%s: is not valid math expression.\n", s)
	done <- 1
}

func makeFloat(in interface{}) float64 {
	switch in := in.(type) {
	case int:
		in = int(in)
		return float64(in)
	case float64:
		return in
	}
	return 0
}

func add(x, y float64, done chan int) {
	defer func() { done <- 1 }()
	result := x + y
	fmt.Fprintf(os.Stdout, "%v + %v = %v\n", x, y, result)
}

func sub(x, y float64, done chan int) {
	defer func() { done <- 1 }()
	result := x - y
	fmt.Fprintf(os.Stdout, "%v - %v = %v\n", x, y, result)
}

func pow(x, y float64, done chan int) {
	defer func() { done <- 1 }()
	result := x * y
	fmt.Fprintf(os.Stdout, "%v * %v = %v\n", x, y, result)
}

func div(x, y float64, done chan int) {
	defer func() { done <- 1 }()
	if y == 0 {
		fmt.Fprintf(os.Stderr, "%v / %v: divizion by zero!\n", x, y)
		return
	}
	result := float64(x / y)
	fmt.Fprintf(os.Stdout, "%v / %v = %v\n", x, y, result)
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
