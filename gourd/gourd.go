package gourd

import (
	"fmt"
	"linson/gourd-cmd/aes"
	"linson/gourd-cmd/db"
	"os"
	"strings"

	"github.com/boltdb/bolt"
)

func InitGourdKey(pwd string) error {
	encKey, err := aes.AesEncrypt(db.GourdKey, pwd)
	if err != nil {
		return err
	}
	if err := db.Client.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(db.DBName)
		if err := b.Put([]byte("aes_key"), []byte(encKey)); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func CheckGourdKey(pwd string) error {
	if err := db.Client.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(db.DBName)
		v := b.Get([]byte("aes_key"))
		decKey, err := aes.AesDecrypt(string(v), pwd)
		if err != nil {
			return err
		}
		if decKey != db.GourdKey {
			fmt.Println("凭证不正确，请重新输入。例如 -p 1234567890123456")
			os.Exit(1)
		}
		return nil
	}); err != nil {
		return err
	}
	fmt.Println("验证成功...")
	return nil
}

func Set(url, account, pwd, key string) error {
	encAcc, err := aes.AesEncrypt(account, key)
	if err != nil {
		return err
	}
	encPwd, err := aes.AesEncrypt(pwd, key)
	if err != nil {
		return err
	}
	if err := db.Client.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(db.DBName)
		acc := strings.Join([]string{encAcc, encPwd}, "|")
		if err := b.Put([]byte(url), []byte(acc)); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func Get(url, key string) (string, string, error) {
	var acc, pwd string
	if err := db.Client.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(db.DBName)
		d := b.Get([]byte(url))
		strs := strings.Split(string(d), "|")
		encAcc, err := aes.AesDecrypt(strs[0], key)
		if err != nil {
			return err
		}
		acc = encAcc
		encPwd, err := aes.AesDecrypt(strs[1], key)
		if err != nil {
			return err
		}
		pwd = encPwd
		return nil
	}); err != nil {
		return "", "", err
	}
	return acc, pwd, nil
}
