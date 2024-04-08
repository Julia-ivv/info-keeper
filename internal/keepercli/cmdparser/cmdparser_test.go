package cmdparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseUserCmd(t *testing.T) {
	tests := []struct {
		name     string
		c        string
		wantCmd  string
		wantArgs UserArgs
		wantErr  bool
	}{
		{
			name:     "reg",
			c:        "--reg -u=name",
			wantCmd:  CmdReg,
			wantArgs: UserArgs{AuthLogin: "name"},
			wantErr:  false,
		},
		{
			name:     "auth",
			c:        "--auth -u=name",
			wantCmd:  CmdAuth,
			wantArgs: UserArgs{AuthLogin: "name"},
			wantErr:  false,
		},
		{
			name:     "addCard",
			c:        "--ncard -p=prompt -n=12345 -e=12/24 -v=111 -m=comment",
			wantCmd:  CmdAddCard,
			wantArgs: UserArgs{Prompt: "prompt", CardNumber: "12345", CardDate: "12/24", CardCode: "111", Note: "comment"},
			wantErr:  false,
		},
		{
			name:     "addLogin",
			c:        "--npwd -p=prompt -l=login -m=comment",
			wantCmd:  CmdAddLogin,
			wantArgs: UserArgs{Prompt: "prompt", Login: "login", Note: "comment"},
			wantErr:  false,
		},
		{
			name:     "addText",
			c:        "--ntext -p=prompt -t=text -m=comment",
			wantCmd:  CmdAddText,
			wantArgs: UserArgs{Prompt: "prompt", Text: "text", Note: "comment"},
			wantErr:  false,
		},
		{
			name:     "addBinary",
			c:        "--nbyte -p=prompt -b=file -m=comment",
			wantCmd:  CmdAddBinary,
			wantArgs: UserArgs{Prompt: "prompt", Binary: "file", Note: "comment"},
			wantErr:  false,
		},
		{
			name:     "updCard",
			c:        "--ucard -p=prompt -n=12345 -e=12/24 -v=111 -m=comment",
			wantCmd:  CmdUpdCard,
			wantArgs: UserArgs{Prompt: "prompt", CardNumber: "12345", CardDate: "12/24", CardCode: "111", Note: "comment"},
			wantErr:  false,
		},
		{
			name:     "updLogin",
			c:        "--upwd -p=prompt -l=login -m=comment",
			wantCmd:  CmdUpdLogin,
			wantArgs: UserArgs{Prompt: "prompt", Login: "login", Note: "comment"},
			wantErr:  false,
		},
		{
			name:     "updText",
			c:        "--utext -p=prompt -t=text -m=comment",
			wantCmd:  CmdUpdText,
			wantArgs: UserArgs{Prompt: "prompt", Text: "text", Note: "comment"},
			wantErr:  false,
		},
		{
			name:     "updBinary",
			c:        "--ubyte -p=prompt -b=file -m=comment",
			wantCmd:  CmdUpdBinary,
			wantArgs: UserArgs{Prompt: "prompt", Binary: "file", Note: "comment"},
			wantErr:  false,
		},
		{
			name:     "getCard",
			c:        "--gcard -n=12345",
			wantCmd:  CmdGetCard,
			wantArgs: UserArgs{CardNumber: "12345"},
			wantErr:  false,
		},
		{
			name:     "getLogin",
			c:        "--gpwd -p=prompt -l=login",
			wantCmd:  CmdGetLogin,
			wantArgs: UserArgs{Prompt: "prompt", Login: "login"},
			wantErr:  false,
		},
		{
			name:     "getText",
			c:        "--gtext -p=prompt",
			wantCmd:  CmdGetText,
			wantArgs: UserArgs{Prompt: "prompt"},
			wantErr:  false,
		},
		{
			name:     "getBinary",
			c:        "--gbyte -p=prompt",
			wantCmd:  CmdGetBinary,
			wantArgs: UserArgs{Prompt: "prompt"},
			wantErr:  false,
		},
		{
			name:     "getCards",
			c:        "--gcards",
			wantCmd:  CmdGetCards,
			wantArgs: UserArgs{},
			wantErr:  false,
		},
		{
			name:     "getLogins",
			c:        "--gpwds",
			wantCmd:  CmdGetLogins,
			wantArgs: UserArgs{},
			wantErr:  false,
		},
		{
			name:     "getTexts",
			c:        "--gtexts",
			wantCmd:  CmdGetTexts,
			wantArgs: UserArgs{},
			wantErr:  false,
		},
		{
			name:     "getBinarys",
			c:        "--gbytes",
			wantCmd:  CmdGetBinarys,
			wantArgs: UserArgs{},
			wantErr:  false,
		},
		{
			name:     "forceAddCardServer",
			c:        "--fcard -n=12345",
			wantCmd:  CmdForceAddCardServer,
			wantArgs: UserArgs{CardNumber: "12345"},
			wantErr:  false,
		},
		{
			name:     "forceAddLoginServer",
			c:        "--fpwd -p=prompt -l=login",
			wantCmd:  CmdForceAddLoginServer,
			wantArgs: UserArgs{Prompt: "prompt", Login: "login"},
			wantErr:  false,
		},
		{
			name:     "forceAddTextServer",
			c:        "--ftext -p=prompt",
			wantCmd:  CmdForceAddTextServer,
			wantArgs: UserArgs{Prompt: "prompt"},
			wantErr:  false,
		},
		{
			name:     "forceAddBinaryServer",
			c:        "--fbyte -p=prompt",
			wantCmd:  CmdForceAddBinaryServer,
			wantArgs: UserArgs{Prompt: "prompt"},
			wantErr:  false,
		},
		{
			name:     "getCardServer",
			c:        "--scard -n=12345",
			wantCmd:  CmdGetCardServer,
			wantArgs: UserArgs{CardNumber: "12345"},
			wantErr:  false,
		},
		{
			name:     "getLoginServer",
			c:        "--spwd -p=prompt -l=login",
			wantCmd:  CmdGetLoginServer,
			wantArgs: UserArgs{Prompt: "prompt", Login: "login"},
			wantErr:  false,
		},
		{
			name:     "getTextServer",
			c:        "--stext -p=prompt",
			wantCmd:  CmdGetTextServer,
			wantArgs: UserArgs{Prompt: "prompt"},
			wantErr:  false,
		},
		{
			name:     "getBinaryServer",
			c:        "--sbyte -p=prompt",
			wantCmd:  CmdGetBinaryServer,
			wantArgs: UserArgs{Prompt: "prompt"},
			wantErr:  false,
		},
		{
			name:     "exit",
			c:        "-x",
			wantCmd:  CmdExit,
			wantArgs: UserArgs{},
			wantErr:  false,
		},
		{
			name:     "long exit",
			c:        "--exit",
			wantCmd:  CmdExit,
			wantArgs: UserArgs{},
			wantErr:  false,
		},
		{
			name:     "unclear",
			c:        "--clear",
			wantCmd:  "",
			wantArgs: UserArgs{},
			wantErr:  false,
		},
		{
			name:     "help",
			c:        "-h",
			wantCmd:  "",
			wantArgs: UserArgs{},
			wantErr:  false,
		},
		{
			name:     "err",
			c:        "err",
			wantCmd:  "",
			wantArgs: UserArgs{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, a, err := ParseUserCmd(tt.c)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.wantArgs, a)
				assert.Equal(t, tt.wantCmd, cmd)
			}
		})
	}
}
