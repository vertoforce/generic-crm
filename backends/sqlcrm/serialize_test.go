package sqlcrm

import (
	"context"
	"testing"
	"time"

	"github.com/segmentio/agecache"
	"github.com/stretchr/testify/require"
)

func TestSerialize(t *testing.T) {
	test := map[string]interface{}{
		"String": "I am a string",
		// Currently not supported
		// "Bytes":  []byte{0, 1, 2},
		"Array": []interface{}{"One", "two", "three"},
		"Struct": struct {
			One string
			Two []byte
		}{One: "One", Two: []byte{0, 1, 2}},
		"Number":           int64(5),
		"float":            float64(5.3),
		"time":             time.Now(),
		"stringQuoted":     "\"test\"",
		"stringInByteForm": []byte("test"),
		"empty":            "",
	}

	// Try to serialize and deserialize and make sure the result is the same as the original
	c := &Client{
		columnsCache: agecache.New(agecache.Config[string, map[string]string]{
			MaxAge:   time.Second * 60,
			Capacity: 1,
		}),
	}
	c.columnsCache.Set("columns", map[string]string{})
	serializedFields, err := c.serializeFields(context.Background(), test)
	require.NoError(t, err)
	processed := deserializeFields(serializedFields)
	test["stringInByteForm"] = interface{}(string("test"))
	require.Equal(t, test, processed)
}
