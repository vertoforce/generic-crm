package googlesheet

import (
	"context"
	"fmt"
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	crm "github.com/vertoforce/generic-crm"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmot"
)

func BenchmarkGetItem(b *testing.B) {
	ctx := context.Background()

	opentracing.SetGlobalTracer(apmot.New())
	defer apm.DefaultTracer.Flush(nil)

	// Add some items to the sheet
	client, err := getTestingClient()
	require.NoError(b, err)

	// Disable sync, we are just testing local performance
	client.WaitToSynchronize = true

	for i := 0; i < 1000; i++ {
		client.CreateItem(ctx, &crm.DefaultItem{Fields: map[string]interface{}{
			"Name": fmt.Sprintf("Name %d", i),
			"Item": fmt.Sprintf("%d", i),
		}})
	}

	b.Run("FindItems", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			items := make(chan crm.Item)
			go func() {
				defer close(items)
				err := client.GetItems(ctx, items, map[string]interface{}{"Name": fmt.Sprintf("impossible to find me %d", i)})
				require.NoError(b, err)
			}()
			for range items {
			}
			assert.NoError(b, err)
		}
	})

}
