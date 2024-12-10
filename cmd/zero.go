package cmd

import (
	"fmt"
	"strconv"
)

func Multiply(first string, second string, shouldRoundUp bool) (result string) {
	num1, err := strconv.ParseFloat(first, 64)
	if err != nil {
		fmt.Println("Error: First value is not a decimal")
		return
	}
	num2, err := strconv.ParseFloat(second, 64)
	if err != nil {
		fmt.Println("Error: Second value is not a decimal")
		return
	}
	if shouldRoundUp {
		return fmt.Sprintf("%.2f", num1*num2)
	}
	return fmt.Sprintf("%f", num1*num2)
}

func Divide(divide string, by string, shouldRoundUp bool) (result string, e error) {
	num1, err := strconv.ParseFloat(divide, 64)
	if err != nil {
		return "", fmt.Errorf("first value is not a number")
	}
	num2, err := strconv.ParseFloat(by, 64)
	if err != nil {
		return "", fmt.Errorf("second value is not a number")
	}
	if shouldRoundUp {
		return fmt.Sprintf("%.2f", num1/num2), nil
	}
	return fmt.Sprintf("%f", num1/num2), nil
}
