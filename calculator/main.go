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
	value := strings.TrimSpace(readLine())

	fmt.Printf("Длина строки %v\n", len(value))
	fmt.Printf("Итоговый ответ: %v\n",calculateExpression(value))
}

func readLine() string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Введите выражение")
	fmt.Print("-> ")

	text, err := reader.ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}

	return text;
}

func calculateExpression(expression string) float64 {

	currentNumber := ""
	var result float64
	transformedValue := []string{}

	const operations = "+-*/"
	chars := strings.Split(expression, "")

	for i := 0; i < len(expression); i++ {
		switch {
		case i+1 == len(expression):
			currentNumber += chars[i]
			transformedValue = append(transformedValue, currentNumber)
			currentNumber = ""
			break
		case !strings.Contains(operations, chars[i]):
			currentNumber += chars[i]
		default:
			transformedValue = append(transformedValue, currentNumber, chars[i])
			currentNumber = ""
		}
	}
	fmt.Printf("Result: %v\n", transformedValue)

	for i := 0; i < len(transformedValue)-1; i++ {
		transformedStrToInt, err := strconv.ParseFloat(transformedValue[i],  64)
		if(err != nil){
			nextValue, _ := strconv.ParseFloat(transformedValue[i+1], 64)
			checkOperations(transformedValue[i], &result, nextValue)
		}
		if(i==0){
			result += transformedStrToInt
			continue
		}
	}

	return result
}


func checkOperations(operation string, result *float64, value float64) float64 {
	switch operation {
	case "+":
		*result += value
	case "-":
		*result -= value
	case "*":
		*result *= value
	case "/":
		*result /= value
	}
	return *result
}
