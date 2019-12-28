package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"log"
	"net/http"
)

func md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func enbase(str []byte) string {
	encodeString := base64.StdEncoding.EncodeToString(str)
	return encodeString
}

func debase(str string) string {
	decodeBytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Fatalln(err)
	}
	return string(decodeBytes)
}

func errorHandle(err error, w http.ResponseWriter) {
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}
