package funcs

import (
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/haproxytech/go-method-gen/pkg/eqdiff"
)

func DiffStrfmtDate(x, y strfmt.Date, opts ...eqdiff.GoMethodGenOptions) map[string][]any {
	return DiffTimeTime(time.Time(x), time.Time(y), opts...)
}

func DiffStrfmtDateTime(x, y strfmt.DateTime, opts ...eqdiff.GoMethodGenOptions) map[string][]any {
	return DiffTimeTime(time.Time(x), time.Time(y), opts...)
}

func DiffTimeTime(x, y time.Time, _ ...eqdiff.GoMethodGenOptions) map[string][]any {
	out := make(map[string][]any)

	xu := x.UTC()
	yu := y.UTC()

	if xu.Year() != yu.Year() {
		out["Year"] = []any{xu.Year(), yu.Year()}
	}
	if xu.Month() != yu.Month() {
		out["Month"] = []any{xu.Month(), yu.Month()}
	}
	if xu.Day() != yu.Day() {
		out["Day"] = []any{xu.Day(), yu.Day()}
	}
	if xu.Hour() != yu.Hour() {
		out["Hour"] = []any{xu.Hour(), yu.Hour()}
	}
	if xu.Minute() != yu.Minute() {
		out["Minute"] = []any{xu.Minute(), yu.Minute()}
	}
	if xu.Second() != yu.Second() {
		out["Second"] = []any{xu.Second(), yu.Second()}
	}
	if xu.Nanosecond() != yu.Nanosecond() {
		out["Nanosecond"] = []any{xu.Nanosecond(), yu.Nanosecond()}
	}

	return out
}
