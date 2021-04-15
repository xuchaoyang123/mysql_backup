//package config
//
//import (
//	"check_list/logfile"
//	"encoding/json"
//	"flag"
//	"io/ioutil"
//	"os"
//	"bufio"
//	"io"
//	"strings"
//)
//
//var (
//	ConfMap  map[string]map[string]string
//	SoarPath string
//)
//
//func init() {
//	var (
//		//savePath string
//		config string
//		file   *os.File
//		bytes  []byte
//		err    error
//	)
//
//	//flag.StringVar(&config, "config", "./conf.json", "指定配置文件")
//	flag.StringVar(&config, "config", "/data/goland/project/src/check_list/config/conf.json", "指定配置文件")
//	flag.Parse()
//
//	if file, err = os.Open(config); err != nil {
//		logfile.Loger.Println("[error] :配置文件找不到,err= ", err) //报错追加日志.
//
//	}
//
//	if bytes, err = ioutil.ReadAll(file); err != nil {
//		logfile.Loger.Println("[error] :读取配置文件错误,err= ", err)
//	}
//
//	//fmt.Println(ConfMap)
//	ConfMap = make(map[string]map[string]string)
//	if err = json.Unmarshal(bytes, &ConfMap); err != nil {
//		logfile.Loger.Println("[error] :解析配置文件错误,err= ", err)
//	}
//
//}

package config

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func InitConfig() map[string]string {

	var (
		//savePath string
		Config string
		err    error
		file   *os.File
	)

	//flag.StringVar(&Config, "config", "./conf.txt", "指定配置文件")
	flag.StringVar(&Config, "config", "/data/goland/project/src/check_list/config/conf.txt", "指定配置文件")
	flag.Parse()

	if file, err = os.Open(Config); err != nil {
		fmt.Println("[error] :配置文件找不到,err= ", err) //报错追加日志.

	}

	if _, err = ioutil.ReadAll(file); err != nil {
		fmt.Println("[error] :读取配置文件错误,err= ", err)
	}

	//初始化
	myMap := make(map[string]string)

	//打开文件指定目录，返回一个文件f和错误信息
	f, err := os.Open(Config)
	defer f.Close()

	//异常处理 以及确保函数结尾关闭文件流
	if err != nil {
		panic(err)
	}

	//创建一个输出流向该文件的缓冲流*Reader
	r := bufio.NewReader(f)
	for {
		//读取，返回[]byte 单行切片给b
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		//去除单行属性两端的空格
		s := strings.TrimSpace(string(b))
		//fmt.Println(s)

		//判断等号=在该行的位置
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}
		//取得等号左边的key值，判断是否为空
		key := strings.TrimSpace(s[:index])
		if len(key) == 0 {
			continue
		}

		//取得等号右边的value值，判断是否为空
		value := strings.TrimSpace(s[index+1:])
		if len(value) == 0 {
			continue
		}
		//这样就成功吧配置文件里的属性key=value对，成功载入到内存中c对象里
		myMap[key] = value
	}
	return myMap
}
