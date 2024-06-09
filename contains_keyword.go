package main

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"strings"
)

type Expression struct {
	Terms []*Term `@@ { "|" @@ }`
}

type Term struct {
	Factors []*Factor `@@ { @@ }`
}

type Factor struct {
	Keyword string      `  @Ident`
	Group   *Expression `| "(" @@ ")"`
}

var keywordLexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "Ident", Pattern: `[a-zA-Z0-9_]+`},
	{Name: "Punct", Pattern: `[()|]`},
	{Name: "Whitespace", Pattern: `\s+`},
})

var parser = participle.MustBuild[Expression](
	participle.Lexer(keywordLexer),
	participle.Elide("Whitespace"),
)

func evaluateExpression(expr *Expression, text string) bool {
	for _, term := range expr.Terms {
		if evaluateTerm(term, text) {
			return true
		}
	}
	return false
}

func evaluateTerm(term *Term, text string) bool {
	for _, factor := range term.Factors {
		if !evaluateFactor(factor, text) {
			return false
		}
	}
	return true
}

func evaluateFactor(factor *Factor, text string) bool {
	if factor.Keyword != "" {
		return strings.Contains(strings.ToLower(text), strings.ToLower(factor.Keyword))
	}
	if factor.Group != nil {
		return evaluateExpression(factor.Group, text)
	}
	return false
}

// ContainsKeyword returns whether the passed text matches by keyword
func ContainsKeyword(text string, keyword string) (bool, error) {
	expr, err := parser.ParseString("", keyword)
	if err != nil {
		return false, err
	}
	return evaluateExpression(expr, text), nil
}
