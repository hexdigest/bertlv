package bertlv

// DO NOT EDIT!
// The code below was generated with http://github.com/gojuno/minimock (2.1.2)

//go:generate minimock -i io.Writer -o ./writer_mock_test.go

import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
)

// WriterMock implements io.Writer
type WriterMock struct {
	t minimock.Tester

	funcWrite          func(p []byte) (n int, err error)
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

	return m
}

type mWriterMockWrite struct {
	mock               *WriterMock
	defaultExpectation *WriterMockWriteExpectation
	expectations       []*WriterMockWriteExpectation
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
func (m *mWriterMockWrite) Expect(p []byte) *mWriterMockWrite {
	if m.mock.funcWrite != nil {
		m.mock.t.Fatalf("WriterMock.Write mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &WriterMockWriteExpectation{}
	}

	m.defaultExpectation.params = &WriterMockWriteParams{p}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by Writer.Write
func (m *mWriterMockWrite) Return(n int, err error) *WriterMock {
	if m.mock.funcWrite != nil {
		m.mock.t.Fatalf("WriterMock.Write mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &WriterMockWriteExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &WriterMockWriteResults{n, err}
	return m.mock
}

//Set uses given function f to mock the Writer.Write method
func (m *mWriterMockWrite) Set(f func(p []byte) (n int, err error)) *WriterMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the Writer.Write method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the Writer.Write method")
	}

	m.mock.funcWrite = f
	return m.mock
}

// When sets expectation for the Writer.Write which will trigger the result defined by the following
// Then helper
func (m *mWriterMockWrite) When(p []byte) *WriterMockWriteExpectation {
	if m.mock.funcWrite != nil {
		m.mock.t.Fatalf("WriterMock.Write mock is already set by Set")
	}

	expectation := &WriterMockWriteExpectation{
		mock:   m.mock,
		params: &WriterMockWriteParams{p},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up Writer.Write return parameters for the expectation previously defined by the When method
func (e *WriterMockWriteExpectation) Then(n int, err error) *WriterMock {
	e.results = &WriterMockWriteResults{n, err}
	return e.mock
}

// Write implements io.Writer
func (m *WriterMock) Write(p []byte) (n int, err error) {
	atomic.AddUint64(&m.beforeWriteCounter, 1)
	defer atomic.AddUint64(&m.afterWriteCounter, 1)

	for _, e := range m.WriteMock.expectations {
		if minimock.Equal(*e.params, WriterMockWriteParams{p}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.n, e.results.err
		}
	}

	if m.WriteMock.defaultExpectation != nil {
		atomic.AddUint64(&m.WriteMock.defaultExpectation.Counter, 1)
		want := m.WriteMock.defaultExpectation.params
		got := WriterMockWriteParams{p}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("WriterMock.Write got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.WriteMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the WriterMock.Write")
		}
		return (*results).n, (*results).err
	}
	if m.funcWrite != nil {
		return m.funcWrite(p)
	}
	m.t.Fatalf("Unexpected call to WriterMock.Write. %v", p)
	return
}

// WriteAfterCounter returns a count of finished WriterMock.Write invocations
func (m *WriterMock) WriteAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterWriteCounter)
}

// WriteBeforeCounter returns a count of WriterMock.Write invocations
func (m *WriterMock) WriteBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeWriteCounter)
}

// MinimockWriteDone returns true if the count of the Write invocations corresponds
// the number of defined expectations
func (m *WriterMock) MinimockWriteDone() bool {
	for _, e := range m.WriteMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.WriteMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterWriteCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcWrite != nil && atomic.LoadUint64(&m.afterWriteCounter) < 1 {
		return false
	}
	return true
}

// MinimockWriteInspect logs each unmet expectation
func (m *WriterMock) MinimockWriteInspect() {
	for _, e := range m.WriteMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to WriterMock.Write with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.WriteMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterWriteCounter) < 1 {
		m.t.Errorf("Expected call to WriterMock.Write with params: %#v", *m.WriteMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcWrite != nil && atomic.LoadUint64(&m.afterWriteCounter) < 1 {
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
func (m *WriterMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-time.After(10 * time.Millisecond):
		}
	}
}

func (m *WriterMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockWriteDone()
}
