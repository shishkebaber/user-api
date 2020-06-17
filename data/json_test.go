package data

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToJSON(t *testing.T) {
	user := []*User{{1, "Roman", "Gavrilov", "shishkebaber", "shishkebaber@gmail.com", "test", "Russia"}}

	b := bytes.NewBufferString("")
	err := ToJson(user, b)
	assert.NoError(t, err)
}
