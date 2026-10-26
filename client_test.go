package backlog

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestNewClient(t *testing.T) {
	baseURL := "https://example.com"
	apiKey := "apiKey"

	t.Run("initialization", func(t *testing.T) {
		t.Parallel()

		c, err := NewClient(baseURL, apiKey)
		require.NoError(t, err)
		assert.NotNil(t, c)

		// Reflection-based safety check
		clientType := reflect.TypeOf(*c)
		clientValue := reflect.ValueOf(*c)
		for i := 0; i < clientType.NumField(); i++ {
			field := clientType.Field(i)
			value := clientValue.Field(i)
			if field.Type.Kind() == reflect.Ptr {
				assert.Falsef(t, value.IsNil(), "%s should not be nil", field.Name)
			}
		}

		assert.NotNil(t, c.httpClient)
		assert.Equal(t, client.Credential{Type: client.CredentialAPIKey, Value: apiKey}, c.httpClient.Credential)
	})

	t.Run("with-doer", func(t *testing.T) {
		mockDoer := &mock.Doer{T: t,
			DoFunc: func(_ *http.Request) (*http.Response, error) { return nil, nil },
		}
		c, err := NewClient(baseURL, apiKey, WithDoer(mockDoer))
		require.NoError(t, err)
		assert.NotNil(t, c)

		assert.Same(t, c.httpClient.Doer, mockDoer)
	})

	t.Run("error-client.NewClient", func(t *testing.T) {
		c, err := NewClient("", "")
		require.Error(t, err)
		assert.IsType(t, &InternalClientError{}, err)
		assert.Nil(t, c)
	})
}

func TestNewOAuthClient(t *testing.T) {
	baseURL := "https://example.com"
	accessToken := "accessToken"

	t.Run("initialization", func(t *testing.T) {
		t.Parallel()

		c, err := NewOAuthClient(baseURL, accessToken)
		require.NoError(t, err)
		assert.NotNil(t, c)

		// Reflection-based safety check
		clientType := reflect.TypeOf(*c)
		clientValue := reflect.ValueOf(*c)
		for i := 0; i < clientType.NumField(); i++ {
			field := clientType.Field(i)
			value := clientValue.Field(i)
			if field.Type.Kind() == reflect.Ptr {
				assert.Falsef(t, value.IsNil(), "%s should not be nil", field.Name)
			}
		}

		assert.NotNil(t, c.httpClient)
		assert.Equal(t, client.Credential{Type: client.CredentialAccessToken, Value: accessToken}, c.httpClient.Credential)
	})

	t.Run("with-doer", func(t *testing.T) {
		mockDoer := &mock.Doer{T: t,
			DoFunc: func(_ *http.Request) (*http.Response, error) { return nil, nil },
		}
		c, err := NewOAuthClient(baseURL, accessToken, WithDoer(mockDoer))
		require.NoError(t, err)
		assert.NotNil(t, c)

		assert.Same(t, c.httpClient.Doer, mockDoer)
	})

	t.Run("error-client.NewOAuthClient", func(t *testing.T) {
		c, err := NewOAuthClient("", "")
		require.Error(t, err)
		assert.IsType(t, &InternalClientError{}, err)
		assert.Nil(t, c)
	})
}
