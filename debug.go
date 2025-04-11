package main

import (
	"log/slog"
	"os"
	"runtime"
	"runtime/debug"
)

func EnableProfiling() {
	// Optional: adjust GC target ratio for tighter memory control
	debug.SetGCPercent(100)

	// Enable block profiling
	runtime.SetBlockProfileRate(1) // 1 = sample every blocking event
	// Enable mutex profiling
	runtime.SetMutexProfileFraction(1) // 1 = sample every mutex contention
}

func setLoggingLevel(level slog.Level) {
	logLevel := &slog.LevelVar{}
	logLevel.Set(level)
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})))
}
