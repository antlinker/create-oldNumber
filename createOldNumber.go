package createoldnumber

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	globalLock sync.RWMutex // 全局锁
	ld         *lastData    // 上次生成单号的信息
)

const filePath = "./createoldnumber.txt" // 储存文件的路径

// 上次生成单号的信息
type lastData struct {
	Date string // 上次生成单号的日期
	Num  int    // 上次生成单号的数字
}

// Init 初始化
func Init() {

	// 判断储存文件是否存在
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			// 文件不存在则创建，其他错误则报错
			file, err := os.Create(filePath)
			if err != nil {
				fmt.Println("\n获取单号初始化失败1：", err)
			}
			file.Close()
		} else {
			fmt.Println("\n获取单号初始化失败2：", err)
		}
	}

	// 获取文件储存的内容
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("\n获取单号初始化失败3：", err)
	}

	ld = &lastData{}
	if string(b) == "" {
		// 如果文件没有内容，默认设置为今天的第0单
		ld.Date = time.Now().Format("060102")
		ld.Num = 1
	} else {
		// 如果文件有内容，解析文件的内容
		err = json.Unmarshal(b, ld)
		if err != nil {
			fmt.Println("\n获取单号初始化失败4：", err)
		}
	}
}

// GetNum 获取单号
func GetNum() (num string, err error) {
	globalLock.Lock()
	defer globalLock.Unlock()

	if ld.Date == time.Now().Format("060102") {
		// 如果上次获取是在今天，单号数字加1
		ld.Num++
	} else {
		// 如果上次获取不是今天
		ld.Date = time.Now().Format("060102")
		ld.Num = 1
	}

	// 单号 (后面不足三位补零)
	if ld.Num < 10 {
		num = ld.Date + "00" + strconv.Itoa(ld.Num)
	} else if ld.Num < 100 {
		num = ld.Date + "0" + strconv.Itoa(ld.Num)
	} else {
		num = ld.Date + strconv.Itoa(ld.Num)
	}

	// 把这次获取单号的信息转成JSON字符串
	jsonstr, err := json.Marshal(ld)
	if err != nil {
		fmt.Println("\n获取单号失败：", err)
		return
	}

	// 把JSON字符串写入文件储存
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Println("\n获取单号失败：", err)
		return
	}
	file.WriteString(string(jsonstr))
	file.Close()

	return
}
