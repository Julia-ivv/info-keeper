package cmdexecutor

import (
	"github.com/Julia-ivv/info-keeper.git/internal/keepercli/cmdparser"
	"github.com/Julia-ivv/info-keeper.git/internal/keepercli/cryptor"
	"github.com/Julia-ivv/info-keeper.git/internal/keepercli/storage"
	pb "github.com/Julia-ivv/info-keeper.git/internal/proto/pb"
)

func cardToPb(c storage.Card) *pb.UserCard {
	return &pb.UserCard{
		Prompt:    c.Prompt,
		Number:    c.Number,
		Date:      c.Date,
		Code:      c.Code,
		Note:      c.Note,
		TimeStamp: c.TimeStamp,
	}
}

func cardsToPb(cs []storage.Card) []*pb.UserCard {
	pbC := make([]*pb.UserCard, 0, len(cs))
	for _, v := range cs {
		pbC = append(pbC, cardToPb(v))
	}
	return pbC
}

func loginToPb(l storage.LoginPwd) *pb.UserLoginPwd {
	return &pb.UserLoginPwd{
		Prompt:    l.Prompt,
		Login:     l.Login,
		Pwd:       l.Pwd,
		Note:      l.Note,
		TimeStamp: l.TimeStamp,
	}
}

func loginsToPb(ls []storage.LoginPwd) []*pb.UserLoginPwd {
	pbL := make([]*pb.UserLoginPwd, 0, len(ls))
	for _, v := range ls {
		pbL = append(pbL, loginToPb(v))
	}
	return pbL
}

func textToPb(t storage.TextRecord) *pb.UserTextRecord {
	return &pb.UserTextRecord{
		Prompt:    t.Prompt,
		Data:      t.Data,
		Note:      t.Note,
		TimeStamp: t.TimeStamp,
	}
}

func textsToPb(ts []storage.TextRecord) []*pb.UserTextRecord {
	pbT := make([]*pb.UserTextRecord, 0, len(ts))
	for _, v := range ts {
		pbT = append(pbT, textToPb(v))
	}
	return pbT
}

func binaryToPb(b storage.BinaryRecord) *pb.UserBinaryRecord {
	return &pb.UserBinaryRecord{
		Prompt:    b.Prompt,
		Data:      b.Data,
		Note:      b.Note,
		TimeStamp: b.TimeStamp,
	}
}

func binarysToPb(bs []storage.BinaryRecord) []*pb.UserBinaryRecord {
	pbB := make([]*pb.UserBinaryRecord, 0, len(bs))
	for _, v := range bs {
		pbB = append(pbB, binaryToPb(v))
	}
	return pbB
}

func pbToCard(c *pb.UserCard) storage.Card {
	return storage.Card{
		Prompt:    c.Prompt,
		Number:    c.Number,
		Date:      c.Date,
		Code:      c.Code,
		Note:      c.Note,
		TimeStamp: c.TimeStamp,
	}
}

func pbToCards(pbCs []*pb.UserCard) []storage.Card {
	cs := make([]storage.Card, 0, len(pbCs))
	for _, v := range pbCs {
		cs = append(cs, pbToCard(v))
	}
	return cs
}

func pbToLogin(l *pb.UserLoginPwd) storage.LoginPwd {
	return storage.LoginPwd{
		Prompt:    l.Prompt,
		Login:     l.Login,
		Pwd:       l.Pwd,
		Note:      l.Note,
		TimeStamp: l.TimeStamp,
	}
}

func pbToLogins(pbLs []*pb.UserLoginPwd) []storage.LoginPwd {
	ls := make([]storage.LoginPwd, 0, len(pbLs))
	for _, v := range pbLs {
		ls = append(ls, pbToLogin(v))
	}
	return ls
}

func pbToText(t *pb.UserTextRecord) storage.TextRecord {
	return storage.TextRecord{
		Prompt:    t.Prompt,
		Data:      t.Data,
		Note:      t.Note,
		TimeStamp: t.TimeStamp,
	}
}

func pbToTexts(pbTs []*pb.UserTextRecord) []storage.TextRecord {
	ts := make([]storage.TextRecord, 0, len(pbTs))
	for _, v := range pbTs {
		ts = append(ts, pbToText(v))
	}
	return ts
}

func pbToBinary(b *pb.UserBinaryRecord) storage.BinaryRecord {
	return storage.BinaryRecord{
		Prompt:    b.Prompt,
		Data:      b.Data,
		Note:      b.Note,
		TimeStamp: b.TimeStamp,
	}
}

func pbToBinarys(pbBs []*pb.UserBinaryRecord) []storage.BinaryRecord {
	bs := make([]storage.BinaryRecord, 0, len(pbBs))
	for _, v := range pbBs {
		bs = append(bs, pbToBinary(v))
	}
	return bs
}

type EncryptArgs struct {
	Prompt     []byte
	Note       []byte
	CardNumber []byte
	CardDate   []byte
	CardCode   []byte
	Login      []byte
	Pwd        []byte
	Text       []byte
	Binary     []byte
}

func encryptArgs(a cmdparser.UserArgs) (enA EncryptArgs, err error) {
	if a.Prompt != "" {
		enA.Prompt, err = cryptor.EncryptsString(a.Prompt)
		if err != nil {
			return
		}
	}
	if a.Note != "" {
		enA.Note, err = cryptor.EncryptsString(a.Note)
		if err != nil {
			return
		}
	}
	if a.CardNumber != "" {
		enA.CardNumber, err = cryptor.EncryptsString(a.CardNumber)
		if err != nil {
			return
		}
	}
	if a.CardDate != "" {
		enA.CardDate, err = cryptor.EncryptsString(a.CardDate)
		if err != nil {
			return
		}
	}
	if a.CardCode != "" {
		enA.CardCode, err = cryptor.EncryptsString(a.CardCode)
		if err != nil {
			return
		}
	}
	if a.Login != "" {
		enA.Login, err = cryptor.EncryptsString(a.Login)
		if err != nil {
			return
		}
	}
	if a.Pwd != "" {
		enA.Pwd, err = cryptor.EncryptsString(a.Pwd)
		if err != nil {
			return
		}
	}
	if a.Text != "" {
		enA.Text, err = cryptor.EncryptsString(a.Text)
		if err != nil {
			return
		}
	}
	if a.Binary != "" {
		enA.Binary, err = cryptor.EncryptsString(a.Binary)
		if err != nil {
			return
		}
	}
	return
}

func decryptCard(c storage.Card) (uc UserCard, err error) {
	uc.Prompt, err = cryptor.Decrypts(c.Prompt)
	if err != nil {
		return
	}
	uc.Number, err = cryptor.Decrypts(c.Number)
	if err != nil {
		return
	}
	uc.Date, err = cryptor.Decrypts(c.Date)
	if err != nil {
		return
	}
	uc.Code, err = cryptor.Decrypts(c.Code)
	if err != nil {
		return
	}
	uc.Note, err = cryptor.Decrypts(c.Note)
	if err != nil {
		return
	}
	uc.TimeStamp = c.TimeStamp
	return
}

func decryptLoginPwd(l storage.LoginPwd) (ul UserLoginPwd, err error) {
	ul.Prompt, err = cryptor.Decrypts(l.Prompt)
	if err != nil {
		return
	}
	ul.Login, err = cryptor.Decrypts(l.Login)
	if err != nil {
		return
	}
	ul.Pwd, err = cryptor.Decrypts(l.Pwd)
	if err != nil {
		return
	}
	ul.Note, err = cryptor.Decrypts(l.Note)
	if err != nil {
		return
	}
	ul.TimeStamp = l.TimeStamp
	return
}

func decryptTextRecord(t storage.TextRecord) (ut UserTextRecord, err error) {
	ut.Prompt, err = cryptor.Decrypts(t.Prompt)
	if err != nil {
		return
	}
	ut.Data, err = cryptor.Decrypts(t.Data)
	if err != nil {
		return
	}
	ut.Note, err = cryptor.Decrypts(t.Note)
	if err != nil {
		return
	}
	ut.TimeStamp = t.TimeStamp
	return
}

func decryptBinaryRecord(b storage.BinaryRecord) (ub UserBinaryRecord, err error) {
	ub.Prompt, err = cryptor.Decrypts(b.Prompt)
	if err != nil {
		return
	}
	ub.Data, err = cryptor.DecryptsInByte(b.Data)
	if err != nil {
		return
	}
	ub.Note, err = cryptor.Decrypts(b.Note)
	if err != nil {
		return
	}
	ub.TimeStamp = b.TimeStamp
	return
}
