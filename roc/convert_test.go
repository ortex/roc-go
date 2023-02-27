package roc

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func Test_go2cBool(t *testing.T) {
	tests := []struct {
		arg  bool
		want uint
	}{
		{true, 1},
		{false, 0},
	}
	for _, tt := range tests {
		t.Run(strconv.FormatBool(tt.arg), func(t *testing.T) {
			assert.Equal(t, tt.want, uint(go2cBool(tt.arg)))
		})
	}
}

func Test_go2cStr_c2goStr(t *testing.T) {
	tests := []struct {
		str string
	}{
		{"str"},
		{""},
	}
	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			assert.Equal(t, tt.str, c2goStr(go2cStr(tt.str)))
		})
	}
}
