package render

import (
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsentMessages(t *testing.T) {
	assert.Len(t, ConsentMessages, 14)
	for k := range ConsentMessages {
		assert.True(t, strings.HasPrefix(k, "consent."), k)
	}
	assert.Equal(t, "Allow", ConsentMessages["consent.allow"])
	assert.Equal(t, "Signed in as %s.", ConsentMessages["consent.signed-in-as"])

	keys := ConsentMessageKeys()
	assert.Len(t, keys, 14)
	assert.True(t, sort.StringsAreSorted(keys))
}
