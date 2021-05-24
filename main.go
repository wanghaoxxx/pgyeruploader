package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

const ReleasePath = "build/app/outputs/apk/release/app-release.apk"
const DebugPath = "build/app/outputs/apk/debug/app-debug.apk"

type UploadResp struct {
	Code int                    `json:"code:code"`
	Data map[string]interface{} `json:"data"`
}

func main() {
	args := os.Args
	var debug = true
	if len(args) > 2 { // 处理参数
		if args[1] == "-t" {
			if args[2] == "release" {
				debug = false
			}
		}
	}
	var uploadFile = DebugPath
	if !debug {
		uploadFile = ReleasePath
	}

	fmt.Println("apk build start ...")
	cmd := exec.Command("flutter", "--no-color", "build", "apk")
	if debug {
		cmd.Args = append(cmd.Args, "--debug")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		fmt.Println("cmd start error", err)
		return
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Println("cmd Wait error", err)
		return
	}
	fmt.Println("apk build end ...")

	fmt.Println("uploading...")
	var params = make(map[string]string)
	params["_api_key"] = "__api_key--" // 改成你自己的api key
	params["buildInstallType"] = "2"
	params["buildPassword"] = "1"
	r, err := uploadApk(params, uploadFile)
	if err != nil || r == nil {
		fmt.Println("uploadApk error", err)
		return
	}
	client := http.DefaultClient
	res, err := client.Do(r)
	if err != nil {
		fmt.Println("uploadApk error", err)
		return
	}

	defer res.Body.Close()
	resp, err := io.ReadAll(res.Body)
	rst := new(UploadResp)
	_ = json.Unmarshal(resp, rst)
	if rst.Code != 0 {
		fmt.Println("uploadApk parse json error", err)
		return
	}
	//应用短连接
	url := "https://www.pgyer.com/" + rst.Data["buildShortcutUrl"].(string)
	fmt.Println("upload successfully url ：" + url)
}

func uploadApk(params map[string]string, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// 实例化multipart
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	// 创建multipart 文件字段
	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		return nil, err
	}
	// 写入文件数据到multipart
	_, err = io.Copy(part, file)
	//将额外参数也写入到multipart
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	//创建请求
	req, err := http.NewRequest("POST", "https://www.pgyer.com/apiv2/app/upload", body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	return req, nil
}
