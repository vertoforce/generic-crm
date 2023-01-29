package sqlcrm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetFieldType(t *testing.T) {
	tests := []struct {
		Value    interface{}
		Expected string
	}{
		{Value: "Hello", Expected: "TEXT(500)"},
		{Value: "10/3/2005", Expected: "Datetime NULL DEFAULT NULL"},
		{Value: "Feb 1st 2022", Expected: "Datetime NULL DEFAULT NULL"},
		{Value: "Aaron, Mary", Expected: "TEXT(500)"},
	}

	for _, test := range tests {
		actual := getFieldType(test.Value)
		require.Equal(t, test.Expected, actual)
	}
}
