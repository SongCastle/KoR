package lib

import (
	"github.com/SongCastle/KoR/lib/encryptor"
	"github.com/SongCastle/KoR/lib/jwt"
)

func SetUp() {
	encryptor.Init()
	jwt.Init()
}
