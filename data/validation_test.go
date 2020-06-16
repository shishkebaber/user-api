package data

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidation(t *testing.T) {
	cases := []struct {
		caseName     string
		user         User
		errorsLength int
	}{
		{"Check Name is exist", User{1, "", "Gavrilov", "shishkebaber", "shishkebaber@gmail.com", "test", "Russia"}, 1},
		{"Check Last Name is exist", User{1, "Roman", "", "shishkebaber", "shishkebaber@gmail.com", "test", "Russia"}, 1},
		{"Check Nickname is exist", User{1, "Roman", "Gavrilov", "", "shishkebaber@gmail.com", "test", "Russia"}, 1},
		{"Check Email is exist", User{1, "Roman", "Gavrilov", "shishkebaber", "", "test", "Russia"}, 1},
		{"Check Email is has correct format", User{1, "Roman", "Gavrilov", "shishkebaber", "shishkebaber", "test", "Russia"}, 1},
		{"Check Password is exist", User{1, "Roman", "Gavrilov", "shishkebaber", "shishkebaber@gmail.com", "", "Russia"}, 1},
		{"Check Country is exist", User{1, "Roman", "Gavrilov", "shishkebaber", "shishkebaber@gmail.com", "test", ""}, 1},
	}

	v := NewValidation()
	for _, tCase := range cases {
		err := v.Validate(tCase.user)
		if len(err) != 1 {
			t.Fatalf("%v is failed", tCase.caseName)
		}
	}
}

func TestToJSON(t *testing.T) {
	user := []*User{{1, "Roman", "Gavrilov", "shishkebaber", "shishkebaber@gmail.com", "test", "Russia"}}

	b := bytes.NewBufferString("")
	err := ToJson(user, b)
	assert.NoError(t, err)
}
