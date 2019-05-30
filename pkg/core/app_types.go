package core

import (
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"
)

type App struct {
	Name string
	// TODO: Consider a name that better describes the intent (e.g. "guardid" or something)
	Hashid AppId
}

type AppId []byte

func (appId AppId) String() string {
	return fmt.Sprintf("%x", string(appId))
}

func DecodeAppId(src string) (AppId, error) {
	bytes, err := hex.DecodeString(src)
	if err != nil {
		return nil, errors.New("app ID must be a valid hex string")
	}

	return bytes, nil
}
