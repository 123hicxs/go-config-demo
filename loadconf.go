package loadconf

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// 导入此包时执行InitConfig函数
func init() {
	InitConfig()
}

var fileconfig map[string]interface{}

func ReadConfig(e string, defval string) string {
	//首先从环境变量中读取REC_CONFIG变量，如果变量不存在则从文件配置中读取
	e_val := os.Getenv(e)
	if e_val == "" {
		e_val = readenvFromfile(e)
		if e_val == "" {
			e_val = defval
		}
	}
	return e_val
}

func readenvFromfile(e string) string {

	// 判断e是否在fileconfig中，如果在的话返回值，不存在的话返回空字符串
	if v, ok := fileconfig[e]; ok {
		return v.(string)
	} else {
		return ""
	}

}

// 初始化file_config
func InitConfig() {

	fileconfig = make(map[string]interface{})
	// 从环境变量读取配置文件路径，没有采取默认值
	globalConfPath := os.Getenv("GLOBAL_CONF_PATH")
	if globalConfPath == "" {
		globalConfPath = "./global_conf"
	}

	file, err := os.Open(globalConfPath)
	defer file.Close()

	if err != nil {
		fmt.Printf("[error] %v \n", err)
		return
	}
	reloadConfig(file)

}

func reloadConfig(file *os.File) {

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		// 根据每行的数据进行分割，分割符m默认为"="或者":"
		line := scanner.Text()
		splited := strings.Split(line, "=")
		if len(splited) != 2 {
			splited = strings.Split(line, ":")
			if len(splited) != 2 {
				fmt.Println("[error] invalid config line:", line)
				continue
			}
		}
		fileconfig[strings.TrimSpace(splited[0])] = strings.TrimSpace(splited[1])
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}
