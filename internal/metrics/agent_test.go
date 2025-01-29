package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricsDatabase_Update(t *testing.T) {

	db := NewDB()

	t.Run("Update", func(t *testing.T) {
		db.poll()
		m := db.GetMetrics()
		assert.NotEqual(t, m, 0, "metrics should not be zero")
		db.poll()
		m2 := db.GetMetrics()

		assert.NotEqual(t, m, m2, "random values for consecutive polls should differ")
	})

}
