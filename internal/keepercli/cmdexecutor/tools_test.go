package cmdexecutor

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Julia-ivv/info-keeper.git/internal/keepercli/cmdparser"
	"github.com/Julia-ivv/info-keeper.git/internal/keepercli/storage"
	pb "github.com/Julia-ivv/info-keeper.git/internal/proto/pb"
)

var (
	testTime = "2024-01-02T15:04:05Z"
	testCard = storage.Card{
		Prompt:    []byte{1, 89, 253, 125, 56, 183, 82, 131, 15, 240, 8, 156, 250, 125, 1, 177, 147, 57, 250, 227, 76, 138},
		Number:    []byte{64, 25, 161, 98, 176, 231, 86, 61, 183, 72, 204, 32, 140, 43, 186, 58, 106, 144, 29},
		Date:      []byte{64, 25, 189, 34, 124, 190, 142, 192, 157, 21, 165, 199, 3, 27, 118, 168, 188, 32, 228, 236, 178},
		Code:      []byte{68, 30, 167, 197, 156, 30, 209, 39, 36, 154, 28, 24, 129, 206, 216, 214, 64, 234, 228},
		Note:      []byte{31, 68, 230, 117, 130, 201, 228, 68, 141, 201, 130, 243, 11, 194, 89, 155, 192, 52, 66, 45},
		TimeStamp: testTime,
	}
	testLoginPwd = storage.LoginPwd{
		Prompt:    []byte{1, 89, 253, 125, 56, 183, 82, 131, 15, 240, 8, 156, 250, 125, 1, 177, 147, 57, 250, 227, 76, 138},
		Login:     []byte{29, 68, 245, 121, 38, 4, 57, 2, 145, 229, 80, 214, 52, 137, 222, 41, 34, 23, 14, 98, 26},
		Pwd:       []byte{1, 92, 246, 57, 219, 189, 107, 91, 188, 80, 32, 85, 96, 176, 13, 173, 169, 251, 107},
		Note:      []byte{31, 68, 230, 117, 130, 201, 228, 68, 141, 201, 130, 243, 11, 194, 89, 155, 192, 52, 66, 45},
		TimeStamp: testTime,
	}
	testTextRecord = storage.TextRecord{
		Prompt:    []byte{1, 89, 253, 125, 56, 183, 82, 131, 15, 240, 8, 156, 250, 125, 1, 177, 147, 57, 250, 227, 76, 138},
		Data:      []byte{5, 78, 234, 100, 213, 202, 147, 13, 147, 58, 30, 16, 43, 161, 175, 6, 200, 36, 87, 250},
		Note:      []byte{31, 68, 230, 117, 130, 201, 228, 68, 141, 201, 130, 243, 11, 194, 89, 155, 192, 52, 66, 45},
		TimeStamp: testTime,
	}
	testBinaryRecord = storage.BinaryRecord{
		Prompt:    []byte{1, 89, 253, 125, 56, 183, 82, 131, 15, 240, 8, 156, 250, 125, 1, 177, 147, 57, 250, 227, 76, 138},
		Data:      []byte{19, 82, 230, 117, 221, 110, 161, 236, 11, 24, 168, 191, 253, 202, 73, 174, 150, 231, 168, 212},
		Note:      []byte{31, 68, 230, 117, 130, 201, 228, 68, 141, 201, 130, 243, 11, 194, 89, 155, 192, 52, 66, 45},
		TimeStamp: testTime,
	}
	testPbCard = &pb.UserCard{
		Prompt:    testCard.Prompt,
		Number:    testCard.Number,
		Date:      testCard.Date,
		Code:      testCard.Code,
		Note:      testCard.Note,
		TimeStamp: testCard.TimeStamp,
	}
	testPbLogin = &pb.UserLoginPwd{
		Prompt:    testLoginPwd.Prompt,
		Login:     testLoginPwd.Login,
		Pwd:       testLoginPwd.Pwd,
		Note:      testLoginPwd.Note,
		TimeStamp: testLoginPwd.TimeStamp,
	}
	testPbTextRecord = &pb.UserTextRecord{
		Prompt:    testTextRecord.Prompt,
		Data:      testTextRecord.Data,
		Note:      testTextRecord.Note,
		TimeStamp: testTextRecord.TimeStamp,
	}
	testPbBinaryRecord = &pb.UserBinaryRecord{
		Prompt:    testBinaryRecord.Prompt,
		Data:      testBinaryRecord.Data,
		Note:      testBinaryRecord.Note,
		TimeStamp: testBinaryRecord.TimeStamp,
	}
	testArgs = cmdparser.UserArgs{
		Prompt:     "prompt",
		Note:       "note",
		CardNumber: "123",
		CardDate:   "12/24",
		CardCode:   "555",
		Login:      "login",
		Pwd:        "pwd",
		Text:       "text",
		Binary:     "byte",
	}
	testEnArgs = EncryptArgs{
		Prompt:     []byte{1, 89, 253, 125, 56, 183, 82, 131, 15, 240, 8, 156, 250, 125, 1, 177, 147, 57, 250, 227, 76, 138},
		Note:       []byte{31, 68, 230, 117, 130, 201, 228, 68, 141, 201, 130, 243, 11, 194, 89, 155, 192, 52, 66, 45},
		CardNumber: []byte{64, 25, 161, 98, 176, 231, 86, 61, 183, 72, 204, 32, 140, 43, 186, 58, 106, 144, 29},
		CardDate:   []byte{64, 25, 189, 34, 124, 190, 142, 192, 157, 21, 165, 199, 3, 27, 118, 168, 188, 32, 228, 236, 178},
		CardCode:   []byte{68, 30, 167, 197, 156, 30, 209, 39, 36, 154, 28, 24, 129, 206, 216, 214, 64, 234, 228},
		Login:      []byte{29, 68, 245, 121, 38, 4, 57, 2, 145, 229, 80, 214, 52, 137, 222, 41, 34, 23, 14, 98, 26},
		Pwd:        []byte{1, 92, 246, 57, 219, 189, 107, 91, 188, 80, 32, 85, 96, 176, 13, 173, 169, 251, 107},
		Text:       []byte{5, 78, 234, 100, 213, 202, 147, 13, 147, 58, 30, 16, 43, 161, 175, 6, 200, 36, 87, 250},
		Binary:     []byte{19, 82, 230, 117, 221, 110, 161, 236, 11, 24, 168, 191, 253, 202, 73, 174, 150, 231, 168, 212},
	}
)

func TestCardToPb(t *testing.T) {
	pbc := cardToPb(testCard)
	assert.Equal(t, testPbCard, pbc)
}

func TestCardsToPb(t *testing.T) {
	pbc := cardsToPb([]storage.Card{testCard, testCard})
	assert.Equal(t, []*pb.UserCard{testPbCard, testPbCard}, pbc)
}

func TestLoginToPb(t *testing.T) {
	pbl := loginToPb(testLoginPwd)
	assert.Equal(t, testPbLogin, pbl)
}

func TestLoginsToPb(t *testing.T) {
	pbl := loginsToPb([]storage.LoginPwd{testLoginPwd, testLoginPwd})
	assert.Equal(t, []*pb.UserLoginPwd{testPbLogin, testPbLogin}, pbl)
}

func TestTextToPb(t *testing.T) {
	pbt := textToPb(testTextRecord)
	assert.Equal(t, testPbTextRecord, pbt)
}

func TestTextsToPb(t *testing.T) {
	pbt := textsToPb([]storage.TextRecord{testTextRecord, testTextRecord})
	assert.Equal(t, []*pb.UserTextRecord{testPbTextRecord, testPbTextRecord}, pbt)
}

func TestBinaryToPb(t *testing.T) {
	pbb := binaryToPb(testBinaryRecord)
	assert.Equal(t, testPbBinaryRecord, pbb)
}

func TestBinarysToPb(t *testing.T) {
	pbb := binarysToPb([]storage.BinaryRecord{testBinaryRecord, testBinaryRecord})
	assert.Equal(t, []*pb.UserBinaryRecord{testPbBinaryRecord, testPbBinaryRecord}, pbb)
}

func TestPbToCard(t *testing.T) {
	d := pbToCard(testPbCard)
	assert.Equal(t, testCard, d)
}

func TestPbToCards(t *testing.T) {
	d := pbToCards([]*pb.UserCard{testPbCard, testPbCard})
	assert.Equal(t, []storage.Card{testCard, testCard}, d)
}

func TestPbToLogin(t *testing.T) {
	d := pbToLogin(testPbLogin)
	assert.Equal(t, testLoginPwd, d)
}

func TestPbToLogins(t *testing.T) {
	d := pbToLogins([]*pb.UserLoginPwd{testPbLogin, testPbLogin})
	assert.Equal(t, []storage.LoginPwd{testLoginPwd, testLoginPwd}, d)
}

func TestPbToText(t *testing.T) {
	d := pbToText(testPbTextRecord)
	assert.Equal(t, testTextRecord, d)
}

func TestPbToTexts(t *testing.T) {
	d := pbToTexts([]*pb.UserTextRecord{testPbTextRecord, testPbTextRecord})
	assert.Equal(t, []storage.TextRecord{testTextRecord, testTextRecord}, d)
}

func TestPbToBinary(t *testing.T) {
	d := pbToBinary(testPbBinaryRecord)
	assert.Equal(t, testBinaryRecord, d)
}

func TestPbToBinarys(t *testing.T) {
	d := pbToBinarys([]*pb.UserBinaryRecord{testPbBinaryRecord, testPbBinaryRecord})
	assert.Equal(t, []storage.BinaryRecord{testBinaryRecord, testBinaryRecord}, d)
}

func TestEncryptArgs(t *testing.T) {
	enA, err := encryptArgs(testArgs)
	if assert.NoError(t, err) {
		assert.Equal(t, testEnArgs, enA)
	}
}

func TestDecryptCard(t *testing.T) {
	d, err := decryptCard(testCard)
	if assert.NoError(t, err) {
		assert.NotEmpty(t, d)
	}
}

func TestDecryptLoginPwd(t *testing.T) {
	d, err := decryptLoginPwd(testLoginPwd)
	if assert.NoError(t, err) {
		assert.NotEmpty(t, d)
	}
}

func TestDecryptTextRecord(t *testing.T) {
	d, err := decryptTextRecord(testTextRecord)
	if assert.NoError(t, err) {
		assert.NotEmpty(t, d)
	}
}

func TestDecryptBinaryRecord(t *testing.T) {
	d, err := decryptBinaryRecord(testBinaryRecord)
	if assert.NoError(t, err) {
		assert.NotEmpty(t, d)
	}
}
