package validator

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"unicode"
)

var (
	reEmptyArgs = regexp.MustCompile(`\(\s*\)`)
	reSelfRef   = regexp.MustCompile(`([a-z0-9])(?P<self>\()(\)|[^\.])`)
)

// patchExprRegex rewrites the expression adding the given fieldRef to any function
// that doesn't explicitly declare one. This version uses regular expressions, uses
// less memory but takes about 5x times than patchExprScanner
//
// For instance the given the fieldRef ".SomeField" and the expression "require()",
// the output expression will become "require(.SomeField)".
func patchExprRegex(expr, fieldRef string) string {
	var (
		buf      strings.Builder
		match    []int
		group    int
		groupIdx int
		pivot    int
	)
	expr = reEmptyArgs.ReplaceAllString(expr, "$1("+fieldRef+")")
	for _, match = range reSelfRef.FindAllStringIndex(expr, -1) {
		for groupIdx, group = range match {
			if groupIdx != 1 {
				continue
			}
			buf.WriteString(expr[pivot : group-1])
			buf.WriteString(fieldRef)
			buf.WriteString(",")
			pivot = group - 1
		}
	}
	buf.WriteString(expr[pivot:])
	return buf.String()
}

// patchExprScanner rewrites the expression adding the given fieldRef to any function
// that doesn't explicitly declare one. This version uses a scanner to perform the job
// uses more memory but is about 5x faster than patchExprRegex
//
// For instance the given the fieldRef ".SomeField" and the expression "require()",
// the output expression will become "require(.SomeField)".
func patchExprScanner(expr, fieldRef string) (string, error) {
	var (
		buf strings.Builder
		ch  rune
		err error
	)
	br := bufio.NewReader(strings.NewReader(expr))
	for {
		ch, _, err = br.ReadRune()
		if errors.Is(err, io.EOF) {
			break
		}

		// start ident
		if unicode.IsLower(ch) {
			buf.WriteRune(ch)

			if err = consumeTo(br, &buf, '('); err == io.EOF {
				return "", fmt.Errorf("consume ident: unexpected EOF")
			}

			// in argument list from here
			ch, err = consumePrefixArg(br)
			if err == io.EOF {
				return "", fmt.Errorf("consume args: unexpected EOF")
			}

			// end of arg body
			if ch == ')' {
				buf.WriteString(fieldRef)
				buf.WriteRune(ch)
				continue
			}

			// existing field ref, consume and move on
			if ch == '.' {
				buf.WriteRune(ch)
			}

			// variable declaration
			if ch == '\'' || unicode.IsDigit(ch) {
				buf.WriteString(fieldRef)
				buf.WriteString(",")
				buf.WriteRune(ch)
			}

			if err = consumeTo(br, &buf, ')'); err == io.EOF {
				return "", fmt.Errorf("consume ref args: unexpected EOF")
			}
		}

		if ch == '(' || ch == ')' || ch == ' ' || ch == '\t' || ch == '&' || ch == '|' {
			buf.WriteRune(ch)
			continue
		}
	}

	return buf.String(), nil
}

// consumeTo consumes the reader up to and including the first instance of the given stopCh rune is found.
func consumeTo(in *bufio.Reader, out *strings.Builder, stopCh rune) error {
	var (
		ch  rune
		err error
	)
	for {
		// consume ident up to open left parenthesis
		ch, _, err = in.ReadRune()
		if err == io.EOF {
			return err
		}
		out.WriteRune(ch)
		if ch == stopCh {
			break
		}
	}
	return nil
}

// consumePrefixArg consumes the reader from the beginning of the argument body to either the end of the arg
// body or to a valid argument. It is used when rewriting the expression to remove trailing whitespaces.
func consumePrefixArg(in *bufio.Reader) (rune, error) {
	var (
		ch  rune
		err error
	)
	for {
		ch, _, err = in.ReadRune()
		if err == io.EOF {
			return 0, nil
		}
		if ch == ')' || ch == '.' || ch == ',' || ch == '\'' || unicode.IsDigit(ch) {
			return ch, nil
		}
	}
}
