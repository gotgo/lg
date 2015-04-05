package lg

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

var workingDir = "/"

func init() {
	wd, err := os.Getwd()
	if err != nil {
		workingDir = "unknown" //LOG IT??
	} else {
		workingDir = filepath.ToSlash(wd)
	}
}

type CallContext struct {
	LineNumber int    `json:"lineNumber"`
	FuncName   string `json:"funcName,omitempty"`
	Filename   string `json:"filename,omitempty"`
	ShortPath  string `json:"shortPath,omitempty"`
	FullPath   string `json:"fullPath,omitepmpty"`
}

func CallerContext(skipFrames int) (*CallContext, error) {
	if skipFrames < 0 {
		return nil, errors.New("negative stack frames not supported")
	}

	pc, fullPath, line, ok := runtime.Caller(skipFrames + 1)
	if !ok {
		return nil, errors.New("error during runtime.Caller")
	}

	shortPath, err := filepath.Rel(workingDir, fullPath)
	if err != nil {
		shortPath = fullPath
	}

	funcName := runtime.FuncForPC(pc).Name()

	if name, err := filepath.Rel(workingDir, funcName); err == nil {
		funcName = name
	}

	_, filename := filepath.Split(fullPath)

	caller := &CallContext{
		LineNumber: line,
		FuncName:   funcName,
		Filename:   filename,
		ShortPath:  shortPath,
		FullPath:   fullPath,
	}

	return caller, nil
}
