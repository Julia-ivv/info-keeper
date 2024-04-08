package cmdexecutor

import (
	"errors"

	"github.com/Julia-ivv/info-keeper.git/internal/keepercli/cmdparser"
	"github.com/Julia-ivv/info-keeper.git/internal/keepercli/storage"
	pb "github.com/Julia-ivv/info-keeper.git/internal/proto/pb"
	"github.com/Julia-ivv/info-keeper.git/pkg/logger"
)

type DataPrinter interface {
	PrintData()
}

var cmds = make(map[string]func(args cmdparser.UserArgs, cl pb.InfoKeeperClient, repo storage.Repositorier) (DataPrinter, error))

func init() {
	cmds[cmdparser.CmdReg] = regExec
	cmds[cmdparser.CmdAuth] = authExec
	cmds[cmdparser.CmdExit] = exitExec

	cmds[cmdparser.CmdAddCard] = addCardExec
	cmds[cmdparser.CmdAddLogin] = addLoginExec
	cmds[cmdparser.CmdAddText] = addTextExec
	cmds[cmdparser.CmdAddBinary] = addBinaryExec

	cmds[cmdparser.CmdUpdCard] = updCardExec
	cmds[cmdparser.CmdUpdLogin] = updLoginExec
	cmds[cmdparser.CmdUpdText] = updTextExec
	cmds[cmdparser.CmdUpdBinary] = updBinaryExec

	cmds[cmdparser.CmdGetCard] = getCardExec
	cmds[cmdparser.CmdGetLogin] = getLoginExec
	cmds[cmdparser.CmdGetText] = getTextExec
	cmds[cmdparser.CmdGetBinary] = getBinaryExec

	cmds[cmdparser.CmdGetCards] = getCardsExec
	cmds[cmdparser.CmdGetLogins] = getLoginsExec
	cmds[cmdparser.CmdGetTexts] = getTextsExec
	cmds[cmdparser.CmdGetBinarys] = getBinarysExec

	cmds[cmdparser.CmdForceAddCardServer] = forceAddCardServerExec
	cmds[cmdparser.CmdForceAddLoginServer] = forceAddLoginServerExec
	cmds[cmdparser.CmdForceAddTextServer] = forceAddTextServerExec
	cmds[cmdparser.CmdForceAddBinaryServer] = forceAddBinaryServerExec

	cmds[cmdparser.CmdGetCardServer] = getCardServerExec
	cmds[cmdparser.CmdGetLoginServer] = getLoginServerExec
	cmds[cmdparser.CmdGetTextServer] = getTextServerExec
	cmds[cmdparser.CmdGetBinaryServer] = getBinaryServerExec
}

func ExecuteCmd(userCmd string, userArgs cmdparser.UserArgs, cl pb.InfoKeeperClient, repo storage.Repositorier) (DataPrinter, error) {
	fn := cmds[userCmd]
	if fn == nil {
		logger.ZapSugar.Infoln("command function not found")
		return nil, errors.New("command function not found")
	}
	res, err := fn(userArgs, cl, repo)
	if err != nil {
		return nil, err
	}

	return res, nil
}
