package funcs

import (
	"time"

	"github.com/go-openapi/strfmt"
)

func DiffStrfmtDate(x, y strfmt.Date) map[string][]interface{} {
	return DiffTimeTime(time.Time(x), time.Time(y))
}

func DiffStrfmtDateTime(x, y strfmt.DateTime) map[string][]interface{} {
	return DiffTimeTime(time.Time(x), time.Time(y))
}

func DiffTimeTime(x, y time.Time) map[string][]interface{} {
	out := make(map[string][]interface{})

	xu := x.UTC()
	yu := y.UTC()

	if xu.Year() != yu.Year() {
		out["Year"] = []interface{}{xu.Year(), yu.Year()}
	}
	if xu.Month() != yu.Month() {
		out["Month"] = []interface{}{xu.Month(), yu.Month()}
	}
	if xu.Day() != yu.Day() {
		out["Day"] = []interface{}{xu.Day(), yu.Day()}
	}
	if xu.Hour() != yu.Hour() {
		out["Hour"] = []interface{}{xu.Hour(), yu.Hour()}
	}
	if xu.Minute() != yu.Minute() {
		out["Minute"] = []interface{}{xu.Minute(), yu.Minute()}
	}
	if xu.Second() != yu.Second() {
		out["Second"] = []interface{}{xu.Second(), yu.Second()}
	}
	if xu.Nanosecond() != yu.Nanosecond() {
		out["Nanosecond"] = []interface{}{xu.Nanosecond(), yu.Nanosecond()}
	}

	return out
}
