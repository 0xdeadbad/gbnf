package lexer

import (
	"bytes"
	"io"
	"log"
	"os"
	"testing"
)

// Tests the lexer's ability to read a sequence of characters correctly
func TestNewLexer(t *testing.T) {
	buffer := []byte("Hello, World!")
	lexer := NewLexer(bytes.NewReader(buffer))

	for i := 0; i < len(buffer); i++ {
		r, err := lexer.nextChar()
		if err != nil {
			t.Fatal(err)
		}
		if r != rune(buffer[i]) {
			t.Fatalf("Expected %c, got %c", buffer[i], r)
		}
	}
}

// Testing if the lexer is able to emit the token Symbol correctly
func TestLexer_NextToken_Symbol(t *testing.T) {
	buffer := []byte("<symbol>")
	lexer := NewLexer(bytes.NewReader(buffer))

	token, err := lexer.NextToken()
	if err != nil {
		if err != io.EOF {
			t.Fatal(err)
		}
	}

	if token.Lexeme != "symbol" {
		t.Fatalf("Expected symbol, got %s", token.Lexeme)
	}

	t.Logf("Token: %s", token)
}

// Testing if the lexer is able to emit the token ProdRule correctly
func TestLexer_NextToken_Assignment(t *testing.T) {
	buffer := []byte("::=")
	lexer := NewLexer(bytes.NewReader(buffer))

	token, err := lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}

	if token.Lexeme != "::=" {
		t.Fatalf("Expected ::=, got %s", token.Lexeme)
	}

	t.Logf("Token: %s", token)
}

// Testing if the lexer is able to emit the token Or correctly
func TestLexer_NextToken_Or(t *testing.T) {
	buffer := []byte("|")
	lexer := NewLexer(bytes.NewReader(buffer))

	token, err := lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}

	if token.Lexeme != "|" {
		t.Fatalf("Expected |, got %s", token.Lexeme)
	}

	t.Logf("Token: %s", token)
}

// Testing if the lexer is able to ignore whitespaces
func TestLexer_NextToken_Whitespace(t *testing.T) {
	buffer := []byte("       	\n\n<abc>")
	lexer := NewLexer(bytes.NewReader(buffer))

	token, err := lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}

	if token.Lexeme != "abc" {
		t.Fatalf("Expected abc, got %s", token.Lexeme)
	}

	t.Logf("Token: %s", token)
}

// Testing if the lexer is able to ignore whitespaces
func TestLexer_NextToken_OneWhitespace(t *testing.T) {
	buffer := []byte(" <abc>")
	lexer := NewLexer(bytes.NewReader(buffer))

	token, err := lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}

	if token.Lexeme != "abc" {
		t.Fatalf("Expected abc, got %s", token.Lexeme)
	}

	t.Logf("Token: %s", token)
}

// Testing if the lexer is able to emit the token String correctly
func TestLexer_NextToken_String(t *testing.T) {
	buffer := []byte("\"Hello, World!\" 'Hello, World!'")
	lexer := NewLexer(bytes.NewReader(buffer))

	token, err := lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	if token.Lexeme != "Hello, World!" {
		t.Fatalf("Expected Hello, World!, got %s", token.Lexeme)
	}
	if token.Type != TerminalSymbol {
		t.Fatalf("Expected type 'String', got '%s'", token.Type)
	}
	t.Logf("Token: %s", token)

	token, err = lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	if token.Lexeme != "Hello, World!" {
		t.Fatalf("Expected Hello, World!, got %s", token.Lexeme)
	}
	if token.Type != TerminalSymbol {
		t.Fatalf("Expected type 'String', got '%s'", token.Type)
	}
	t.Logf("Token: %s", token)
}

// Testing if the lexer is able to emit the token String correctly
func TestLexer_NextToken_Action(t *testing.T) {
	buffer := []byte("{Fn(Hello, World)}")
	lexer := NewLexer(bytes.NewReader(buffer))

	token, err := lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	if token.Lexeme != "Fn" {
		t.Fatalf("Expected Fn, got %s", token.Lexeme)
	}

	/*token, err = lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	if token.Lexeme != "(" {
		t.Fatalf("Expected (, got %s", token.Lexeme)
	}*/

	token, err = lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	if token.Lexeme != "Hello" {
		t.Fatalf("Expected Hello, got %s", token.Lexeme)
	}

	token, err = lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	if token.Lexeme != "World" {
		t.Fatalf("Expected World, got %s", token.Lexeme)
	}

	/*token, err = lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	if token.Lexeme != ")" {
		t.Fatalf("Expected ), got %s", token.Lexeme)
	}*/

	t.Logf("Tokens: %s", lexer.Tokens)
}

// Testing if the lexer is able to ignore whitespaces between and newline characters between tokens
func TestLexer_NextToken_Newline(t *testing.T) {
	buffer := []byte("\"Hello World!\"\n<term>\n   	<abc>\n\n::=\nwhile | true | false")
	lexer := NewLexer(bytes.NewReader(buffer))

	token, err := lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}

	if token.Lexeme != "Hello World!" {
		t.Fatalf("Expected Hello World!, got %s", token.Lexeme)
	}
	if token.Type != TerminalSymbol {
		t.Fatalf("Expected type 'String', got '%s'", token.Type)
	}
	t.Logf("Token: %s", token)

	token, err = lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}

	if token.Lexeme != "term" {
		t.Fatalf("Expected term, got %s", token.Lexeme)
	}
	if token.Type != NonTerminalSymbol {
		t.Fatalf("Expected type 'NonTerminalSymbol', got '%s'", token.Type)
	}
	t.Logf("Token: %s", token)

	token, err = lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}

	if token.Lexeme != "abc" {
		t.Fatalf("Expected abc, got %s", token.Lexeme)
	}
	if token.Type != NonTerminalSymbol {
		t.Fatalf("Expected type 'NonTerminalSymbol', got '%s'", token.Type)
	}
	t.Logf("Token: %s", token)

	token, err = lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}

	if token.Lexeme != "::=" {
		t.Fatalf("Expected ::=, got %s", token.Lexeme)
	}
	if token.Type != ProdRule {
		t.Fatalf("Expected type 'ProdRule', got '%s'", token.Type)
	}
	t.Logf("Token: %s", token)

	token, err = lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}

	if token.Lexeme != "while" {
		t.Fatalf("Expected while, got %s", token.Lexeme)
	}
	if token.Type != TerminalSymbol {
		t.Fatalf("Expected type 'TerminalSymbol', got '%s'", token.Type)
	}
	t.Logf("Token: %s", token)

	token, err = lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}

	if token.Lexeme != "|" {
		t.Fatalf("Expected |, got %s", token.Lexeme)
	}
	if token.Type != Or {
		t.Fatalf("Expected type 'Or', got '%s'", token.Type)
	}
	t.Logf("Token: %s", token)

	token, err = lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}

	if token.Lexeme != "true" {
		t.Fatalf("Expected true, got %s", token.Lexeme)
	}
	if token.Type != TerminalSymbol {
		t.Fatalf("Expected type 'TerminalSymbol', got '%s'", token.Type)
	}
	t.Logf("Token: %s", token)

	token, err = lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}

	if token.Lexeme != "|" {
		t.Fatalf("Expected |, got %s", token.Lexeme)
	}
	if token.Type != Or {
		t.Fatalf("Expected type 'Or', got '%s'", token.Type)
	}
	t.Logf("Token: %s", token)

	token, err = lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}

	if token.Lexeme != "false" {
		t.Fatalf("Expected false, got %s", token.Lexeme)
	}
	if token.Type != TerminalSymbol {
		t.Fatalf("Expected type 'TerminalSymbol', got '%s'", token.Type)
	}
	t.Logf("Token: %s", token)
}

// Testing if the lexer is able to emit the tokens correctly of a correct syntax
func TestLexer_NextToken_EntireLine(t *testing.T) {
	var err error

	buffer := []byte("<expr> ::= <term> \"+\" <expr> |  <term>\n<expr>   ::=   <term> \"+\"  <expr> |        <term>")
	lexer := NewLexer(bytes.NewReader(buffer))

	for i := 0; i < 14; i++ {
		_, err = lexer.NextToken()
		if err != nil {
			t.Fatal(err)
		}
	}

	if len(lexer.Tokens) != 14 {
		t.Fatalf("Expected 14 token, got %d", len(lexer.Tokens))
	}

	// High quality piece of code
	// I'm too lazy to write a loop to check all the tokens
	if lexer.Tokens[0].Lexeme != "expr" || lexer.Tokens[7].Lexeme != "expr" {
		t.Fatalf("Expected expr, got %s, %s", lexer.Tokens[0].Lexeme, lexer.Tokens[7].Lexeme)
	}
	if lexer.Tokens[1].Lexeme != "::=" || lexer.Tokens[8].Lexeme != "::=" {
		t.Fatalf("Expected ::=, got %s, %s", lexer.Tokens[1].Lexeme, lexer.Tokens[8].Lexeme)
	}
	if lexer.Tokens[2].Lexeme != "term" || lexer.Tokens[9].Lexeme != "term" {
		t.Fatalf("Expected term, got %s, %s", lexer.Tokens[2].Lexeme, lexer.Tokens[9].Lexeme)
	}
	if lexer.Tokens[3].Lexeme != "+" || lexer.Tokens[10].Lexeme != "+" {
		t.Fatalf("Expected +, got %s, %s", lexer.Tokens[3].Lexeme, lexer.Tokens[10].Lexeme)
	}
	if lexer.Tokens[4].Lexeme != "expr" || lexer.Tokens[11].Lexeme != "expr" {
		t.Fatalf("Expected expr, got %s, %s", lexer.Tokens[4].Lexeme, lexer.Tokens[11].Lexeme)
	}
	if lexer.Tokens[5].Lexeme != "|" || lexer.Tokens[12].Lexeme != "|" {
		t.Fatalf("Expected |, got %s, %s", lexer.Tokens[5].Lexeme, lexer.Tokens[12].Lexeme)
	}
	if lexer.Tokens[6].Lexeme != "term" || lexer.Tokens[13].Lexeme != "term" {
		t.Fatalf("Expected term, got %s, %s", lexer.Tokens[6].Lexeme, lexer.Tokens[13].Lexeme)
	}

	t.Logf("Tokens: %v", lexer.Tokens)
}

func TestLexer_NextToken_FromFile(t *testing.T) {
	// Create a temporary file, so we can test the lexer with a real file instead of an in memory buffer
	f, err := os.Create("test_bnf_07ce.txt")
	if err != nil {
		t.Fatal(err)
	}

	// Deferred function to close and remove the file after the test ends
	defer func(f *os.File) {
		// Close and remove file after the function ends
		name := f.Name()
		err := f.Close()
		if err != nil {
			log.Println(err)
		}
		err = os.Remove(name)
		if err != nil {
			log.Println(err)
		}
	}(f)

	data := []byte("<expr> ::= <term> \"+\" <expr> |  <term>\n<expr>   ::=   <term> \"+\"  <expr> |        <term>")

	// Here write moves the file pointer to the end of the written data
	n, err := f.Write(data)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(data) {
		t.Fatalf("Expected %d bytes written, got %d", len(data), n)
	}

	// Reset file pointer to the beginning of the file because writing moves the pointer to the end
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		t.Fatal(err)
	}

	lexer := NewLexer(f)

	// Reusing the same piece of code from TestLexer_NextToken_EntireLine because I'm lazy
	for i := 0; i < 14; i++ {
		_, err = lexer.NextToken()
		if err != nil {
			if err == io.EOF {
				continue
			}
		}
	}

	if len(lexer.Tokens) != 14 {
		t.Fatalf("Expected 14 token, got %d", len(lexer.Tokens))
	}

	t.Logf("Tokens: %v", lexer.Tokens)
}

func TestLexer_NextToken_Assign(t *testing.T) {
	buffer := []byte("=")
	lexer := NewLexer(bytes.NewReader(buffer))

	token, err := lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Token: %s", token)
}

func TestLexer_NextToken_LineWithAssign(t *testing.T) {
	buffer := []byte("<expr> ::= x=<expr> + y=<term> { Add(x, y) }\n\t| <term>\n<term> ::= NAME")
	lexer := NewLexer(bytes.NewReader(buffer))

	token, err := lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	for token != nil {
		t.Logf("Token: %s", token)
		token, err = lexer.NextToken()
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Logf("Tokens: %v", lexer.Tokens)
}

func TestLexer_NextToken_Sequence(t *testing.T) {
	buffer := []byte("...")
	lexer := NewLexer(bytes.NewReader(buffer))

	token, err := lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}

	if token.Lexeme != "..." {
		t.Fatalf("Expected ..., got %s", token.Lexeme)
	}
	if token.Type != Sequence {
		t.Fatalf("Expected type 'Sequence', got '%s'", token.Type)
	}

	t.Logf("Token: %s", token)
}

func TestLexer_CommonLispBNF(t *testing.T) {
	buffer := []byte(`
<s_expression> ::= <atomic_symbol>
	| "(" <s_expression> "."<s_expression> ")"
	| <list>
<list> ::= "(" <s_expression> "<" <s_expression> ">" ")"
<atomic_symbol> ::= <letter> <atom_part>
<atom_part> ::= <empty> | <letter> <atom_part> | <number> <atom_part>
<letter> ::= "a" ... "z"
<number> ::= "1" ... "9"
<empty> ::= " "`)

	lexer := NewLexer(bytes.NewReader(buffer))

	token, err := lexer.NextToken()
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}
	t.Logf("Token: %s", token)
	for token != nil {
		token, err = lexer.NextToken()
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Fatal(err)
		}
		t.Logf("Token: %s", token)
	}

	t.Logf("Tokens: %v", lexer.Tokens)
}

func TestLexer_NextToken_MarkReset(t *testing.T) {
	buffer := []byte("<expr> ::= <term> \"+\" <expr> |  <term>")
	lexer := NewLexer(bytes.NewReader(buffer))

	token, err := lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Token: %s", token)

	mark := lexer.Mark()
	token, err = lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Token: %s (Marked)", token)

	token, err = lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Token: %s", token)

	t.Logf("Resetting to mark")
	lexer.Reset(mark)
	token, err = lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Token: %s", token)

	token, err = lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Token: %s", token)

}

func TestLexer_PeekToken(t *testing.T) {
	buffer := []byte("<expr> ::= <term> \"+\" <expr> |  <term>")
	lexer := NewLexer(bytes.NewReader(buffer))

	token, err := lexer.PeekToken()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Peeked Token: %s", token)

	token, err = lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Token: %s", token)

	mark := lexer.Mark()
	token, err = lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Token: %s (Marked)", token)

	token, err = lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Token: %s", token)

	t.Logf("Resetting to mark")
	lexer.Reset(mark)
	token, err = lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Token: %s", token)

	token, err = lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Token: %s", token)

	token, err = lexer.PeekToken()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Peeked Token: %s", token)

	token, err = lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Token: %s", token)

}

func TestLexer_NextToken_EndOfRule(t *testing.T) {
	buffer := []byte("<expr> ::= <term> \"+\" <expr> |  <term>!!")
	lexer := NewLexer(bytes.NewReader(buffer))

	token, err := lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Token: %s", token)

	for token.Type != EndOfRule {
		token, err = lexer.NextToken()
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Token: %s", token)
	}
}

func TestLexer_NextToken_GroupsAndNot(t *testing.T) {
	buffer := []byte("<expr> ::= !(<term> & \"+\" & <expr>) |  [<term> & + & \"test\" & ![a-Z]]!!")
	lexer := NewLexer(bytes.NewReader(buffer))

	token, err := lexer.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Token: %s", token)

	for token.Type != EndOfRule {
		token, err = lexer.NextToken()
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Token: %s", token)
	}
}
