package ad

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAd(t *testing.T) {
	t.Run("Test Ad struct", func(t *testing.T) {
		content := "Test ad content"
		metadata := map[string]string{"key": "value"}
		ad := Ad{
			Content:  content,
			Metadata: metadata,
		}

		assert.Equal(t, content, ad.Content)
		assert.Equal(t, metadata, ad.Metadata)
	})

	t.Run("Test SubscriptionHandler Ad storage", func(t *testing.T) {

		ad1 := Ad{Content: "Ad 1"}
		ad2 := Ad{Content: "Ad 2"}

		assert.NotEqual(t, ad1, ad2, nil)
		/// wip
	})
}

func toJSON(t *testing.T, ad Ad) []byte {
	data, err := json.Marshal(ad)
	require.NoError(t, err)
	return data
}
