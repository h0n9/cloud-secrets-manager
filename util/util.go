package util

import (
	"encoding/base64"
	"io/ioutil"
	"os"
)

func GetEnv(key, fallback string) string {
	if value, exist := os.LookupEnv(key); exist {
		return value
	}
	return fallback
}

func EncodeBase64StrToStr(input string) string {
	return base64.RawStdEncoding.EncodeToString([]byte(input))
}

func DecodeBase64StrToStr(input string) (string, error) {
	output, err := DecodeBase64StrToBytes(input)
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func DecodeBase64StrToBytes(input string) ([]byte, error) {
	return base64.RawStdEncoding.DecodeString(input)
}

func ReadFileToBytes(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func ReadFileToStr(filename string) (string, error) {
	data, err := ReadFileToBytes(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func WriteBytesToFile(filename string, data []byte) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func WriteStrToFile(filename string, data string) error {
	return WriteBytesToFile(filename, []byte(data))
}
