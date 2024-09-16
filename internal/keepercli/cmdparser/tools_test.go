package cmdparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitCmd(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		s := splitCmd("abf dhj")
		assert.Equal(t, []string{"abf dhj"}, s)
	})
	t.Run("ok too", func(t *testing.T) {
		s := splitCmd("abf - dhj")
		assert.Equal(t, []string{"abf", "- dhj"}, s)
	})
	t.Run("ok else", func(t *testing.T) {
		s := splitCmd("--reg -n=name")
		assert.Equal(t, []string{"--reg", "-n=name"}, s)
	})
	t.Run("ok more", func(t *testing.T) {
		s := splitCmd("--reg -n='name'")
		assert.Equal(t, []string{"--reg", "-n='name'"}, s)
	})
}

func TestClearOpt(t *testing.T) {
	t.Run("ok test", func(t *testing.T) {
		o := Options{
			Reg:                  true,
			Auth:                 true,
			AddCard:              true,
			AddLogin:             true,
			AddText:              true,
			AddBinary:            true,
			UpdCard:              true,
			UpdLogin:             true,
			UpdText:              true,
			UpdBinary:            true,
			GetCard:              true,
			GetLogin:             true,
			GetText:              true,
			GetBinary:            true,
			GetCards:             true,
			GetLogins:            true,
			GetTexts:             true,
			GetBinarys:           true,
			ForceAddCardServer:   true,
			ForceAddLoginServer:  true,
			ForceAddTextServer:   true,
			ForceAddBinaryServer: true,
			GetCardServer:        true,
			GetLoginServer:       true,
			GetTextServer:        true,
			GetBinaryServer:      true,
			UserLogin:            "q",
			Prompt:               "q",
			Login:                "q",
			Note:                 "q",
			CardNumber:           "q",
			CardDate:             "q",
			CardCode:             "q",
			Text:                 "q",
			Binary:               "q",
			Exit:                 true,
		}
		err := clearOpt(&o)
		if assert.NoError(t, err) {
			assert.Equal(t, Options{}, o)
		}
	})
}
