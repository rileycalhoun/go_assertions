package assert

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime/debug"
)

// Much of this code has been completely stolen from ThePrimagen!
type AssertData interface {
	Dump() string
}

type AssertFlush interface {
	Flush()
}

var flushes []AssertFlush = []AssertFlush{}
var assertData map[string]AssertData = map[string]AssertData{}
var writer io.Writer

func AddAssertData(key string, value AssertData) {
	assertData[key] = value
}

func RemoveAssertData(key string) {
	delete(assertData, key)
}

func AddAssertFlush(flusher AssertFlush) {
	flushes = append(flushes, flusher)
}

func ToWriter(w io.Writer) {
	writer = w
}

func assertFail(msg string, args ...interface{}) {

	for _, f := range flushes {
		f.Flush()
	}

	slogValues := []interface{} {
		"msg",
		msg,
		"area",
		"Assert",
	}

	slogValues = append(slogValues, args...)
	fmt.Fprintf(os.Stderr, "ARGS: %+v\n", args)
	for k, v := range assertData {
		slogValues = append(slogValues, k, v.Dump())
	}

	fmt.Fprintf(os.Stderr, "ASSERT\n")
	for i := 0; i < len(slogValues); i += 2 {
		fmt.Fprintf(os.Stderr, "  %s=%v\n", slogValues[i], slogValues[i+1])
	}

	fmt.Fprintln(os.Stderr, string(debug.Stack()))
	os.Exit(1)
}

func Assert(truth bool, msg string, data ...any) {
	if !truth {
		assertFail(msg, data...)
	}
}

func NotEqual(unexpected any, actual any, msg string, data ...any) {
	slog.Info("Neq Check", "unexpected", unexpected, "actual", actual)
	if unexpected != actual {
		return
	}

	slog.Error("Neq#eq encountered")
	assertFail(msg, data)
}

func Nil(item any, msg string, data ...any) {
	slog.Info("Nil Check", "item", item)
	if item == nil {
		return
	}

	slog.Error("Nil#not nil enountered")
	assertFail(msg, data)
}

func NotNil(item any, msg string, data ...any) {
	slog.Info("Not Nil Check", "item", item)
	if item != nil {
		return
	}

	slog.Error("Not Nil#nil enountered")
	assertFail(msg, data)
}

func NoError(err error, msg string, data ...any) {
	if err != nil {
		data = append(data, "error", err)
		assertFail(msg, data...)
	}
}
