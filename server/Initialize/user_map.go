package Initialize

import (
	"bufio"
	"fmt"
	"im/global"
	"io"
	"os"
	"strings"
)

func InitUserMap() {
	f, err := os.Open("./user.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	r := bufio.NewReader(f)
	//i := 0
	for {
		//fmt.Println(i)
		str, err := r.ReadString('\n')
		if err == io.EOF {
			fmt.Println(err)
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}
		flagAndData := strings.Split(str, ":")
		//fmt.Println("u--------",flagAndData)
		if len(flagAndData) != 2 {
			continue
		}
		userName := strings.TrimRight(flagAndData[1], "\n")
		userName = strings.TrimRight(userName, "\r")
		if flagAndData[0] == "u" {
			str, err := r.ReadString('\n')

			flagAndData := strings.Split(str, ":")
			//fmt.Println("p-------",flagAndData)
			if len(flagAndData) != 2 {
				continue
			}
			var psw = strings.TrimRight(flagAndData[1], "\n")
			psw = strings.TrimRight(psw, "\r")
			if flagAndData[0] == "p" {
				global.UserMap.Store(userName, psw)
				if err == io.EOF {
					fmt.Println(err)
					break
				}
				if err != nil {
					fmt.Println(err)
					break
				}
			} else {
				continue
			}
		} else {
			continue
		}

	}
}
