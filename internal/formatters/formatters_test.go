package formatters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFormatter(t *testing.T) {
	t.Run("Different formats", func(t *testing.T) {
		f, err := NewFormatter("")
		assert.NoError(t, err)
		assert.NotNil(t, f)
		assert.Equal(t, &StylishFormatter{}, f)

		f, err = NewFormatter(StylishOutputFormat)
		assert.NoError(t, err)
		assert.NotNil(t, f)
		assert.Equal(t, &StylishFormatter{}, f)

		f, err = NewFormatter(PlainOutputFormat)
		assert.NoError(t, err)
		assert.NotNil(t, f)
		assert.Equal(t, &PlainFormatter{}, f)

		_, err = NewFormatter("Some format")
		assert.Error(t, err)
		assert.Equal(t, "format is not supported", err.Error())
	})
}
