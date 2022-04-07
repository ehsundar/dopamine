package token

import (
	"errors"
	"io/fs"
	"io/ioutil"
	"math/rand"
	"os"
)

const (
	signingKeyFileName = "dopamine-secret.txt"
	letterBytes        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func LoadSigningKey() []byte {
	contents, err := ioutil.ReadFile(signingKeyFileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			key := randStringBytes(64)
			err = ioutil.WriteFile(signingKeyFileName, key, fs.ModePerm)
			if err != nil {
				panic(err)
			}
			return key
		}
		panic(err)
	}
	return contents
}

func randStringBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return b
}
