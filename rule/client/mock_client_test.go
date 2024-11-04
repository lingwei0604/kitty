package client

import (
	"testing"

	"github.com/lingwei0604/kitty/rule/dto"
	"github.com/stretchr/testify/assert"
)

func TestMockEngine(t *testing.T) {
	// Arrange
	rule := NewMockRule(func(pl *dto.Payload) dto.Data {
		return dto.Data{"foo": "bar"}
	})
	engine := NewMockEngine(map[string]Tenanter{
		"conf-prod": rule,
	})

	// Act
	c, err := engine.Of("conf-prod").Payload(nil)

	// asserts
	assert.NoError(t, err)
	assert.Equal(t, "bar", c.Get("foo"))
}
