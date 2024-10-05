package albiononline_dumper

import (
	"bytes"
	"compress/gzip"
	"crypto/cipher"
	"crypto/des"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	xj "github.com/basgys/goxml2json"
)

var (
	Key = []byte{48, 239, 114, 71, 66, 242, 4, 50}
	Iv  = []byte{14, 166, 220, 137, 219, 237, 220, 79}
)

func GetBinFolder(gameFolder string) string {
	return filepath.Join(gameFolder, "./Albion-Online_Data/StreamingAssets/GameData")
}

func Dump(gameFolder string, outputFolder string) {
	binFolder := GetBinFolder(gameFolder)

	err := filepath.Walk(binFolder, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		// Skipping dirs
		if info.IsDir() {
			return nil
		}

		// Skipping other files than .bin
		if info.Name()[len(info.Name())-4:] != ".bin" {
			return nil
		}

		fmt.Println("Extracting", path[len(binFolder)+1:])
		OutputFile(binFolder, outputFolder, path, DecryptFile(path))
		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", binFolder, err)
		return
	}
}

func DecryptFile(path string) []byte {
	data, err := os.ReadFile(path)

	if err != nil {
		panic(err)
	}

	block, err := des.NewCipher(Key)

	if err != nil {
		panic(err)
	}

	blockMode := cipher.NewCBCDecrypter(block, Iv)
	origData := make([]byte, len(data))
	blockMode.CryptBlocks(origData, data)
	origData = PKCS5UnPadding(origData)

	r, err := gzip.NewReader(bytes.NewReader(origData))
	if err != nil {
		panic(err)
	}

	var res bytes.Buffer
	_, err = res.ReadFrom(r)
	if err != nil {
		panic(err)
	}

	return res.Bytes()
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

func OutputFile(gamePath string, outputFolder string, binPath string, data []byte) {
	// Check if dir exist, if not create it
	res, err := os.Stat(filepath.Join(outputFolder, filepath.Dir(binPath[len(gamePath)+1:])))
	if os.IsNotExist(err) {
		os.MkdirAll(filepath.Join(outputFolder, filepath.Dir(binPath)[len(gamePath)+1:]), 0777)
	} else if !res.IsDir() {
		panic("Output folder is not a directory")
	}

	// Write XML file
	os.WriteFile(
		filepath.Join(outputFolder, binPath[len(gamePath)+1:len(binPath)-4]+".xml"),
		data,
		0666,
	)

	json, err := xj.Convert(bytes.NewReader(data))
	if err != nil {
		panic("Couldn't convert XML to JSON")
	}

	// Write JSON file
	os.WriteFile(
		filepath.Join(outputFolder, binPath[len(gamePath)+1:len(binPath)-4]+".json"),
		json.Bytes(),
		0666,
	)
}
