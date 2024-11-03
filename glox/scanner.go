package glox

type Scanner struct {
	source  string
	start   int
	current int
	length  int
	line    int
}

var scanner Scanner

func (scanner *Scanner) initScanner(source string) {
	scanner.source = source
	scanner.length = len(source)
	scanner.start = 0
	scanner.current = 0
	scanner.line = 1
}

func (scanner *Scanner) isAtEnd() bool {
	return scanner.current >= scanner.length
}

func (scanner *Scanner) makeToken(tokenType TokenType) Token {
	token := Token{}

	token.tokenType = tokenType
	token.start = scanner.start
	token.length = scanner.current - scanner.start
	token.line = scanner.line
	token.value = scanner.source[scanner.start:scanner.current]

	return token
}

func (scanner *Scanner) errorToken(err string) Token {
	token := Token{}

	token.tokenType = TOKEN_ERROR
	token.value = err
	token.start = 0
	token.line = scanner.line

	return token
}

func (scanner *Scanner) skipWhitespace() {
	for {
		if scanner.isAtEnd() {
			return
		}

		c := scanner.source[scanner.current]
		switch c {
		case ' ', '\r', '\t':
			scanner.advance()
		case '\n':
			scanner.line += 1
			scanner.advance()
		case '/':
			if scanner.peekNext() == '/' {
				// A comment goes until the end of the line.
				for scanner.peek() != '\n' && !scanner.isAtEnd() {
					scanner.advance()
				}
			} else {
				return
			}
		default:
			return
		}
	}
}

func (scanner *Scanner) string() Token {
	for scanner.peek() != '"' && !scanner.isAtEnd() {
		if scanner.peek() == '\n' {
			scanner.line += 1
		}
		scanner.advance()
	}

	if scanner.isAtEnd() {
		return scanner.errorToken("Unterminated string.")
	}

	// The closing ".
	scanner.advance()
	return scanner.makeToken(TOKEN_STRING)
}

func (scanner *Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (scanner *Scanner) number() Token {
	for scanner.isDigit(scanner.peek()) {
		scanner.advance()
	}

	if scanner.peek() == '.' && scanner.isDigit(scanner.peekNext()) {
		scanner.advance()
		for scanner.isDigit(scanner.peek()) {
			scanner.advance()
		}
	}

	return scanner.makeToken(TOKEN_NUMBER)
}

func (scanner *Scanner) isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func (scanner *Scanner) identifier() Token {
	for scanner.isAlpha(scanner.peek()) || scanner.isDigit(scanner.peek()) {
		scanner.advance()
	}

	return scanner.makeToken(scanner.identifierType())
}

func (scanner *Scanner) checkKeyword(start, end int, last string, tokenType TokenType) TokenType {
	got := scanner.source[scanner.start+start : scanner.start+end]
	if got == last {
		return tokenType
	}

	return TOKEN_IDENTIFIER
}

func (scanner *Scanner) identifierType() TokenType {
	switch scanner.source[scanner.start] {
	case 'a':
		return scanner.checkKeyword(1, 3, "nd", TOKEN_AND)
	case 'c':
		return scanner.checkKeyword(1, 5, "lass", TOKEN_CLASS)
	case 'e':
		return scanner.checkKeyword(1, 4, "lse", TOKEN_ELSE)
	case 'i':
		return scanner.checkKeyword(1, 2, "f", TOKEN_IF)
	case 'n':
		return scanner.checkKeyword(1, 3, "il", TOKEN_NIL)
	case 'o':
		return scanner.checkKeyword(1, 2, "r", TOKEN_OR)
	case 'p':
		return scanner.checkKeyword(1, 5, "rint", TOKEN_PRINT)
	case 'r':
		return scanner.checkKeyword(1, 6, "eturn", TOKEN_RETURN)
	case 's':
		return scanner.checkKeyword(1, 5, "uper", TOKEN_SUPER)
	case 'v':
		return scanner.checkKeyword(1, 3, "ar", TOKEN_VAR)
	case 'w':
		return scanner.checkKeyword(1, 5, "hile", TOKEN_WHILE)
	case 'f':
		if scanner.current-scanner.start > 1 {
			switch scanner.source[scanner.start+1] {
			case 'a':
				return scanner.checkKeyword(2, 3, "lse", TOKEN_FALSE)
			case 'o':
				return scanner.checkKeyword(2, 2, "or", TOKEN_FOR)
			case 'u':
				return scanner.checkKeyword(2, 2, "un", TOKEN_FUN)
			}
		}
	case 't':
		if scanner.current-scanner.start > 1 {
			switch scanner.source[scanner.start+1] {
			case 'h':
				return scanner.checkKeyword(2, 2, "is", TOKEN_THIS)
			case 'r':
				return scanner.checkKeyword(2, 2, "ue", TOKEN_TRUE)
			}
		}
	}

	return TOKEN_IDENTIFIER
}

func (scanner *Scanner) advance() byte {
	scanner.current += 1
	return scanner.source[scanner.current-1]
}

func (scanner *Scanner) match(expected byte) bool {
	if scanner.isAtEnd() {
		return false
	}

	if scanner.source[scanner.current] != expected {
		return false
	}

	scanner.current += 1
	return true
}

func (scanner *Scanner) peek() byte {
	if scanner.isAtEnd() {
		return '\000'
	}
	return scanner.source[scanner.current]
}

func (scanner *Scanner) peekNext() byte {
	if scanner.current+1 >= len(scanner.source) {
		return '\000'
	}

	return scanner.source[scanner.current+1]
}

func (scanner *Scanner) scanToken() Token {
	scanner.skipWhitespace()

	scanner.start = scanner.current
	// scanner.line = 1

	if scanner.isAtEnd() {
		return scanner.makeToken(TOKEN_EOF)
	}

	c := scanner.advance()
	if scanner.isAlpha(c) {
		return scanner.identifier()
	}
	if scanner.isDigit(c) {
		return scanner.number()
	}

	switch c {
	case '(':
		return scanner.makeToken(TOKEN_LEFT_PAREN)
	case ')':
		return scanner.makeToken(TOKEN_RIGHT_PAREN)
	case '{':
		return scanner.makeToken(TOKEN_LEFT_BRACE)
	case '}':
		return scanner.makeToken(TOKEN_RIGHT_BRACE)
	case ';':
		return scanner.makeToken(TOKEN_SEMICOLON)
	case ',':
		return scanner.makeToken(TOKEN_COMMA)
	case '.':
		return scanner.makeToken(TOKEN_DOT)
	case '-':
		return scanner.makeToken(TOKEN_MINUS)
	case '+':
		return scanner.makeToken(TOKEN_PLUS)
	case '/':
		return scanner.makeToken(TOKEN_SLASH)
	case '*':
		return scanner.makeToken(TOKEN_STAR)
	case '!':
		if scanner.match('=') {
			return scanner.makeToken(TOKEN_BANG_EQUAL)
		}
		return scanner.makeToken(TOKEN_BANG)
	case '=':
		if scanner.match('=') {
			return scanner.makeToken(TOKEN_EQUAL_EQUAL)
		}
		return scanner.makeToken(TOKEN_EQUAL)
	case '<':
		if scanner.match('=') {
			return scanner.makeToken(TOKEN_LESS_EQUAL)
		}
		return scanner.makeToken(TOKEN_LESS)
	case '>':
		if scanner.match('=') {
			return scanner.makeToken(TOKEN_GREATER_EQUAL)
		}
		return scanner.makeToken(TOKEN_GREATER)
	case '"':
		return scanner.string()
	case '\000':
		return scanner.makeToken(TOKEN_EOF)
	}

	return scanner.errorToken("Unexpected character.")
}
