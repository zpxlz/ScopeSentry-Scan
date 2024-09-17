// utils-------------------------------------
// @file      : utils.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2024/9/6 22:34
// -------------------------------------------

package utils

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

type UtilTools struct{}

var Tools *UtilTools

func InitializeTools() {
	Tools = &UtilTools{}
}

// ReadYAMLFile 读取 YAML 文件并将其解析为目标结构体
func (t *UtilTools) ReadYAMLFile(filePath string, target interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(byteValue, target)
	if err != nil {
		return err
	}

	return nil
}

// WriteYAMLFile 将目标结构体序列化为 YAML 并写入到文件
func (t *UtilTools) WriteYAMLFile(filePath string, data interface{}) error {
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, yamlData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (t *UtilTools) GenerateRandomString(length int) string {
	// 定义字符集
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// 构建随机字符串
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func (t *UtilTools) GetSystemUsage() (int, float64) {
	// 获取CPU使用率
	percent, err := cpu.Percent(3*time.Second, false)
	if err != nil {
		fmt.Println("Failed to get CPU usage:", err)
		return 0, 0
	}
	cpuNum := 0
	if len(percent) > 0 {
		cpuNum = int(percent[0])
	}
	// 获取内存使用率
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("Failed to get memory usage:", err)
		return 0, 0
	}
	return cpuNum, memInfo.UsedPercent
}

func (t *UtilTools) WriteContentFile(filPath string, fileContent string) error {
	// 将字符串写入文件
	return t.WriteByteContentFile(filPath, []byte(fileContent))
}

func (t *UtilTools) WriteByteContentFile(filPath string, fileContent []byte) error {
	// 将字符串写入文件
	if err := ioutil.WriteFile(filPath, fileContent, 0666); err != nil {
		fmt.Printf("Failed to create filPath: %s - %s", filPath, err)
		return err
	}
	return nil
}

// MarshalYAMLToString 将目标结构体序列化为 YAML 字符串
func (t *UtilTools) MarshalYAMLToString(data interface{}) (string, error) {
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(yamlData), nil
}

func (t *UtilTools) StructToJSON(data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// JSONToStruct 将 JSON 字符串反序列化为结构体
func (t *UtilTools) JSONToStruct(jsonStr []byte, result interface{}) error {
	return json.Unmarshal(jsonStr, result)
}

func (t *UtilTools) ParseArgs(args string, keys ...string) (map[string]string, error) {
	// 将参数字符串分割为切片
	argsSlice := strings.Fields(args)

	// 创建一个 FlagSet 对象来解析参数
	fs := flag.NewFlagSet("ParseArgs", flag.ContinueOnError)

	// 创建一个 map 用于存储 flag 的值
	values := make(map[string]*string)
	for _, key := range keys {
		// 初始化 map 并创建对应的 flag
		value := ""
		values[key] = &value
		fs.String(key, "", "a placeholder") // 创建一个 flag 用于存储 key 的值
	}

	// 解析参数
	err := fs.Parse(argsSlice)
	if err != nil {
		return nil, err
	}

	// 获取 key 对应的值并填充到结果 map 中
	result := make(map[string]string)
	for _, key := range keys {
		if valuePtr, ok := values[key]; ok {
			result[key] = *valuePtr
		} else {
			result[key] = ""
		}
	}

	return result, nil
}

func (t *UtilTools) GetParameter(Parameters map[string]map[string]interface{}, module string, plugin string) (string, bool) {
	// 查找 module 是否存在
	if plugins, modOk := Parameters[module]; modOk {
		// 查找 plugin 是否存在
		if param, plugOk := plugins[plugin]; plugOk {
			return param.(string), true
		}
	}
	// 没有找到对应的参数，返回 false
	return "", false
}
