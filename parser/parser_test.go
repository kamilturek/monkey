package parser_test

import (
	"strconv"
	"testing"

	"github.com/kamilturek/monkey/ast"
	"github.com/kamilturek/monkey/lexer"
	"github.com/kamilturek/monkey/parser"
)

func testLetStatement(t *testing.T, stmt ast.Statement, name string) bool {
	t.Helper()

	if stmt.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", stmt.TokenLiteral())

		return false
	}

	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", stmt)

		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)

		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s", name, letStmt.Name.TokenLiteral())

		return false
	}

	return true
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, value int64) bool {
	t.Helper()

	il, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("exp not *ast.IntegerLiteral. got=%T", exp)

		return false
	}

	if il.Value != value {
		t.Errorf("il.Value not %d. got=%d", value, il.Value)

		return false
	}

	if il.TokenLiteral() != strconv.FormatInt(value, 10) {
		t.Errorf("il.TokenLiteral not %d. got=%s", value, il.TokenLiteral())
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	t.Helper()

	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)

		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)

		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())

		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	t.Helper()

	ident, ok := exp.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)

		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %t. got=%t", value, ident.Value)

		return false
	}

	if ident.TokenLiteral() != strconv.FormatBool(value) {
		t.Errorf("ident.TokenLiteral not %t. got=%s", value, ident.TokenLiteral())

		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	t.Helper()

	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}

	t.Errorf("type of exp not handled. got=%T", exp)

	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	t.Helper()

	ie, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T", exp)

		return false
	}

	if !testLiteralExpression(t, ie.Left, left) {
		return false
	}

	if ie.Operator != operator {
		return false
	}

	if !testLiteralExpression(t, ie.Right, right) {
		return false
	}

	return true
}

func testParserErrors(t *testing.T, p *parser.Parser) {
	t.Helper()

	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}

func testProgramStatements(t *testing.T, program *ast.Program, expectedCount int) bool {
	t.Helper()

	if len(program.Statements) != expectedCount {
		t.Errorf("program.Statements does not contain %d statement. got=%d", expectedCount, (program.Statements))

		return false
	}

	return true
}

func TestLetStatements(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = 10;", "y", 10},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := parser.NewParser(l)
		program := p.ParseProgram()
		testParserErrors(t, p)

		if !testProgramStatements(t, program, 1) {
			return
		}

		stmt, ok := program.Statements[0].(*ast.LetStatement)
		if !ok {
			t.Fatalf("program.Statements[0] not *ast.LetStatement. got=%T", program.Statements[0])
		}

		val := stmt.Value

		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestReturnStatement(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return 10;", 10},
		{"return y;", "y"},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := parser.NewParser(l)
		program := p.ParseProgram()
		testParserErrors(t, p)

		if !testProgramStatements(t, program, 1) {
			return
		}

		stmt, ok := program.Statements[0].(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("program.Statements[0] not *ast.ReturnStatement. got=%T", program.Statements[0])
		}

		val := stmt.Value

		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	t.Parallel()

	input := "foobar;"

	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	testParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	testLiteralExpression(t, stmt.Expression, "foobar")
}

func TestIntegerLiteralExpression(t *testing.T) {
	t.Parallel()

	input := "5;"

	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	testParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	testLiteralExpression(t, stmt.Expression, 5)
}

func TestParsingPrefixExpressions(t *testing.T) {
	t.Parallel()

	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.NewLexer(tt.input)
		p := parser.NewParser(l)
		program := p.ParseProgram()
		testParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.PrefixExpression. got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		testLiteralExpression(t, exp.Right, tt.value)
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	t.Parallel()

	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "==", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.NewLexer(tt.input)
		p := parser.NewParser(l)
		program := p.ParseProgram()
		testParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue)
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 > 4 != 3 < 4",
			"((5 > 4) != (3 < 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"a + - b",
			"(a + (-b))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := parser.NewParser(l)
		program := p.ParseProgram()
		testParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestBoolean(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := parser.NewParser(l)
		program := p.ParseProgram()
		testParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		testLiteralExpression(t, stmt.Expression, tt.expected)
	}
}

func TestIfExpression(t *testing.T) {
	t.Parallel()

	input := "if (x < y) { x }"

	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	testParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression not *ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("exp.Consequence.Statements is not 1 statements. got=%d", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Consequence.Statements[0] not an *ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative was not nil. got=%+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	t.Parallel()

	input := "if (x < y) { x } else { y }"

	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	testParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression not *ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("exp.Consequence.Statements is not 1 statements. got=%d", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Consequence.Statements[0] not an *ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.alternative.Statements[0] not an *ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestFunctionLiteral(t *testing.T) {
	t.Parallel()

	input := "fn(x, y) { x + y; }"

	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	testParserErrors(t, p)

	if !testProgramStatements(t, program, 1) {
		return
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression not *ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function.Parameters does not contain 2 parameters. got=%d", len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements does not contain 1 statement. got=%d", len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function.Body.Statements[0] not *ast.ExpressionStatement. got=%T", function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameters(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input          string
		expectedParams []string
	}{
		{"fn() {};", []string{}},
		{"fn(x) {};", []string{"x"}},
		{"fn(x, y) {};", []string{"x", "y"}},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := parser.NewParser(l)
		program := p.ParseProgram()
		testParserErrors(t, p)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		function, ok := stmt.Expression.(*ast.FunctionLiteral)
		if !ok {
			t.Fatalf("stmt.Expression not *ast.FunctionLiteral. got=%T", stmt.Expression)
		}

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Fatalf(
				"function.Parameters does not contain %d parameters. got=%d",
				len(tt.expectedParams),
				len(function.Parameters),
			)
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpression(t *testing.T) {
	t.Parallel()

	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	testParserErrors(t, p)

	if !testProgramStatements(t, program, 1) {
		return
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	callExp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression not *ast.CallExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, callExp.Function, "add") {
		return
	}

	if len(callExp.Arguments) != 3 {
		t.Fatalf("callExp.Arguments does not contain 3 arguments. got=%d", len(callExp.Arguments))
	}

	testLiteralExpression(t, callExp.Arguments[0], 1)
	testInfixExpression(t, callExp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, callExp.Arguments[2], 4, "+", 5)
}

func TestCallExpressionArguments(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input             string
		expectedArguments []string
	}{
		{"add();", []string{}},
		{"add(x);", []string{"x"}},
		{"add(x, y);", []string{"x", "y"}},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := parser.NewParser(l)
		program := p.ParseProgram()
		testParserErrors(t, p)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		callExp, ok := stmt.Expression.(*ast.CallExpression)
		if !ok {
			t.Fatalf("stmt.Expression not *ast.CallExpression. got=%T", stmt.Expression)
		}

		if len(callExp.Arguments) != len(tt.expectedArguments) {
			t.Fatalf(
				"callExp.Arguments does not contain %d arguments. got=%d",
				len(tt.expectedArguments),
				len(callExp.Arguments),
			)
		}

		for i, ident := range tt.expectedArguments {
			testLiteralExpression(t, callExp.Arguments[i], ident)
		}
	}
}
