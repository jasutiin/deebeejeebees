package lexer

import (
	"strings"
	"unicode"

	"github.com/jasutiin/deebeejeebees/internal/tokens"
)

// AnalyzeString tokenizes the input SQL string and returns the list of tokens.
func AnalyzeString(input string) []string {
    var parsedTokens []string

    i := 0
    for i < len(input) {
        character := rune(input[i])

        if unicode.IsSpace(character) {
            i++
            continue
        }

        if character == '\'' {
            str, newIndex := analyzeStringLiteralToken(input, i)
            parsedTokens = append(parsedTokens, str)
            i = newIndex + 1
            continue
        }

        symbol, newIndex, ok := tryAnalyzeSymbol(input, i);

        if ok {
            parsedTokens = append(parsedTokens, symbol)
            i = newIndex + 1
            continue
        }

        if unicode.IsLetter(character) || character == '_' {
            start := i

            // read the entire word
            for i < len(input) && (unicode.IsLetter(rune(input[i])) || unicode.IsDigit(rune(input[i])) || rune(input[i]) == '_') {
                i++
            }

            word := strings.ToUpper(input[start:i]) // all of the reserved keywords are uppercase

            // if it's a reserved keyword, add it. if not, then take its original form and add it
            if tokens.ReservedKeywords[word] {
                parsedTokens = append(parsedTokens, word)
            } else {
                parsedTokens = append(parsedTokens, input[start:i])
            }

            continue
        }

        // checks for numbers. the reason why this is separate from the logic above is because it takes into 
        // account identifiers that start with numbers: '1age', which is not allowed in SQL.
        // identifiers must start with a letter.
        if unicode.IsDigit(character) {
            start := i

            for i < len(input) && unicode.IsDigit(rune(input[i])) {
                i++
            }

            parsedTokens = append(parsedTokens, input[start:i])
            continue
        }

        i++
    }

    return parsedTokens
}

// tryAnalyzeSymbol takes in the input string and the index we were on to check if the symbol is a reserved symbol.
// It also checks if it is a double symbol like '!=' or '<>'.
func tryAnalyzeSymbol(input string, index int) (string, int, bool) {
    if index + 1 < len(input) {
        two := input[index : index + 2]
        if _, ok := tokens.ReservedSymbols[two]; ok {
            return two, index + 1, true
        }
    }

    one := string(input[index])
    if _, ok := tokens.ReservedSymbols[one]; ok {
        return one, index, true
    }
    return "", index, false
}

// analyzeStringLiteralToken takes in the input string and the index we were on to check for the entire sentence
// inside the string literals.
func analyzeStringLiteralToken(input string, index int) (string, int) {
    var sb strings.Builder
    sb.WriteRune('\'')
    i := index + 1
    for i < len(input) {
        ch := rune(input[i])
        sb.WriteRune(ch)
        if ch == '\'' {
            return sb.String(), i
        }
        i++
    }

    return sb.String(), i - 1
}