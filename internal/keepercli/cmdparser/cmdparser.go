package cmdparser

import (
	"errors"

	flags "github.com/jessevdk/go-flags"
)

type UserCommandName = string

const (
	CmdReg  UserCommandName = "reg"
	CmdAuth UserCommandName = "auth"

	CmdAddCard   UserCommandName = "addCard"
	CmdAddLogin  UserCommandName = "addLogin"
	CmdAddText   UserCommandName = "addText"
	CmdAddBinary UserCommandName = "addBinary"

	CmdUpdCard   UserCommandName = "updCard"
	CmdUpdLogin  UserCommandName = "updLogin"
	CmdUpdText   UserCommandName = "updText"
	CmdUpdBinary UserCommandName = "updBinary"

	CmdGetCard   UserCommandName = "getCard"
	CmdGetLogin  UserCommandName = "getLogin"
	CmdGetText   UserCommandName = "getText"
	CmdGetBinary UserCommandName = "getBinary"

	CmdGetCards   UserCommandName = "getCards"
	CmdGetLogins  UserCommandName = "getLogins"
	CmdGetTexts   UserCommandName = "getTexts"
	CmdGetBinarys UserCommandName = "getBinarys"

	CmdForceAddCardServer   UserCommandName = "forceAddCardServer"
	CmdForceAddLoginServer  UserCommandName = "forceAddLoginServer"
	CmdForceAddTextServer   UserCommandName = "forceAddTextServer"
	CmdForceAddBinaryServer UserCommandName = "forceAddBinaryServer"

	CmdGetCardServer   UserCommandName = "getCardServer"
	CmdGetLoginServer  UserCommandName = "getLoginServer"
	CmdGetTextServer   UserCommandName = "getTextServer"
	CmdGetBinaryServer UserCommandName = "getBinaryServer"

	CmdExit UserCommandName = "syncExit"
)

type Options struct {
	Reg  bool `long:"reg" description:"registration"`
	Auth bool `long:"auth" description:"authentication"`

	AddCard   bool `long:"ncard" description:"add new card"`
	AddLogin  bool `long:"npwd" description:"add new pair login-password"`
	AddText   bool `long:"ntext" description:"add new text data"`
	AddBinary bool `long:"nbyte" description:"add new binary data"`

	UpdCard   bool `long:"ucard" description:"update card"`
	UpdLogin  bool `long:"upwd" description:"update pair login-password"`
	UpdText   bool `long:"utext" description:"update text data"`
	UpdBinary bool `long:"ubyte" description:"update binary data"`

	GetCard   bool `long:"gcard" description:"get card using number"`
	GetLogin  bool `long:"gpwd" description:"get pair login-password using prompt and login"`
	GetText   bool `long:"gtext" description:"get text data using prompt"`
	GetBinary bool `long:"gbyte" description:"get binary data using prompt"`

	GetCards   bool `long:"gcards" description:"get all cards"`
	GetLogins  bool `long:"gpwds" description:"get all pairs login-password"`
	GetTexts   bool `long:"gtexts" description:"get all text data"`
	GetBinarys bool `long:"gbytes" description:"get all binary data"`

	ForceAddCardServer   bool `long:"fcard" description:"add card to the server without verification"`
	ForceAddLoginServer  bool `long:"fpwd" description:"add pair login-password to the server without verification"`
	ForceAddTextServer   bool `long:"ftext" description:"add text data to the server without verification"`
	ForceAddBinaryServer bool `long:"fbyte" description:"add binary data to the server without verification"`

	GetCardServer   bool `long:"scard" description:"get card from server using number"`
	GetLoginServer  bool `long:"spwd" description:"get pair login-password from server using prompt and login"`
	GetTextServer   bool `long:"stext" description:"get text data from server using prompt"`
	GetBinaryServer bool `long:"sbyte" description:"get binary data from server using prompt"`

	UserLogin  string `short:"u" long:"userlogin" description:"user login"`
	Prompt     string `short:"p" long:"prompt" description:"hint for users data"`
	Login      string `short:"l" long:"login" description:"login for a login-password pair"`
	Note       string `short:"m" long:"comment" description:"data description or comment"`
	CardNumber string `short:"n" long:"number" description:"card number"`
	CardDate   string `short:"e" long:"date" description:"card expiry date"`
	CardCode   string `short:"v" long:"code" description:"card code"`
	Text       string `short:"t" long:"text" description:"text data"`
	Binary     string `short:"b" long:"byte" description:"binary data"`

	Exit bool `short:"x" long:"exit" description:"sync and exit"`
}

type UserArgs struct {
	AuthLogin  string
	Prompt     string
	Note       string
	CardNumber string
	CardDate   string
	CardCode   string
	Login      string
	Pwd        string
	Text       string
	Binary     string
}

var opt Options
var parser = flags.NewParser(&opt, flags.Default)

func ParseUserCmd(c string) (cmdName string, args UserArgs, err error) {
	cSpl := splitCmd(c)

	_, err = parser.ParseArgs(cSpl)
	if err != nil {
		var fErr *flags.Error
		if errors.As(err, &fErr) && fErr.Type == flags.ErrHelp {
			return "", UserArgs{}, nil
		}
		return "", UserArgs{}, err
	}

	switch {
	case opt.Reg:
		cmdName = CmdReg
		args = UserArgs{AuthLogin: opt.UserLogin}
		err = nil
	case opt.Auth:
		cmdName = CmdAuth
		args = UserArgs{AuthLogin: opt.UserLogin}
		err = nil

	case opt.AddCard:
		cmdName = CmdAddCard
		args = UserArgs{
			Prompt: opt.Prompt, Note: opt.Note, CardNumber: opt.CardNumber,
			CardDate: opt.CardDate, CardCode: opt.CardCode,
		}
		err = nil
	case opt.AddLogin:
		cmdName = CmdAddLogin
		args = UserArgs{Prompt: opt.Prompt, Note: opt.Note, Login: opt.Login}
		err = nil
	case opt.AddText:
		cmdName = CmdAddText
		args = UserArgs{Prompt: opt.Prompt, Note: opt.Note, Text: opt.Text}
		err = nil
	case opt.AddBinary:
		cmdName = CmdAddBinary
		args = UserArgs{Prompt: opt.Prompt, Note: opt.Note, Binary: opt.Binary}
		err = nil

	case opt.UpdCard:
		cmdName = CmdUpdCard
		args = UserArgs{
			Prompt: opt.Prompt, Note: opt.Note, CardNumber: opt.CardNumber,
			CardDate: opt.CardDate, CardCode: opt.CardCode,
		}
		err = nil
	case opt.UpdLogin:
		cmdName = CmdUpdLogin
		args = UserArgs{Prompt: opt.Prompt, Note: opt.Note, Login: opt.Login}
		err = nil
	case opt.UpdText:
		cmdName = CmdUpdText
		args = UserArgs{Prompt: opt.Prompt, Note: opt.Note, Text: opt.Text}
		err = nil
	case opt.UpdBinary:
		cmdName = CmdUpdBinary
		args = UserArgs{Prompt: opt.Prompt, Note: opt.Note, Binary: opt.Binary}
		err = nil

	case opt.GetCard:
		cmdName = CmdGetCard
		args = UserArgs{CardNumber: opt.CardNumber}
		err = nil
	case opt.GetLogin:
		cmdName = CmdGetLogin
		args = UserArgs{Prompt: opt.Prompt, Login: opt.Login}
		err = nil
	case opt.GetText:
		cmdName = CmdGetText
		args = UserArgs{Prompt: opt.Prompt}
		err = nil
	case opt.GetBinary:
		cmdName = CmdGetBinary
		args = UserArgs{Prompt: opt.Prompt}
		err = nil

	case opt.GetCards:
		cmdName = CmdGetCards
		args = UserArgs{}
		err = nil
	case opt.GetLogins:
		cmdName = CmdGetLogins
		args = UserArgs{}
		err = nil
	case opt.GetTexts:
		cmdName = CmdGetTexts
		args = UserArgs{}
		err = nil
	case opt.GetBinarys:
		cmdName = CmdGetBinarys
		args = UserArgs{}
		err = nil

	case opt.ForceAddCardServer:
		cmdName = CmdForceAddCardServer
		args = UserArgs{CardNumber: opt.CardNumber}
		err = nil
	case opt.ForceAddLoginServer:
		cmdName = CmdForceAddLoginServer
		args = UserArgs{Prompt: opt.Prompt, Login: opt.Login}
		err = nil
	case opt.ForceAddTextServer:
		cmdName = CmdForceAddTextServer
		args = UserArgs{Prompt: opt.Prompt}
		err = nil
	case opt.ForceAddBinaryServer:
		cmdName = CmdForceAddBinaryServer
		args = UserArgs{Prompt: opt.Prompt}
		err = nil

	case opt.GetCardServer:
		cmdName = CmdGetCardServer
		args = UserArgs{CardNumber: opt.CardNumber}
		err = nil
	case opt.GetLoginServer:
		cmdName = CmdGetLoginServer
		args = UserArgs{Prompt: opt.Prompt, Login: opt.Login}
		err = nil
	case opt.GetTextServer:
		cmdName = CmdGetTextServer
		args = UserArgs{Prompt: opt.Prompt}
		err = nil
	case opt.GetBinaryServer:
		cmdName = CmdGetBinaryServer
		args = UserArgs{Prompt: opt.Prompt}
		err = nil

	case opt.Exit:
		cmdName = CmdExit
		args = UserArgs{}
		err = nil

	default:
		cmdName = ""
		args = UserArgs{}
		err = errors.New("unclear command")
	}

	clearOpt(&opt)
	return
}
