package main

import (
	"flag"
	"fmt"
	"linson/gourd-cmd/db"
	"linson/gourd-cmd/gourd"
	"os"
	"strings"
	"unicode/utf8"
)

var option = flag.String("o", "", "操作选项\nverify: 验证密钥并进入命令行界面\n")
var key = flag.String("p", "", "保存帐号时加密的凭证, 需要牢记\n")

func init() {
	flag.Parse()
	if utf8.RuneCountInString(*key) != 16 {
		fmt.Println("> 请输入16位字符并牢记, 例如: -p 1234567890123456")
		os.Exit(1)
	}
}

func main() {
	ok := db.Init()
	if !ok {
		if err := gourd.InitGourdKey(*key); err != nil {
			fmt.Println(err)
		}
	} else {
		switch *option {
		case "verify":
			if err := gourd.CheckGourdKey(*key); err != nil {
				fmt.Println(err)
				os.Exit(10001)
			}
		default:
			fmt.Println("> 不支持的操作")
			os.Exit(1)
		}
	}
opt:
	var o string
	fmt.Println("------------------------------------------------------------------")
	fmt.Println("> 请输入操作，支持 set, get 后续将支持:search, list. 例如: set")
	fmt.Println("------------------------------------------------------------------")
	fmt.Scanln(&o)
	switch o {
	case "set":
	set:
		var f = func() []string {
			fmt.Println("------------------------------------------------------------------")
			fmt.Println("> 请输入存储的网址, 账户, 密码 使用'|'分隔, 例如: www.baidu.com|admin|admin")
			fmt.Println("> 建议使用的密码:")
			for i := 1; i < 4; i++ {
				fmt.Println(fmt.Sprintf("%d：%s", i, gourd.GetRandStr()))
			}
			fmt.Println("------------------------------------------------------------------")
			var d string
			fmt.Scanln(&d)
			return strings.Split(d, "|")
		}
		strs := f()
		if len(strs) < 3 {
			fmt.Println("> 格式错误,请重新输入...")
			goto set
		}
		if err := gourd.Set(strs[0], strs[1], strs[2], *key); err != nil {
			fmt.Println("> 保存账号时出现错误...")
			goto set
		}
		fmt.Println("保存成功...")
		fmt.Println("------------------------------------------------------------------")
		goto opt
	case "get":
	get:
		var f = func() string {
			fmt.Println("------------------------------------------------------------------")
			fmt.Println("> 请输入存储的网址, 例如: www.baidu.com")
			fmt.Println("------------------------------------------------------------------")
			var d string
			fmt.Scanln(&d)
			return d
		}
		url := f()
		if url == "" {
			fmt.Println("> 格式错误,请重新输入...")
			goto get
		}
		acc, pwd, err := gourd.Get(url, *key)
		if err != nil {
			fmt.Println("> 获取账号时出现错误...")
			goto get
		}
		fmt.Println("========账号信息========")
		fmt.Println("> 账号:", acc)
		fmt.Println("> 密码:", pwd)
		fmt.Println("========================")
		goto opt
	default:
		fmt.Println("> 错误的操作类型...")
		goto opt
	}
}
