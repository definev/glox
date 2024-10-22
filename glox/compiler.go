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
		TOKEN_BANG:          {prefix: parser.unary, infix: nil, precedence: PREC_NONE},
		TOKEN_BANG_EQUAL:    {prefix: nil, infix: parser.binary, precedence: PREC_EQUALITY},
		TOKEN_EQUAL:         {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_EQUAL_EQUAL:   {prefix: nil, infix: parser.binary, precedence: PREC_EQUALITY},
		TOKEN_GREATER:       {prefix: nil, infix: parser.binary, precedence: PREC_COMPARISON},
		TOKEN_GREATER_EQUAL: {prefix: nil, infix: parser.binary, precedence: PREC_COMPARISON},
		TOKEN_LESS:          {prefix: nil, infix: parser.binary, precedence: PREC_COMPARISON},
		TOKEN_LESS_EQUAL:    {prefix: nil, infix: parser.binary, precedence: PREC_COMPARISON},
		TOKEN_IDENTIFIER:    {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_STRING:        {prefix: parser.string, infix: nil, precedence: PREC_NONE},
		TOKEN_NUMBER:        {prefix: parser.number, infix: nil, precedence: PREC_NONE},
		TOKEN_AND:           {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_CLASS:         {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_ELSE:          {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_FALSE:         {prefix: parser.literal, infix: nil, precedence: PREC_NONE},
		TOKEN_FOR:           {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_FUN:           {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_IF:            {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_NIL:           {prefix: parser.literal, infix: nil, precedence: PREC_NONE},
		TOKEN_OR:            {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_PRINT:         {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_RETURN:        {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_SUPER:         {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_THIS:          {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_TRUE:          {prefix: parser.literal, infix: nil, precedence: PREC_NONE},
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
		fmt.Fprintf(os.Stderr, " at '%s'", token.value)
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
	parser.complierChunk.Write(Byte, parser.current.line)
}

func (parser *Parser) emitBytes(bytes ...byte) {
	for _, Byte := range bytes {
		parser.complierChunk.Write(Byte, parser.current.line)
	}
}

func (parser *Parser) endCompiler() {
	parser.emitReturn()
	if DEBUG_PRINT_CODE {
		if !parser.hadError {
			parser.complierChunk.DisassembleChunk("code")
		}
	}
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
	parser.emitConstant(NewNumberVal(value))
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
	case TOKEN_BANG:
		parser.emitByte(OP_NOT)
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
	case TOKEN_EQUAL_EQUAL:
		parser.emitByte(OP_EQUAL)
	case TOKEN_BANG_EQUAL:
		parser.emitBytes(OP_EQUAL, OP_NOT)
	case TOKEN_GREATER:
		parser.emitByte(OP_GREATER)
	case TOKEN_GREATER_EQUAL:
		parser.emitBytes(OP_LESS, OP_NOT)
	case TOKEN_LESS:
		parser.emitByte(OP_LESS)
	case TOKEN_LESS_EQUAL:
		parser.emitBytes(OP_GREATER, OP_NOT)
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

func (parser *Parser) literal() {
	switch parser.previous.tokenType {
	case TOKEN_NIL:
		parser.emitByte(OP_NIL)
	case TOKEN_FALSE:
		parser.emitByte(OP_FALSE)
	case TOKEN_TRUE:
		parser.emitByte(OP_TRUE)
	}
}

func (parser *Parser) string() {
	str := parser.previous.value[1 : len(parser.previous.value)-1]
	parser.emitConstant(NewObjVal(NewObjString(str)))
}
