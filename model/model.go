package model

import (
	capnp_go "blumer-ms-refers/capnp"
	"bytes"
	"capnproto.org/go/capnp/v3"
	"strings"
)

type Profile struct {
	UserID   string  `json:"user_id"`
	Username string  `json:"username"`
	IsActive bool    `json:"is_active"`
	Reward   float64 `json:"reward"`
}

type UserCtx struct {
	UserID string `json:"user"`
	Role   string `json:"role"`
	Ip     string `json:"X-Forward-Ip"`
}

func DecodeProfile(msg []byte) (*Profile, error) {
	reader := strings.NewReader(string(msg))
	msgCompiled, err := capnp.NewDecoder(reader).Decode()
	if err != nil && err.Error() != "EOF" {
		return nil, err
	}

	msgDecompiled, err := capnp_go.ReadRootProfileInfo(msgCompiled)
	if err != nil {
		return nil, err
	}

	userID, _ := msgDecompiled.UserId()
	username, _ := msgDecompiled.Username()
	isActive := msgDecompiled.IsActive()
	reward := msgDecompiled.Reward()

	return &Profile{
		UserID:   userID,
		Username: username,
		IsActive: isActive,
		Reward:   reward,
	}, nil
}

func DecodeUserID(msg []byte) (*string, error) {
	reader := strings.NewReader(string(msg))
	msgCompiled, err := capnp.NewDecoder(reader).Decode()
	if err != nil && err.Error() != "EOF" {
		return nil, err
	}

	msgDecompiled, err := capnp_go.ReadRootProfileWallet(msgCompiled)

	userID, _ := msgDecompiled.UserId()

	return &userID, nil
}

func EncodeWalletReward(profile Profile) ([]byte, error) {
	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		return nil, err
	}

	wallet, err := capnp_go.NewRootProfileInfo(seg)
	if err != nil {
		return nil, err
	}

	_ = wallet.SetUserId(profile.UserID)
	wallet.SetReward(profile.Reward)
	buf := new(bytes.Buffer)
	err = capnp.NewEncoder(buf).Encode(msg)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
