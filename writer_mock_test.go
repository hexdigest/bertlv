package bertlv

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

//go:generate minimock -i io.Writer -o ./writer_mock_test.go

import (
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
)

// WriterMock implements io.Writer
type WriterMock struct {
	t minimock.Tester

	funcWrite          func(p []byte) (n int, err error)
	inspectFuncWrite   func(p []byte)
	afterWriteCounter  uint64
	beforeWriteCounter uint64
	WriteMock          mWriterMockWrite
}

// NewWriterMock returns a mock for io.Writer
func NewWriterMock(t minimock.Tester) *WriterMock {
	m := &WriterMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.WriteMock = mWriterMockWrite{mock: m}
	m.WriteMock.callArgs = []*WriterMockWriteParams{}

	return m
}

type mWriterMockWrite struct {
	mock               *WriterMock
	defaultExpectation *WriterMockWriteExpectation
	expectations       []*WriterMockWriteExpectation

	callArgs []*WriterMockWriteParams
	mutex    sync.RWMutex
}

// WriterMockWriteExpectation specifies expectation struct of the Writer.Write
type WriterMockWriteExpectation struct {
	mock    *WriterMock
	params  *WriterMockWriteParams
	results *WriterMockWriteResults
	Counter uint64
}

// WriterMockWriteParams contains parameters of the Writer.Write
type WriterMockWriteParams struct {
	p []byte
}

// WriterMockWriteResults contains results of the Writer.Write
type WriterMockWriteResults struct {
	n   int
	err error
}

// Expect sets up expected params for Writer.Write
func (mmWrite *mWriterMockWrite) Expect(p []byte) *mWriterMockWrite {
	if mmWrite.mock.funcWrite != nil {
		mmWrite.mock.t.Fatalf("WriterMock.Write mock is already set by Set")
	}

	if mmWrite.defaultExpectation == nil {
		mmWrite.defaultExpectation = &WriterMockWriteExpectation{}
	}

	mmWrite.defaultExpectation.params = &WriterMockWriteParams{p}
	for _, e := range mmWrite.expectations {
		if minimock.Equal(e.params, mmWrite.defaultExpectation.params) {
			mmWrite.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmWrite.defaultExpectation.params)
		}
	}

	return mmWrite
}

// Inspect accepts an inspector function that has same arguments as the Writer.Write
func (mmWrite *mWriterMockWrite) Inspect(f func(p []byte)) *mWriterMockWrite {
	if mmWrite.mock.inspectFuncWrite != nil {
		mmWrite.mock.t.Fatalf("Inspect function is already set for WriterMock.Write")
	}

	mmWrite.mock.inspectFuncWrite = f

	return mmWrite
}

// Return sets up results that will be returned by Writer.Write
func (mmWrite *mWriterMockWrite) Return(n int, err error) *WriterMock {
	if mmWrite.mock.funcWrite != nil {
		mmWrite.mock.t.Fatalf("WriterMock.Write mock is already set by Set")
	}

	if mmWrite.defaultExpectation == nil {
		mmWrite.defaultExpectation = &WriterMockWriteExpectation{mock: mmWrite.mock}
	}
	mmWrite.defaultExpectation.results = &WriterMockWriteResults{n, err}
	return mmWrite.mock
}

//Set uses given function f to mock the Writer.Write method
func (mmWrite *mWriterMockWrite) Set(f func(p []byte) (n int, err error)) *WriterMock {
	if mmWrite.defaultExpectation != nil {
		mmWrite.mock.t.Fatalf("Default expectation is already set for the Writer.Write method")
	}

	if len(mmWrite.expectations) > 0 {
		mmWrite.mock.t.Fatalf("Some expectations are already set for the Writer.Write method")
	}

	mmWrite.mock.funcWrite = f
	return mmWrite.mock
}

// When sets expectation for the Writer.Write which will trigger the result defined by the following
// Then helper
func (mmWrite *mWriterMockWrite) When(p []byte) *WriterMockWriteExpectation {
	if mmWrite.mock.funcWrite != nil {
		mmWrite.mock.t.Fatalf("WriterMock.Write mock is already set by Set")
	}

	expectation := &WriterMockWriteExpectation{
		mock:   mmWrite.mock,
		params: &WriterMockWriteParams{p},
	}
	mmWrite.expectations = append(mmWrite.expectations, expectation)
	return expectation
}

// Then sets up Writer.Write return parameters for the expectation previously defined by the When method
func (e *WriterMockWriteExpectation) Then(n int, err error) *WriterMock {
	e.results = &WriterMockWriteResults{n, err}
	return e.mock
}

// Write implements io.Writer
func (mmWrite *WriterMock) Write(p []byte) (n int, err error) {
	mm_atomic.AddUint64(&mmWrite.beforeWriteCounter, 1)
	defer mm_atomic.AddUint64(&mmWrite.afterWriteCounter, 1)

	if mmWrite.inspectFuncWrite != nil {
		mmWrite.inspectFuncWrite(p)
	}

	mm_params := &WriterMockWriteParams{p}

	// Record call args
	mmWrite.WriteMock.mutex.Lock()
	mmWrite.WriteMock.callArgs = append(mmWrite.WriteMock.callArgs, mm_params)
	mmWrite.WriteMock.mutex.Unlock()

	for _, e := range mmWrite.WriteMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.n, e.results.err
		}
	}

	if mmWrite.WriteMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmWrite.WriteMock.defaultExpectation.Counter, 1)
		mm_want := mmWrite.WriteMock.defaultExpectation.params
		mm_got := WriterMockWriteParams{p}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmWrite.t.Errorf("WriterMock.Write got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmWrite.WriteMock.defaultExpectation.results
		if mm_results == nil {
			mmWrite.t.Fatal("No results are set for the WriterMock.Write")
		}
		return (*mm_results).n, (*mm_results).err
	}
	if mmWrite.funcWrite != nil {
		return mmWrite.funcWrite(p)
	}
	mmWrite.t.Fatalf("Unexpected call to WriterMock.Write. %v", p)
	return
}

// WriteAfterCounter returns a count of finished WriterMock.Write invocations
func (mmWrite *WriterMock) WriteAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmWrite.afterWriteCounter)
}

// WriteBeforeCounter returns a count of WriterMock.Write invocations
func (mmWrite *WriterMock) WriteBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmWrite.beforeWriteCounter)
}

// Calls returns a list of arguments used in each call to WriterMock.Write.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmWrite *mWriterMockWrite) Calls() []*WriterMockWriteParams {
	mmWrite.mutex.RLock()

	argCopy := make([]*WriterMockWriteParams, len(mmWrite.callArgs))
	copy(argCopy, mmWrite.callArgs)

	mmWrite.mutex.RUnlock()

	return argCopy
}

// MinimockWriteDone returns true if the count of the Write invocations corresponds
// the number of defined expectations
func (m *WriterMock) MinimockWriteDone() bool {
	for _, e := range m.WriteMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.WriteMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterWriteCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcWrite != nil && mm_atomic.LoadUint64(&m.afterWriteCounter) < 1 {
		return false
	}
	return true
}

// MinimockWriteInspect logs each unmet expectation
func (m *WriterMock) MinimockWriteInspect() {
	for _, e := range m.WriteMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to WriterMock.Write with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.WriteMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterWriteCounter) < 1 {
		if m.WriteMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to WriterMock.Write")
		} else {
			m.t.Errorf("Expected call to WriterMock.Write with params: %#v", *m.WriteMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcWrite != nil && mm_atomic.LoadUint64(&m.afterWriteCounter) < 1 {
		m.t.Error("Expected call to WriterMock.Write")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *WriterMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockWriteInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *WriterMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *WriterMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockWriteDone()
}
