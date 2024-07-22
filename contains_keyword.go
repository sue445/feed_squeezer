package main

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/cockroachdb/errors"
	"strings"
)

type expression struct {
	Terms []*term `@@ { "|" @@ }` //nolint:govet
}

type term struct {
	Factors []*factor `@@ { @@ }` //nolint:govet
}

type factor struct {
	Keyword string      `  @Ident`     //nolint:govet
	Group   *expression `| "(" @@ ")"` //nolint:govet
}

var keywordLexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "Ident", Pattern: `[^\s()|]+`},
	{Name: "Punct", Pattern: `[()|]`},
	{Name: "Whitespace", Pattern: `\s+`},
})

var parser = participle.MustBuild[expression](
	participle.Lexer(keywordLexer),
	participle.Elide("Whitespace"),
)

func evaluateExpression(expr *expression, text string) bool {
	for _, term := range expr.Terms {
		if evaluateTerm(term, text) {
			return true
		}
	}
	return false
}

func evaluateTerm(term *term, text string) bool {
	for _, factor := range term.Factors {
		if !evaluateFactor(factor, text) {
			return false
		}
	}
	return true
}

func evaluateFactor(factor *factor, text string) bool {
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
	if keyword == "" {
		return true, nil
	}

	expr, err := parser.ParseString("", keyword)
	if err != nil {
		return false, errors.WithStack(err)
	}
	return evaluateExpression(expr, text), nil
}
