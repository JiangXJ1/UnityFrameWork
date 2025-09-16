package provider

import (
	"context"
	"sync"
)

type RunInfo struct {
	context  context.Context
	writer   chan string
	errChan  chan *ExportError
	errMutex sync.Mutex
	hasError bool
}

func (r *RunInfo) HasError() bool {
	r.errMutex.Lock()
	defer r.errMutex.Unlock()
	return r.hasError
}

func (r *RunInfo) SetError() {
	r.errMutex.Lock()
	defer r.errMutex.Unlock()
	r.hasError = true
}

func NewRunInfo(ctx context.Context) *RunInfo {
	return &RunInfo{
		context:  ctx,
		writer:   make(chan string, 1000),
		errChan:  make(chan *ExportError, 100),
		hasError: false,
	}
}

func (r *RunInfo) StartListen() {

}
