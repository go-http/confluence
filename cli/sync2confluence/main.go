package main

import (
	"flag"
	"log"

	"confluence"
)

func main() {
	var addr, user, pass, space, dir string

	flag.StringVar(&addr, "addr", "https://www.confluence.com", "Confluence访问地址")
	flag.StringVar(&user, "u", "", "用户名")
	flag.StringVar(&pass, "p", "", "密码")
	flag.StringVar(&space, "s", "", "Confluence空间标识")
	flag.StringVar(&dir, "d", "", "要导入的目录")

	flag.Parse()

	client := confluence.New(addr, user, pass)
	err := client.SpaceContentImportFrom(space, dir)
	if err != nil {
		log.Fatal(err)
	}
}
