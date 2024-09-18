package glox

import "fmt"

func Compile(source string) {
	scanner.initScanner(source)
	var line int = -1
	for {
		token := scanner.scanToken()
		if token.line != line {
			fmt.Printf("%4d ", token.line)
			line = token.line
		} else {
			fmt.Printf("   | ")
		}
		fmt.Printf("%2d '%s'\n", token.tokenType, token.value)

		if token.tokenType == TOKEN_EOF {
			break
		}
	}
}
