package ports

import (
	"bytes"
	"cmp"
	"errors"
	"fmt"
	"slices"

	"github.com/aereal/iter/seq"
	"github.com/aereal/mt"
	"github.com/aereal/nikki/backend/domain"
)

func asConvertMTEntryError(articleID domain.ArticleID, articleRevisionID domain.ArticleRevisionID, errs ...error) *ConvertMTEntryError {
	if len(errs) == 0 {
		return nil
	}
	accum := make([]error, 0, len(errs)+1)
	accum = append(accum, errs...)
	return &ConvertMTEntryError{
		ArticleID:         articleID,
		ArticleRevisionID: articleRevisionID,
		errs:              accum,
	}
}

type ConvertMTEntryError struct {
	ArticleID         domain.ArticleID
	ArticleRevisionID domain.ArticleRevisionID

	errs []error
}

var _ error = (*ConvertMTEntryError)(nil)

func (e *ConvertMTEntryError) Error() string {
	buf := new(bytes.Buffer)
	buf.WriteString("failed to convert MT entry: ")
	var seen bool
	for _, err := range e.errs {
		if seen {
			buf.WriteString("; ")
		}
		buf.WriteString(err.Error())
		seen = true
	}
	return buf.String()
}

func (e *ConvertMTEntryError) Unwrap() []error {
	if e == nil {
		return nil
	}
	return slices.Clone(e.errs)
}

func WrapInvalidMTExportEntryError(errs ...error) *InvalidMTExportEntryError {
	return &InvalidMTExportEntryError{errs: errs}
}

type InvalidMTExportEntryError struct {
	errs []error
}

var _ error = (*InvalidMTExportEntryError)(nil)

func (e *InvalidMTExportEntryError) Error() string {
	buf := new(bytes.Buffer)
	buf.WriteString("invalid MT export entry: ")
	var written bool
	for _, err := range e.errs {
		if written {
			buf.WriteString(", ")
		}
		buf.WriteString(err.Error())
		written = true
	}
	return buf.String()
}

func (e *InvalidMTExportEntryError) Unwrap() []error {
	return e.errs
}

func (e *InvalidMTExportEntryError) Is(other error) bool {
	rhs := new(InvalidMTExportEntryError)
	if !errors.As(other, &rhs) {
		return false
	}
	if len(e.errs) != len(rhs.errs) {
		return false
	}
	leftErrs := slices.SortedFunc(
		slices.Values(e.errs),
		cmpErr,
	)
	rightErrs := slices.SortedFunc(
		slices.Values(rhs.errs),
		cmpErr,
	)
	for a, b := range seq.Zip(slices.Values(leftErrs), slices.Values(rightErrs)) {
		if !errors.Is(a, b) {
			return false
		}
	}
	return true
}

type UnsupportedConvertBreaksError struct {
	Value mt.ConvertBreaks
}

var _ error = (*UnsupportedConvertBreaksError)(nil)

func (e *UnsupportedConvertBreaksError) Error() string {
	return fmt.Sprintf("unsupported convert breaks: %s", e.Value)
}

func (e *UnsupportedConvertBreaksError) Is(other error) bool {
	rhs := new(UnsupportedConvertBreaksError)
	if !errors.As(other, &rhs) {
		return false
	}
	return e.Value == rhs.Value
}

var ErrEmptyBasename EmptyBasenameError

type EmptyBasenameError struct{}

func (EmptyBasenameError) Error() string { return "empty Basename" }

var ErrEmptyDate EmptyDateError

type EmptyDateError struct{}

func (EmptyDateError) Error() string { return "empty Date" }

func cmpErr(a, b error) int { return cmp.Compare(a.Error(), b.Error()) }
