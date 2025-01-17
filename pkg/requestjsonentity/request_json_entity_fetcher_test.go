package requestjsonentity

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testStruct struct {
	Name string `json:"name" valid:"minstringlength(3)"`
	Val  int    `json:"val"`
}

func TestFetcher_Fetch(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    *testStruct
		wantErr bool
	}{
		{
			name:    "Test with decode err",
			data:    "../",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Test with validate err",
			data:    "{\"name\":\"t1\",\"val\":10}",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Test OK",
			data:    "{\"name\":\"test1\",\"val\":10}",
			want:    &testStruct{Name: "test1", Val: 10},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Fetcher[testStruct]{}
			got, err := f.Fetch(strings.NewReader(tt.data))
			require.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFetcher_FetchSlice(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    []testStruct
		wantErr bool
	}{
		{
			name:    "Test with decode err",
			data:    "../",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Test with validate err",
			data:    "[{\"name\":\"t1\",\"val\":10}]",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Test OK",
			data:    "[{\"name\":\"test1\",\"val\":10}]",
			want:    []testStruct{{Name: "test1", Val: 10}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Fetcher[testStruct]{}
			got, err := f.FetchSlice(strings.NewReader(tt.data))
			require.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
