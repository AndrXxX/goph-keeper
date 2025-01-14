package filestorage

import (
	"io"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	id1 = uuid.New()
)

const storagePath = "./"

type testHg struct {
	str string
}

func (g *testHg) Generate(_ []byte) string {
	return g.str
}

func Test_storage_Delete(t *testing.T) {
	tests := []struct {
		name         string
		hg           hashGenerator
		id           uuid.UUID
		beforeDelete func()
		wantErr      bool
	}{
		{
			name:         "Test with error",
			hg:           &testHg{str: "test.log"},
			id:           id1,
			beforeDelete: func() {},
			wantErr:      true,
		},
		{
			name: "Test ok",
			hg:   &testHg{str: "test.log"},
			id:   id1,
			beforeDelete: func() {
				f, _ := os.OpenFile(path.Join(storagePath, "test.log"), os.O_RDONLY|os.O_CREATE, os.ModePerm)
				_ = f.Close()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := New(storagePath, tt.hg)
			tt.beforeDelete()
			err := s.Delete(tt.id)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func Test_storage_FileId(t *testing.T) {
	tests := []struct {
		name string
		hg   hashGenerator
		id   uuid.UUID
		want string
	}{
		{
			name: "Test with HG test",
			hg:   &testHg{str: "test.log"},
			id:   id1,
			want: "test.log",
		},
		{
			name: "Test with HG from id1",
			hg:   &testHg{str: id1.String()},
			id:   id1,
			want: id1.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := New(storagePath, tt.hg)
			assert.Equal(t, tt.want, s.FileId(tt.id))
		})
	}
}

func Test_storage_Get(t *testing.T) {
	tests := []struct {
		name      string
		hg        hashGenerator
		beforeGet func()
		want      string
		wantErr   bool
	}{
		{
			name: "Test OK with content test",
			hg:   &testHg{str: "test.log"},
			beforeGet: func() {
				f, _ := os.Create(path.Join(storagePath, "test.log"))
				_, _ = f.WriteString("test")
				_ = f.Close()
			},
			want:    "test",
			wantErr: false,
		},
		{
			name:      "Test with error",
			hg:        &testHg{str: "test.log"},
			beforeGet: func() {},
			want:      "",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := New(storagePath, tt.hg)
			tt.beforeGet()
			f, err := s.Get(id1)
			require.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				data, err := io.ReadAll(f)
				assert.NoError(t, err)
				assert.Equal(t, tt.want, string(data))
			}
			_ = os.Remove(path.Join(s.path, tt.hg.Generate(nil)))
		})
	}
}

func Test_storage_IsExist(t *testing.T) {
	tests := []struct {
		name        string
		hg          hashGenerator
		beforeCheck func()
		want        bool
	}{
		{
			name: "Test with exist",
			hg:   &testHg{str: "test.log"},
			beforeCheck: func() {
				f, _ := os.Create(path.Join(storagePath, "test.log"))
				_ = f.Close()
			},
			want: true,
		},
		{
			name:        "Test with not exist",
			hg:          &testHg{str: "test.log"},
			beforeCheck: func() {},
			want:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := New(storagePath, tt.hg)
			tt.beforeCheck()
			assert.Equal(t, tt.want, s.IsExist(id1))
			_ = os.Remove(path.Join(s.path, tt.hg.Generate(nil)))
		})
	}
}

func Test_storage_Store(t *testing.T) {
	tests := []struct {
		name        string
		hg          hashGenerator
		beforeStore func()
		data        string
		want        string
		wantErr     bool
	}{
		{
			name:        "Test with OK store data",
			hg:          &testHg{str: "test.log"},
			beforeStore: func() {},
			data:        "test",
			want:        "test",
			wantErr:     false,
		},
		{
			name: "Test with error on write file",
			hg:   &testHg{str: "test.log"},
			beforeStore: func() {
				_, _ = os.OpenFile(path.Join(storagePath, "test.log"), os.O_CREATE, 0111)
			},
			data:    "test",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := New(storagePath, tt.hg)
			tt.beforeStore()
			err := s.Store(strings.NewReader(tt.data), id1)
			require.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				f, _ := os.OpenFile(path.Join(s.path, tt.hg.Generate(nil)), os.O_RDONLY, os.ModePerm)
				data, err := io.ReadAll(f)
				assert.NoError(t, err)
				assert.Equal(t, tt.want, string(data))
			}
			_ = os.Remove(path.Join(s.path, tt.hg.Generate(nil)))
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		hg      hashGenerator
		before  func()
		after   func()
		wantNil bool
		wantErr bool
	}{
		{
			name: "Test with error on create storage",
			path: path.Join(storagePath, "test.log"),
			hg:   &testHg{str: "test.log"},
			before: func() {
				f, _ := os.Create(path.Join(storagePath, "test.log"))
				_ = f.Close()
			},
			after: func() {
				_ = os.Remove(path.Join(storagePath, "test.log"))
			},
			wantNil: true,
			wantErr: true,
		},
		{
			name:   "Test with OK",
			path:   path.Join(storagePath, "testX"),
			hg:     &testHg{},
			before: func() {},
			after: func() {
				_ = os.Remove(path.Join(storagePath, "testX"))
			},
			wantNil: false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()
			got, err := New(tt.path, tt.hg)
			tt.after()
			require.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantNil, got == nil)
		})
	}
}
