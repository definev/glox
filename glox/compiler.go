package glox

import (
	"fmt"
	"os"
	"strconv"
)

type Precedence byte

const (
	PREC_NONE       Precedence = iota
	PREC_ASSIGNMENT            // =
	PREC_OR                    // or
	PREC_AND                   // and
	PREC_EQUALITY              // == !=
	PREC_COMPARISON            // < > <= >=
	PREC_TERM                  // + -
	PREC_FACTOR                // * /
	PREC_UNARY                 // ! -
	PREC_CALL                  // . ()
	PREC_PRIMARY
)

type ParseFn func()

type ParseRule struct {
	prefix     ParseFn
	infix      ParseFn
	precedence Precedence
}

var rules map[TokenType]ParseRule

func init() {
	rules = map[TokenType]ParseRule{
		TOKEN_LEFT_PAREN:    {prefix: parser.grouping, infix: nil, precedence: PREC_NONE},
		TOKEN_RIGHT_PAREN:   {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_LEFT_BRACE:    {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_RIGHT_BRACE:   {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_COMMA:         {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_DOT:           {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_MINUS:         {prefix: parser.unary, infix: parser.binary, precedence: PREC_TERM},
		TOKEN_PLUS:          {prefix: nil, infix: parser.binary, precedence: PREC_TERM},
		TOKEN_SEMICOLON:     {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_SLASH:         {prefix: nil, infix: parser.binary, precedence: PREC_FACTOR},
		TOKEN_STAR:          {prefix: nil, infix: parser.binary, precedence: PREC_FACTOR},
		TOKEN_BANG:          {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_BANG_EQUAL:    {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_EQUAL:         {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_EQUAL_EQUAL:   {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_GREATER:       {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_GREATER_EQUAL: {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_LESS:          {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_LESS_EQUAL:    {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_IDENTIFIER:    {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_STRING:        {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_NUMBER:        {prefix: parser.number, infix: nil, precedence: PREC_NONE},
		TOKEN_AND:           {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_CLASS:         {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_ELSE:          {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_FALSE:         {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_FOR:           {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_FUN:           {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_IF:            {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_NIL:           {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_OR:            {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_PRINT:         {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_RETURN:        {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_SUPER:         {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_THIS:          {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_TRUE:          {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_VAR:           {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_WHILE:         {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_ERROR:         {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_EOF:           {prefix: nil, infix: nil, precedence: PREC_NONE},
	}
}

func getRule(tokenType TokenType) ParseRule {
	return rules[tokenType]
}

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

type Parser struct {
	previous      Token
	current       Token
	hadError      bool
	panicMode     bool
	complierChunk *Chunk
}

var parser Parser

func (parser *Parser) advance() {
	parser.previous = parser.current

	for {
		parser.current = scanner.scanToken()
		if parser.current.tokenType != TOKEN_ERROR {
			break
		}

		parser.errorAtCurrent(parser.current.value)
	}
}

func (parser *Parser) errorAtCurrent(message string) {
	parser.errorAt(&parser.current, message)
}

func (parser *Parser) error(message string) {
	parser.errorAt(&parser.previous, message)
}

func (parser *Parser) errorAt(token *Token, message string) {
	if parser.panicMode {
		return
	}
	parser.panicMode = true

	fmt.Fprintf(os.Stderr, "[line %d] Error", token.line)

	if token.tokenType == TOKEN_EOF {
		fmt.Fprintf(os.Stderr, " at end")
	} else if token.tokenType == TOKEN_ERROR {
		// Nothing.
	} else {
		fmt.Fprintf(os.Stderr, " at '%s'", token.start)
	}

	fmt.Fprintf(os.Stderr, ": %s\n", message)
}

func (parser *Parser) consume(tokenType TokenType, message string) {
	if parser.current.tokenType == tokenType {
		parser.advance()
		return
	}

	parser.errorAtCurrent(message)
}

func (parser *Parser) emitByte(Byte byte) {
	chunk := parser.complierChunk
	chunk.Write(Byte, parser.current.line)
}

func (parser *Parser) emitBytes(byte1 byte, byte2 byte) {
	parser.emitByte(byte1)
	parser.emitByte(byte2)
}

func (parser *Parser) endCompiler() {
	if DEBUG_PRINT_CODE {
		if !parser.hadError {
			parser.complierChunk.DisassembleChunk("code")
		}
	}
	parser.emitReturn()
}

func (parser *Parser) emitReturn() {
	parser.emitByte(OP_RETURN)
}

func (parser *Parser) emitConstant(value Value) {
	chunk := parser.complierChunk
	chunk.WriteConstant(value, parser.previous.line)
}

func (parser *Parser) expression() {
	parser.parsePrecedence(PREC_ASSIGNMENT)
}

func (parser *Parser) number() {
	value, _ := strconv.ParseFloat(parser.previous.value, 64)
	parser.emitConstant(Value(value))
}

func (parser *Parser) grouping() {
	parser.expression()
	parser.consume(TOKEN_RIGHT_PAREN, "Expect ')' after expression.")
}

func (parser *Parser) unary() {
	operationType := parser.previous.tokenType

	// Compile the operand.
	parser.parsePrecedence(PREC_UNARY)

	// Emit the operator instruction.
	switch operationType {
	case TOKEN_MINUS:
		parser.emitByte(OP_NEGATE)
	default:
		return // Unreachable.
	}
}

func (parser *Parser) binary() {
	operatorType := parser.previous.tokenType

	rule := getRule(operatorType)
	parser.parsePrecedence(Precedence(rule.precedence + 1))

	switch operatorType {
	case TOKEN_PLUS:
		parser.emitByte(OP_ADD)
	case TOKEN_MINUS:
		parser.emitByte(OP_SUBTRACT)
	case TOKEN_STAR:
		parser.emitByte(OP_MULTIPLY)
	case TOKEN_SLASH:
		parser.emitByte(OP_DIVIDE)
	default:
		return // Unreachable.
	}
}

func (parser *Parser) parsePrecedence(precedence Precedence) {
	parser.advance()
	prefix := getRule(parser.previous.tokenType).prefix
	if prefix == nil {
		parser.error("Expect expression.")
		return
	}
	prefix()

	for precedence <= getRule(parser.current.tokenType).precedence {
		parser.advance()
		infix := getRule(parser.previous.tokenType).infix
		infix()
	}
}
