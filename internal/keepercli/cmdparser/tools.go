package cmdparser

import "strings"

func splitCmd(cmd string) []string {
	cmd = strings.TrimSpace(cmd)
	sl := strings.Split(cmd, " -")
	for i := 1; i < len(sl); i++ {
		sl[i] = "-" + sl[i]
	}
	return sl
}

func clearOpt(opt *Options) error {
	opt.AddBinary = false
	opt.AddCard = false
	opt.AddLogin = false
	opt.AddText = false
	opt.Auth = false
	opt.Binary = ""
	opt.CardCode = ""
	opt.CardDate = ""
	opt.CardDate = ""
	opt.CardNumber = ""
	opt.Exit = false
	opt.Force = false
	opt.ForceAddBinaryServer = false
	opt.ForceAddCardServer = false
	opt.ForceAddLoginServer = false
	opt.ForceAddTextServer = false
	opt.GetBinary = false
	opt.GetBinaryServer = false
	opt.GetBinarys = false
	opt.GetCard = false
	opt.GetCardServer = false
	opt.GetCards = false
	opt.GetLogin = false
	opt.GetLoginServer = false
	opt.GetLogins = false
	opt.GetText = false
	opt.GetTextServer = false
	opt.GetTexts = false
	opt.Login = ""
	opt.Note = ""
	opt.Prompt = ""
	opt.Reg = false
	opt.Text = ""
	opt.UpdBinary = false
	opt.UpdCard = false
	opt.UpdLogin = false
	opt.UpdText = false
	opt.UserLogin = ""

	return nil
}
