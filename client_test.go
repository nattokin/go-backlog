package backlog

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient_initialization(t *testing.T) {
	t.Parallel()

	baseURL := "https://example.com"
	token := "token123"

	c, err := NewClient(baseURL, token)
	require.NoError(t, err)
	assert.NotNil(t, c)

	assert.NotNil(t, c.core)
	assert.Equal(t, token, c.core.Token)

	// Reflection-based safety check
	clientType := reflect.TypeOf(*c)
	clientValue := reflect.ValueOf(*c)
	for i := 0; i < clientType.NumField(); i++ {
		field := clientType.Field(i)
		value := clientValue.Field(i)
		if field.Type.Kind() == reflect.Ptr && field.Name != "httpClient" {
			assert.Falsef(t, value.IsNil(), "%s should not be nil", field.Name)
		}
	}
}
