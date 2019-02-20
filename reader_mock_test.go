package bertlv

// DO NOT EDIT!
// The code below was generated with http://github.com/gojuno/minimock (2.1.2)

//go:generate minimock -i io.Reader -o ./reader_mock_test.go

import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
)

// ReaderMock implements io.Reader
type ReaderMock struct {
	t minimock.Tester

	funcRead          func(p []byte) (n int, err error)
	afterReadCounter  uint64
	beforeReadCounter uint64
	ReadMock          mReaderMockRead
}

// NewReaderMock returns a mock for io.Reader
func NewReaderMock(t minimock.Tester) *ReaderMock {
	m := &ReaderMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}
	m.ReadMock = mReaderMockRead{mock: m}

	return m
}

type mReaderMockRead struct {
	mock               *ReaderMock
	defaultExpectation *ReaderMockReadExpectation
	expectations       []*ReaderMockReadExpectation
}

// ReaderMockReadExpectation specifies expectation struct of the Reader.Read
type ReaderMockReadExpectation struct {
	mock    *ReaderMock
	params  *ReaderMockReadParams
	results *ReaderMockReadResults
	Counter uint64
}

// ReaderMockReadParams contains parameters of the Reader.Read
type ReaderMockReadParams struct {
	p []byte
}

// ReaderMockReadResults contains results of the Reader.Read
type ReaderMockReadResults struct {
	n   int
	err error
}

// Expect sets up expected params for Reader.Read
func (m *mReaderMockRead) Expect(p []byte) *mReaderMockRead {
	if m.mock.funcRead != nil {
		m.mock.t.Fatalf("ReaderMock.Read mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ReaderMockReadExpectation{}
	}

	m.defaultExpectation.params = &ReaderMockReadParams{p}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by Reader.Read
func (m *mReaderMockRead) Return(n int, err error) *ReaderMock {
	if m.mock.funcRead != nil {
		m.mock.t.Fatalf("ReaderMock.Read mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ReaderMockReadExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &ReaderMockReadResults{n, err}
	return m.mock
}

//Set uses given function f to mock the Reader.Read method
func (m *mReaderMockRead) Set(f func(p []byte) (n int, err error)) *ReaderMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the Reader.Read method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the Reader.Read method")
	}

	m.mock.funcRead = f
	return m.mock
}

// When sets expectation for the Reader.Read which will trigger the result defined by the following
// Then helper
func (m *mReaderMockRead) When(p []byte) *ReaderMockReadExpectation {
	if m.mock.funcRead != nil {
		m.mock.t.Fatalf("ReaderMock.Read mock is already set by Set")
	}

	expectation := &ReaderMockReadExpectation{
		mock:   m.mock,
		params: &ReaderMockReadParams{p},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up Reader.Read return parameters for the expectation previously defined by the When method
func (e *ReaderMockReadExpectation) Then(n int, err error) *ReaderMock {
	e.results = &ReaderMockReadResults{n, err}
	return e.mock
}

// Read implements io.Reader
func (m *ReaderMock) Read(p []byte) (n int, err error) {
	atomic.AddUint64(&m.beforeReadCounter, 1)
	defer atomic.AddUint64(&m.afterReadCounter, 1)

	for _, e := range m.ReadMock.expectations {
		if minimock.Equal(*e.params, ReaderMockReadParams{p}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.n, e.results.err
		}
	}

	if m.ReadMock.defaultExpectation != nil {
		atomic.AddUint64(&m.ReadMock.defaultExpectation.Counter, 1)
		want := m.ReadMock.defaultExpectation.params
		got := ReaderMockReadParams{p}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("ReaderMock.Read got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.ReadMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the ReaderMock.Read")
		}
		return (*results).n, (*results).err
	}
	if m.funcRead != nil {
		return m.funcRead(p)
	}
	m.t.Fatalf("Unexpected call to ReaderMock.Read. %v", p)
	return
}

// ReadAfterCounter returns a count of finished ReaderMock.Read invocations
func (m *ReaderMock) ReadAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterReadCounter)
}

// ReadBeforeCounter returns a count of ReaderMock.Read invocations
func (m *ReaderMock) ReadBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeReadCounter)
}

// MinimockReadDone returns true if the count of the Read invocations corresponds
// the number of defined expectations
func (m *ReaderMock) MinimockReadDone() bool {
	for _, e := range m.ReadMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ReadMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterReadCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcRead != nil && atomic.LoadUint64(&m.afterReadCounter) < 1 {
		return false
	}
	return true
}

// MinimockReadInspect logs each unmet expectation
func (m *ReaderMock) MinimockReadInspect() {
	for _, e := range m.ReadMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to ReaderMock.Read with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ReadMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterReadCounter) < 1 {
		m.t.Errorf("Expected call to ReaderMock.Read with params: %#v", *m.ReadMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcRead != nil && atomic.LoadUint64(&m.afterReadCounter) < 1 {
		m.t.Error("Expected call to ReaderMock.Read")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *ReaderMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockReadInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *ReaderMock) MinimockWait(timeout time.Duration) {
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

func (m *ReaderMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockReadDone()
}
