package main

import "log"

func main() {
	log.Println(addDigits(18))
}

func divmode(num, divNum int) (int, int) {
	mod := num % divNum
	div := num / divNum
	return div, mod
}

func addDigits(num int) int {
	var sum int
	for num > 9 {

		temp := num
		for temp != 0 {
			div, mod := divmode(temp, 10)
			temp = div
			sum += mod
		}

		num = sum

	}
	return num
}
