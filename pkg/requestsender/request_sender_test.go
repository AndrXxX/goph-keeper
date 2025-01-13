package requestsender

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/AndrXxX/goph-keeper/pkg/hashgenerator"
	"github.com/AndrXxX/goph-keeper/pkg/requestsender/dto"
)

type closableReadableBodyMock struct {
	mock.Mock
	io.Reader
}

func (m *closableReadableBodyMock) Close() error {
	return nil
}

func (m *closableReadableBodyMock) Read(_ []byte) (n int, err error) {
	return 0, nil
}

type mockClient struct {
	mock.Mock
}

func (m *mockClient) Do(r *http.Request) (*http.Response, error) {
	args := m.Called(r)
	resp, _ := args.Get(0).(*http.Response)
	return resp, args.Error(1)
}

type dataCompressorMock struct {
	mock.Mock
}

func (m *dataCompressorMock) Compress(in []byte) (io.Reader, error) {
	args := m.Called(in)
	resp, _ := args.Get(0).(io.Reader)
	return resp, args.Error(1)
}

type readerMock struct {
	mock.Mock
}

func (r *readerMock) Read(in []byte) (n int, err error) {
	args := r.Called(in)
	return args.Int(0), args.Error(1)
}

func TestRequestSender_Post(t *testing.T) {
	comp := func() *dataCompressorMock {
		c := dataCompressorMock{}
		c.On("Compress", mock.Anything).Return(&bytes.Buffer{}, nil)
		return &c
	}
	type fields struct {
		c    client
		comp dataCompressor
	}
	tests := []struct {
		name    string
		fields  fields
		url     string
		data    []byte
		wantErr bool
	}{
		{
			name: "Positive test #1",
			fields: fields{
				c: func() *mockClient {
					c := mockClient{}
					c.On("Do", mock.Anything).Return(nil, nil)
					return &c
				}(),
				comp: comp(),
			},
			wantErr: false,
		},
		{
			name: "Positive test #2 with body",
			fields: fields{
				c: func() *mockClient {
					c := mockClient{}
					c.On("Do", mock.Anything).Return(&http.Response{Header: http.Header{}, Body: &closableReadableBodyMock{}}, nil)
					return &c
				}(),
				comp: comp(),
			},
			wantErr: false,
		},
		{
			name: "Positive test #3 with data",
			fields: fields{
				c: func() *mockClient {
					c := mockClient{}
					c.On("Do", mock.Anything).Return(nil, nil)
					return &c
				}(),
				comp: comp(),
			},
			data:    []byte("test"),
			wantErr: false,
		},
		{
			name: "Error on create request",
			fields: fields{
				c: func() *mockClient {
					c := mockClient{}
					c.On("Do", mock.Anything).Return(nil, nil)
					return &c
				}(),
				comp: comp(),
			},
			url:     string(rune(0x1B)),
			wantErr: true,
		},
		{
			name: "Error on do request",
			fields: fields{
				c: func() *mockClient {
					c := mockClient{}
					c.On("Do", mock.Anything).Return(nil, errors.New("error from web server"))
					return &c
				}(),
				comp: comp(),
			},
			wantErr: true,
		},
		{
			name: "Error on compress data",
			fields: fields{
				c: func() *mockClient {
					c := mockClient{}
					c.On("Do", mock.Anything).Return(nil, nil)
					return &c
				}(),
				comp: func() *dataCompressorMock {
					c := dataCompressorMock{}
					c.On("Compress", mock.Anything).Return(nil, errors.New("error"))
					return &c
				}(),
			},
			wantErr: true,
		},
		{
			name: "Error on read compressed data",
			fields: fields{
				c: func() *mockClient {
					c := mockClient{}
					c.On("Do", mock.Anything).Return(nil, nil)
					return &c
				}(),
				comp: func() *dataCompressorMock {
					r := readerMock{}
					r.On("Read", mock.Anything).Return(0, errors.New("error on read"))
					c := dataCompressorMock{}
					c.On("Compress", mock.Anything).Return(&r, nil)
					return &c
				}(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.fields.c, WithGzip(tt.fields.comp), WithSHA256(hashgenerator.Factory().SHA256("test")))
			_, err := s.Post(tt.url, "", bytes.NewBuffer(tt.data))
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestRequestSender_Get(t *testing.T) {
	tests := []struct {
		name    string
		c       client
		opts    []Option
		url     string
		wantErr bool
	}{
		{
			name: "Test with error on run option",
			opts: []Option{
				func(dto *dto.ParamsDto) error {
					return errors.New("error")
				},
			},
			wantErr: true,
		},
		{
			name:    "Error on create request",
			url:     string(rune(0x1B)),
			wantErr: true,
		},
		{
			name: "Positive test #1",
			c: func() *mockClient {
				c := mockClient{}
				c.On("Do", mock.Anything).Return(nil, nil)
				return &c
			}(),
			wantErr: false,
		},
		{
			name: "Error on do request",
			c: func() *mockClient {
				c := mockClient{}
				c.On("Do", mock.Anything).Return(nil, errors.New("error from web server"))
				return &c
			}(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.c, tt.opts...)
			_, err := s.Get(tt.url, "")
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
