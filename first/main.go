package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	weight := transformStringToFloat(readLine("Введите ваш вес"));
	height := transformStringToFloat(readLine("Введите ваш рост"));

	imt := weight / math.Pow(height, 2)
	

	switch {
		case imt < 16:
			fmt.Println("Недостаточный вес")
		case imt > 16 && imt < 18.5:
			fmt.Println("Дефицит массы тела")
		case imt > 18.5 && imt < 25:
			fmt.Println("Нормальный вес")
		case imt > 25 && imt < 30:
			fmt.Println("Избыточная масса")
		case imt > 30 && imt < 35:
			fmt.Println("1-я степень ожирения")
		case imt > 35 && imt < 40:
			fmt.Println("2-я степень ожирения")
		case imt > 40:
			fmt.Println("3-я степень ожирения")
	}
}

func transformStringToFloat (text string) float64{
	number, err := strconv.ParseFloat(strings.TrimSpace(text), 64);

	if err != nil {
		log.Fatal(err)
	}

	return number;
}


func readLine(inputText string) string{
	reader := bufio.NewReader(os.Stdin)
	
	fmt.Print(inputText);
	fmt.Print("-> ")

	text, err := reader.ReadString('\n')

	if err != nil {
        log.Fatal(err)
    }

	return text;
}