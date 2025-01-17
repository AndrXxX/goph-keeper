package entities

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/AndrXxX/goph-keeper/internal/client/formats"
)

func TestPasswordItem(t *testing.T) {
	date := time.Now()
	tests := []struct {
		name            string
		entity          PasswordItem
		wantFilterValue string
		wantTitle       string
		wantDescription string
	}{
		{
			name:            "Test with empty fields",
			entity:          PasswordItem{StoredItem: StoredItem{UpdatedAt: date}},
			wantFilterValue: "",
			wantTitle:       "",
			wantDescription: date.Format(formats.FullDate),
		},
		{
			name:            "Test with filled fields",
			entity:          PasswordItem{Login: "test_login", StoredItem: StoredItem{UpdatedAt: date, Desc: "test_desc"}},
			wantFilterValue: "test_logintest_desc",
			wantTitle:       "test_login",
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
