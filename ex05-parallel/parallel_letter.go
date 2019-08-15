package letter

func Frequency(text string) map[rune]int {
	countOfSym := make(map[rune]int)

	for _, i := range text {
		countOfSym[i] += 1
	}

	return countOfSym
}

func ConcurrentFrequency(text []string) map[rune]int {
	countOfSym := make(map[rune]int)
	chanel := make(chan map[rune]int, len(text))

	for _, i := range text {
		go func(str string) {
			chanel <- Frequency(str)
		}(i)
	}

	for i := 0; i < len(text); i++ {
		for symbol, number := range <-chanel {
			countOfSym[symbol] += number
		}
	}

	return countOfSym
}
