package entities

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/AndrXxX/goph-keeper/internal/client/formats"
)

func TestNoteItem(t *testing.T) {
	date := time.Now()
	tests := []struct {
		name            string
		entity          NoteItem
		wantFilterValue string
		wantTitle       string
		wantDescription string
	}{
		{
			name:            "Test with empty fields",
			entity:          NoteItem{StoredItem: StoredItem{UpdatedAt: date}},
			wantFilterValue: "",
			wantTitle:       "",
			wantDescription: date.Format(formats.FullDate),
		},
		{
			name:            "Test with filled fields",
			entity:          NoteItem{Text: "test very long text", StoredItem: StoredItem{UpdatedAt: date, Desc: "test_desc"}},
			wantFilterValue: "test very long texttest_desc",
			wantTitle:       "test very  ...",
			wantDescription: fmt.Sprintf("test_desc [%s]", date.Format(formats.FullDate)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantFilterValue, tt.entity.FilterValue())
			assert.Equal(t, tt.wantTitle, tt.entity.Title())
			assert.Equal(t, tt.wantDescription, tt.entity.Description())
		})
	}
}
