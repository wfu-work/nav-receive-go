package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// ToJson 将对象序列化为 JSON 字符串
func ToJson(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ToJsonIndent 将对象序列化为格式化 JSON 字符串（带缩进）
func ToJsonIndent(v interface{}) (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FromJson 将 JSON 字符串反序列化到对象
func FromJson[T any](jsonStr string) (*T, error) {
	var obj T
	err := json.Unmarshal([]byte(jsonStr), &obj)
	if err != nil {
		return nil, err
	}
	return &obj, nil
}

// MustFromJson 与 FromJSON 类似，但失败时 panic
func MustFromJson[T any](jsonStr string) *T {
	obj, err := FromJson[T](jsonStr)
	if err != nil {
		panic(fmt.Sprintf("JSON 解析失败: %v, 内容: %s", err, jsonStr))
	}
	return obj
}

// DeepCopy 通过 JSON 实现深拷贝
func DeepCopy[T any](src T) (*T, error) {
	jsonBytes, err := json.Marshal(src)
	if err != nil {
		return nil, err
	}
	var dst T
	err = json.Unmarshal(jsonBytes, &dst)
	if err != nil {
		return nil, err
	}
	return &dst, nil
}

// PrettyPrint 打印格式化 JSON
func PrettyPrint(v interface{}) error {
	jsonBytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(jsonBytes))
	return nil
}

// ToJsonBuffer 返回 bytes.Buffer，方便 io.Writer 使用
func ToJsonBuffer(v interface{}) (*bytes.Buffer, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(data), nil
}
