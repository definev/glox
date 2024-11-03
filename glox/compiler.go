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
		TOKEN_LEFT_PAREN:    {prefix: compiler.grouping, infix: nil, precedence: PREC_NONE},
		TOKEN_RIGHT_PAREN:   {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_LEFT_BRACE:    {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_RIGHT_BRACE:   {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_COMMA:         {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_DOT:           {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_MINUS:         {prefix: compiler.unary, infix: compiler.binary, precedence: PREC_TERM},
		TOKEN_PLUS:          {prefix: nil, infix: compiler.binary, precedence: PREC_TERM},
		TOKEN_SEMICOLON:     {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_SLASH:         {prefix: nil, infix: compiler.binary, precedence: PREC_FACTOR},
		TOKEN_STAR:          {prefix: nil, infix: compiler.binary, precedence: PREC_FACTOR},
		TOKEN_BANG:          {prefix: compiler.unary, infix: nil, precedence: PREC_NONE},
		TOKEN_BANG_EQUAL:    {prefix: nil, infix: compiler.binary, precedence: PREC_EQUALITY},
		TOKEN_EQUAL:         {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_EQUAL_EQUAL:   {prefix: nil, infix: compiler.binary, precedence: PREC_EQUALITY},
		TOKEN_GREATER:       {prefix: nil, infix: compiler.binary, precedence: PREC_COMPARISON},
		TOKEN_GREATER_EQUAL: {prefix: nil, infix: compiler.binary, precedence: PREC_COMPARISON},
		TOKEN_LESS:          {prefix: nil, infix: compiler.binary, precedence: PREC_COMPARISON},
		TOKEN_LESS_EQUAL:    {prefix: nil, infix: compiler.binary, precedence: PREC_COMPARISON},
		TOKEN_IDENTIFIER:    {prefix: compiler.variable, infix: nil, precedence: PREC_NONE},
		TOKEN_STRING:        {prefix: compiler.string, infix: nil, precedence: PREC_NONE},
		TOKEN_NUMBER:        {prefix: compiler.number, infix: nil, precedence: PREC_NONE},
		TOKEN_AND:           {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_CLASS:         {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_ELSE:          {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_FALSE:         {prefix: compiler.literal, infix: nil, precedence: PREC_NONE},
		TOKEN_FOR:           {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_FUN:           {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_IF:            {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_NIL:           {prefix: compiler.literal, infix: nil, precedence: PREC_NONE},
		TOKEN_OR:            {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_PRINT:         {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_RETURN:        {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_SUPER:         {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_THIS:          {prefix: nil, infix: nil, precedence: PREC_NONE},
		TOKEN_TRUE:          {prefix: compiler.literal, infix: nil, precedence: PREC_NONE},
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

type Compiler struct {
	previous      Token
	current       Token
	hadError      bool
	panicMode     bool
	complierChunk *Chunk
}

var compiler Compiler

func (compiler *Compiler) advance() {
	compiler.previous = compiler.current

	for {
		compiler.current = scanner.scanToken()
		if compiler.current.tokenType != TOKEN_ERROR {
			break
		}

		compiler.errorAtCurrent(compiler.current.value)
	}
}

func (compiler *Compiler) errorAtCurrent(message string) {
	compiler.errorAt(&compiler.current, message)
}

func (compiler *Compiler) error(message string) {
	compiler.errorAt(&compiler.previous, message)
}

func (compiler *Compiler) errorAt(token *Token, message string) {
	if compiler.panicMode {
		return
	}
	compiler.panicMode = true

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

func (compiler *Compiler) check(tokenType TokenType) bool {
	return compiler.current.tokenType == tokenType
}

func (compiler *Compiler) match(tokenType TokenType) bool {
	if !compiler.check(tokenType) {
		return false
	}
	compiler.advance()
	return true
}

func (compiler *Compiler) consume(tokenType TokenType, message string) {
	if compiler.current.tokenType == tokenType {
		compiler.advance()
		return
	}

	compiler.errorAtCurrent(message)
}

func (compiler *Compiler) emitByte(Byte byte) {
	compiler.complierChunk.Write(Byte, compiler.current.line)
}

func (compiler *Compiler) emitBytes(bytes ...byte) {
	for _, Byte := range bytes {
		compiler.complierChunk.Write(Byte, compiler.current.line)
	}
}

func (compiler *Compiler) endCompiler() {
	compiler.emitReturn()
	if DEBUG_PRINT_CODE {
		if !compiler.hadError {
			compiler.complierChunk.DisassembleChunk("code")
		}
	}
}

func (compiler *Compiler) emitReturn() {
	compiler.emitByte(OP_RETURN)
}

func (compiler *Compiler) emitConstant(value Value) {
	compiler.complierChunk.WriteConstant(value, compiler.previous.line)
}

func (compiler *Compiler) identifierConstant(token *Token) int {
	return compiler.complierChunk.AddConstant(NewObjVal(NewObjString(token.value)))
}

func (compiler *Compiler) defineVariable(global int) {
	constant0, constant1, constant2 := SplitConstant(global)
	compiler.emitBytes(OP_DEFINE_GLOBAL, constant0, constant1, constant2)
}

func (compiler *Compiler) parseVariable(errorMessage string) int {
	compiler.consume(TOKEN_IDENTIFIER, errorMessage)
	return compiler.identifierConstant(&compiler.previous)
}

func (compiler *Compiler) varDeclaration() {
	global := compiler.parseVariable("Expect variable name.")

	if compiler.match(TOKEN_EQUAL) {
		compiler.expression()
	} else {
		compiler.emitByte(OP_NIL)
	}

	compiler.consume(TOKEN_SEMICOLON, "Expect ';' after variable declaration.")
	compiler.defineVariable(global)
}

func (compiler *Compiler) declaration() {
	if compiler.match(TOKEN_VAR) {
		compiler.varDeclaration()
	} else {
		compiler.statement()
	}
	if compiler.panicMode {
		compiler.synchronize()
	}
}

func (compiler *Compiler) printStatement() {
	compiler.expression()
	compiler.consume(TOKEN_SEMICOLON, "Expect ';' after value.")
	compiler.emitByte(OP_PRINT)
}

func (compiler *Compiler) synchronize() {
	compiler.panicMode = false
	for compiler.current.tokenType != TOKEN_EOF {
		if compiler.previous.tokenType == TOKEN_SEMICOLON {
			return
		}
		switch compiler.current.tokenType {
		case TOKEN_CLASS, TOKEN_FUN, TOKEN_VAR, TOKEN_FOR, TOKEN_IF, TOKEN_WHILE, TOKEN_PRINT, TOKEN_RETURN:
			return
		}

		compiler.advance()
	}
}

func (compiler *Compiler) expressionStatement() {
	compiler.expression()
	compiler.consume(TOKEN_SEMICOLON, "Expect ';' after expression.")
	compiler.emitByte(OP_POP)
}

func (compiler *Compiler) expression() {
	compiler.parsePrecedence(PREC_ASSIGNMENT)
}

func (compiler *Compiler) statement() {
	if compiler.match(TOKEN_PRINT) {
		compiler.printStatement()
	} else {
		compiler.expressionStatement()
	}
}

func (compiler *Compiler) number() {
	value, _ := strconv.ParseFloat(compiler.previous.value, 64)
	compiler.emitConstant(NewNumberVal(value))
}

func (compiler *Compiler) grouping() {
	compiler.expression()
	compiler.consume(TOKEN_RIGHT_PAREN, "Expect ')' after expression.")
}

func (compiler *Compiler) unary() {
	operationType := compiler.previous.tokenType

	// Compile the operand.
	compiler.parsePrecedence(PREC_UNARY)

	// Emit the operator instruction.
	switch operationType {
	case TOKEN_BANG:
		compiler.emitByte(OP_NOT)
	case TOKEN_MINUS:
		compiler.emitByte(OP_NEGATE)
	default:
		return // Unreachable.
	}
}

func (compiler *Compiler) binary() {
	operatorType := compiler.previous.tokenType

	rule := getRule(operatorType)
	compiler.parsePrecedence(Precedence(rule.precedence + 1))

	switch operatorType {
	case TOKEN_PLUS:
		compiler.emitByte(OP_ADD)
	case TOKEN_MINUS:
		compiler.emitByte(OP_SUBTRACT)
	case TOKEN_STAR:
		compiler.emitByte(OP_MULTIPLY)
	case TOKEN_SLASH:
		compiler.emitByte(OP_DIVIDE)
	case TOKEN_EQUAL_EQUAL:
		compiler.emitByte(OP_EQUAL)
	case TOKEN_BANG_EQUAL:
		compiler.emitBytes(OP_EQUAL, OP_NOT)
	case TOKEN_GREATER:
		compiler.emitByte(OP_GREATER)
	case TOKEN_GREATER_EQUAL:
		compiler.emitBytes(OP_LESS, OP_NOT)
	case TOKEN_LESS:
		compiler.emitByte(OP_LESS)
	case TOKEN_LESS_EQUAL:
		compiler.emitBytes(OP_GREATER, OP_NOT)
	default:
		return // Unreachable.
	}
}

func (compiler *Compiler) parsePrecedence(precedence Precedence) {
	compiler.advance()
	prefix := getRule(compiler.previous.tokenType).prefix
	if prefix == nil {
		compiler.error("Expect expression.")
		return
	}
	prefix()

	for precedence <= getRule(compiler.current.tokenType).precedence {
		compiler.advance()
		infix := getRule(compiler.previous.tokenType).infix
		infix()
	}
}

func (compiler *Compiler) literal() {
	switch compiler.previous.tokenType {
	case TOKEN_NIL:
		compiler.emitByte(OP_NIL)
	case TOKEN_FALSE:
		compiler.emitByte(OP_FALSE)
	case TOKEN_TRUE:
		compiler.emitByte(OP_TRUE)
	}
}

func (compiler *Compiler) string() {
	str := compiler.previous.value[1 : len(compiler.previous.value)-1]
	compiler.emitConstant(NewObjVal(NewObjString(str)))
}

func (compiler *Compiler) variable() {
	compiler.namedVariable(compiler.previous)
}

func (compiler *Compiler) namedVariable(name Token) {
	nameConstant := compiler.identifierConstant(&name)
	constant0, constant1, constant2 := SplitConstant(nameConstant)
	compiler.emitBytes(OP_GET_GLOBAL, constant0, constant1, constant2)
}
