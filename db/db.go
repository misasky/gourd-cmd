package db

import (
	"flag"
	"fmt"
	"linson/gourd-cmd/utils"
	"os"

	"github.com/boltdb/bolt"
)

const (
	FilePath   = "./gourd"
	GourdKey   = "Q7C^23oPC09Cm12E"
	fileLength = 500000
)

var DBName = []byte("gourd")
var Client *bolt.DB

func Init() bool {
	flag.Parse()
	exist, err := utils.PathExists(FilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(10001)
	}
	if !exist {
		err := os.Mkdir(FilePath, os.ModePerm)
		if err != nil {
			panic(err)
		}
		f, err := os.Create(FilePath + "/grourd")
		if err != nil {
			panic(err)
		}
		if err := f.Close(); err != nil {
			panic(err)
		}
	}
	d, err := bolt.Open(FilePath+"/grourd", 0600, nil)
	if err != nil {
		panic(err)
	}
	Client = d
	if err := Client.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(DBName)
		if b == nil {
			_, err := tx.CreateBucket(DBName)
			if err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
	return exist
}
