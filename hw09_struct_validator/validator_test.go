package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type UserRole string

type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          Response{100, ""},
			expectedErr: errors.New("not in these mass of int , "),
		},
		{
			in: User{
				ID:     "1002340",
				Age:    67,
				Email:  "info@example.com",
				Role:   "admin",
				Phones: []string{"79061234567", "79012345678"},
			},
			expectedErr: errors.New("len is not equal 36, greater then 50, "),
		},
		{
			in: User{
				ID:     "1002345640",
				Age:    30,
				Email:  "info@examplecom",
				Role:   "admin",
				Phones: []string{"79061234567", "79012345678"},
				meta: []byte(`[
					{"Space": "YCbCr", "Point": {"Y": 255, "Cb": 0, "Cr": -10}},
					{"Space": "RGB",   "Point": {"R": 98, "G": 218, "B": 255}}	
					]`),
			},
			expectedErr: errors.New("len is not equal 36, string is not matched regexp expression , "),
		},
		{
			in: User{
				ID:     "10024567890340",
				Age:    37,
				Email:  "info@example.com",
				Role:   "manager",
				Phones: []string{"79061234567", "7901234678"},
			},
			expectedErr: errors.New("len is not equal 36, not in these mass of strings , len is not equal 11, "),
		},
		{
			in:          App{"v.10"},
			expectedErr: errors.New("len is not equal 5, "),
		},
		{
			in:          Response{Code: 505},
			expectedErr: errors.New("not in these mass of int , "),
		},
		{
			in:          Token{Header: []byte{4, 5}},
			expectedErr: errors.New(""),
		},
	}

	t.Run("not struct case", func(t *testing.T) {
		in := "not struct"
		err := Validate(in)
		require.EqualError(t, err, ErrNotStruct.Error())
	})

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)

			if tt.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			}
		})
	}
}
