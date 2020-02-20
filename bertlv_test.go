package bertlv

import (
	"bytes"
	"io"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestEncodeInt(t *testing.T) {
	tests := []struct {
		name  string
		in    int
		want1 []byte
	}{
		{
			name:  "one byte",
			in:    111,
			want1: []byte{111},
		},
		{
			name:  "two byte",
			in:    411,
			want1: []byte{0x01, 0x9b},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1 := EncodeInt(tt.in)
			assert.Equal(t, tt.want1, got1, "EncodeInt returned unexpected result")
		})
	}
}

func TestTagValue_Len(t *testing.T) {
	tests := []struct {
		name  string
		tv    TagValue
		want1 []byte
	}{
		{
			name:  "one byte len",
			tv:    TagValue{V: []byte{1, 2, 3}},
			want1: []byte{3},
		},
		{
			name:  "tree byte len",
			tv:    TagValue{V: make([]byte, 132768)},
			want1: []byte{0x83, 0x02, 0x06, 0xa0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1 := tt.tv.Len()
			assert.Equal(t, tt.want1, got1, "TagValue.Len returned unexpected result")
		})
	}
}

func TestTagValue_Tag(t *testing.T) {
	tests := []struct {
		name  string
		tv    TagValue
		want1 []byte
	}{
		{
			name:  "0xbf0c",
			tv:    TagValue{T: 0xbf0c},
			want1: []byte{0xbf, 0x0c},
		},
		{
			name:  "0x4f",
			tv:    TagValue{T: 0x4f},
			want1: []byte{0x4f},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1 := tt.tv.Tag()
			assert.Equal(t, tt.want1, got1, "TagValue.Tag returned unexpected result")
		})
	}
}

func TestReadTag(t *testing.T) {
	tests := []struct {
		name    string
		r       io.Reader
		want1   int
		want2   int
		wantErr bool
	}{
		{
			name:    "0xbf0c",
			r:       bytes.NewReader([]byte{0xbf, 0x0c}),
			want1:   0xbf0c,
			want2:   2,
			wantErr: false,
		},
		{
			name:    "0x4f",
			r:       bytes.NewReader([]byte{0x4f}),
			want1:   0x4f,
			want2:   1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, got2, err := ReadTag(tt.r)

			assert.Equal(t, tt.want1, got1, "ReadTag returned unexpected result")
			assert.Equal(t, tt.want2, got2, "ReadTag returned unexpected result")

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestReadLen(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name string
		args func(t minimock.Tester) args

		want1      int
		want2      int
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation
	}{
		{
			name: "one byte",
			args: func(t minimock.Tester) args {
				return args{r: bytes.NewReader([]byte{0x05})}
			},
			want1:   5,
			want2:   1,
			wantErr: false,
		},
		{
			name: "tree byte",
			args: func(t minimock.Tester) args {
				return args{r: bytes.NewReader([]byte{0x83, 0x02, 0x06, 0xa0})}
			},
			want1:   132768,
			want2:   4,
			wantErr: false,
		},
		{
			name: "indefinite length",
			args: func(t minimock.Tester) args {
				return args{r: bytes.NewReader([]byte{0x80, 0x02, 0x06, 0xa0})}
			},
			want1:   0,
			want2:   1,
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, errIndefiniteLength, err)
			},
		},
		{
			name: "invalid length",
			args: func(t minimock.Tester) args {
				return args{r: bytes.NewReader([]byte{0x85, 0x02, 0x06, 0xa0})}
			},
			want1:   0,
			want2:   1,
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, errInvalidLength, err)
			},
		},
		{
			name: "first read error",
			args: func(t minimock.Tester) args {
				return args{r: NewReaderMock(t).ReadMock.Return(2, io.EOF)}
			},
			want1:   0,
			want2:   2,
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, io.EOF, err)
			},
		},
		{
			name: "second read error",
			args: func(t minimock.Tester) args {
				var calls uint32
				return args{r: NewReaderMock(t).ReadMock.Set(func(p []byte) (int, error) {
					switch atomic.AddUint32(&calls, 1) {
					case 1:
						p[0] = 0x83
						return 1, nil
					case 2:
						assert.Len(t, p, 3)
						return 0, io.ErrUnexpectedEOF
					}
					t.Fatalf("too many calls to reader mock")
					return 0, nil
				})}
			},
			want1:   0,
			want2:   1,
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, io.ErrUnexpectedEOF, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Finish()

			tArgs := tt.args(mc)

			got1, got2, err := ReadLen(tArgs.r)

			assert.Equal(t, tt.want1, got1, "ReadLen returned unexpected result")

			assert.Equal(t, tt.want2, got2, "ReadLen returned unexpected result")

			if tt.wantErr {
				if assert.Error(t, err) && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestTagValue_WriteTo(t *testing.T) {
	type args struct {
		w io.Writer
	}
	tests := []struct {
		name string
		tv   TagValue

		args func(t minimock.Tester) args

		want1      int64
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation
	}{
		{
			name: "success",
			tv:   TagValue{T: 0x9f38, V: []byte("hello world")},
			args: func(t minimock.Tester) args {
				return args{
					w: NewWriterMock(t).WriteMock.Return(14, nil),
				}
			},
			want1:   14,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Wait(time.Second)

			tArgs := tt.args(mc)

			got1, err := tt.tv.WriteTo(tArgs.w)

			assert.Equal(t, tt.want1, got1, "TagValue.WriteTo returned unexpected result")

			if tt.wantErr {
				if assert.Error(t, err) && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestTagValue_ReadFrom(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		init    func(t minimock.Tester) TagValue
		inspect func(r TagValue, t *testing.T) //inspects TagValue after execution of ReadFrom

		args func(t minimock.Tester) args

		want1      int64
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation
	}{
		{
			name: "tag read error",
			init: func(t minimock.Tester) TagValue { return TagValue{} },
			args: func(t minimock.Tester) args {
				return args{
					r: NewReaderMock(t).ReadMock.Return(0, io.EOF),
				}
			},
			want1:   0,
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, io.EOF, err)
			},
		},
		{
			name: "len read error",
			init: func(t minimock.Tester) TagValue { return TagValue{} },
			args: func(t minimock.Tester) args {
				var calls uint32

				return args{
					r: NewReaderMock(t).ReadMock.Set(func(p []byte) (int, error) {
						switch atomic.AddUint32(&calls, 1) {
						case 1:
							p[0] = 0x4f
							return 1, nil
						case 2:
							return 0, io.ErrUnexpectedEOF
						}

						t.Fatalf("too many calls to reader")
						return 0, nil
					}),
				}
			},
			want1:   1,
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Contains(t, err.Error(), "failed to read length")
			},
		},
		{
			name: "success",
			init: func(t minimock.Tester) TagValue { return TagValue{} },
			args: func(t minimock.Tester) args {
				var calls uint32

				return args{
					r: NewReaderMock(t).ReadMock.Set(func(p []byte) (int, error) {
						switch atomic.AddUint32(&calls, 1) {
						case 1:
							p[0] = 0x4f
							return 1, nil
						case 2:
							p[0] = 11
							return 1, nil
						case 3:
							copy(p, []byte("hello world"))
							return 11, nil
						}

						t.Fatalf("too many calls to reader")
						return 0, nil
					}),
				}
			},
			want1:   13,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Wait(time.Second)

			tArgs := tt.args(mc)
			receiver := tt.init(mc)

			got1, err := receiver.ReadFrom(tArgs.r)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			assert.Equal(t, tt.want1, got1, "TagValue.ReadFrom returned unexpected result")

			if tt.wantErr {
				if assert.Error(t, err) && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		name string
		p    []byte

		want1      []TagValue
		wantErr    bool
		inspectErr func(err error, t *testing.T)
	}{
		{
			name:    "unexpected EOF",
			p:       []byte{0x4f, 0x05, 0, 1, 2},
			want1:   nil,
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, io.ErrUnexpectedEOF, err)
			},
		},
		{
			name:    "EOF",
			p:       []byte{0x4f, 0x03, 0, 1, 2},
			want1:   []TagValue{{T: 0x4f, V: []byte{0, 1, 2}}},
			wantErr: false,
		},
		{
			name: "multiple records",
			p:    []byte{0x4f, 0x03, 0, 1, 2, 0x4f, 0x02, 0, 1},
			want1: []TagValue{
				{T: 0x4f, V: []byte{0, 1, 2}},
				{T: 0x4f, V: []byte{0, 1}},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Wait(time.Second)

			got1, err := Decode(tt.p)

			assert.Equal(t, tt.want1, got1, "Decode returned unexpected result")

			if tt.wantErr {
				if assert.Error(t, err) && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestFind(t *testing.T) {
	//p, err := hex.DecodeString("a528500a566973612044656269745f2d047275656e9f380f9f66049f02069f37045f2a029f1a02bf0c00")
	//p, err := hex.DecodeString("500a566973612044656269745f2d047275656e9f380f9f66049f02069f37045f2a029f1a02bf0c00")
	//p, err := hex.DecodeString("500a56697361204465626974")
	//p, err := hex.DecodeString("5f2d047275656e9f380f9f66049f02069f37045f2a029f1a02bf0c00")
	//require.NoError(t, err)

	tests := []struct {
		name string
		tag  int
		p    []byte

		want1      *TagValue
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation
	}{
		{
			name:    "decode error",
			tag:     0x2f,
			p:       []byte{0x2f, 0x05, 0x00, 0x00},
			want1:   nil,
			wantErr: true,
		},
		{
			name:    "no tag found",
			tag:     0x2f,
			p:       []byte{0x4f, 0x02, 0x00, 0x00},
			want1:   nil,
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, ErrNotFound, err)
			},
		},
		{
			name:    "tag found",
			tag:     0x4f,
			p:       []byte{0x4f, 0x02, 0x00, 0x00},
			want1:   &TagValue{T: 0x4f, V: []byte{0, 0}},
			wantErr: false,
		},
		{
			name:    "tag found recursively",
			tag:     0x4f,
			p:       []byte{0xbf, 0x0c, 0x04, 0x4f, 0x02, 55, 44},
			want1:   &TagValue{T: 0x4f, V: []byte{55, 44}},
			wantErr: false,
		},
		{
			name:    "error during recursive search",
			tag:     0x4f,
			p:       []byte{0xbf, 0x0c, 0x04, 0x4f, 0x05, 55, 44},
			want1:   nil,
			wantErr: true,
		},
		/*
			{
				name:    "EOF",
				p:       p,
				want1:   nil,
				wantErr: false,
			},
		*/
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Wait(time.Second)

			got1, err := Find(tt.tag, tt.p)

			assert.Equal(t, tt.want1, got1, "Find returned unexpected result")

			if tt.wantErr {
				if assert.Error(t, err) && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestTagValue_MarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		tv   TagValue

		want1      []byte
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation
	}{
		{
			name:    "leaf",
			tv:      TagValue{T: 0x4f, V: []byte{0x12, 0x34}},
			want1:   []byte(`{"4f":"1234"}`),
			wantErr: false,
		},
		{
			name:    "recursive",
			tv:      TagValue{T: 0xbf0c, V: []byte{0x4f, 0x02, 0x12, 0x34}},
			want1:   []byte(`{"bf0c":[{"4f":"1234"}]}`),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Wait(time.Second)

			got1, err := tt.tv.MarshalJSON()

			assert.Equal(t, tt.want1, got1, "TagValue.MarshalJSON returned unexpected result")

			if tt.wantErr {
				if assert.Error(t, err) && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else {
				assert.NoError(t, err)
			}

		})
	}
}
