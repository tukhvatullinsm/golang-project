package metrics

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExportMetrics(t *testing.T) {
	testObj := &MyMetrics{}
	tests := []struct {
		name string
	}{
		{name: "Empty return"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.NotNil(t, testObj.ExportMetrics())
		})
	}
}
