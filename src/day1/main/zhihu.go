package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type captureResponse struct {
	show_captcha bool
}

/**
知乎登陆是否需要验证码
*/
func GetIfNeedCapture() bool {
	client := http.Client{}
	var r io.Reader
	req, err := http.NewRequest("GET", "https://www.zhihu.com/api/v3/oauth/captcha?lang=en", r)
	if err != nil {
		log.Printf("make request error...")
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("http request error...")
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	var cr captureResponse
	err = json.Unmarshal(b, &cr)
	if err != nil {
		log.Printf("capture request error...")
	}
	return cr.show_captcha
}

func getCwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic("获取 CWD 失败：" + err.Error())
	}
	return cwd
}

func getCapture() string {
	client := http.Client{}
	var r io.Reader
	req, err := http.NewRequest("PUT", "https://www.zhihu.com/api/v3/oauth/captcha?lang=en", r)
	if err != nil {
		log.Printf("getCapture request error...")
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()

	fileExt := strings.Split(resp.Header.Get("Content-Type"), "/")[1]
	verifyImg := filepath.Join(getCwd(), "verify."+fileExt)
	fd, err := os.OpenFile(verifyImg, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		panic("打开验证码文件失败：" + err.Error())
	}
	defer fd.Close()

	io.Copy(fd, resp.Body) // 保存验证码文件
	openCaptchaFile(verifyImg)
	captcha := readCaptchaInput() // 读取用户输入
	fmt.Println(captcha)
	return captcha
}

/**
发送验证码
*/
func postCapture(capture string) {

}

/**
计算签名
*/
func getSinature(timestamp time.Time) string {
	key := "d1b964811afb40118a12068ff74a12f4"
	a := hmac.New(sha1.New, []byte(key))
	a.Write([]byte("password"))
	a.Write([]byte("c3cef7c66a1843f8b3a9e6a1e3160e20"))
	a.Write([]byte("com.zhihu.web"))
	a.Write([]byte(string(timestamp.UnixNano())))
	signature := a.Sum(nil)
	return string(signature)
}

/**
知乎登陆
*/
func login(timestamp time.Time, signature string, username string, password string, captcha string) {
	var param map[string]string
	param = make(map[string]string)
	param["client_id"] = "c3cef7c66a1843f8b3a9e6a1e3160e20"
	param["grant_type"] = "password"
	param["timestamp"] = string(timestamp.UnixNano())
	param["source"] = "com.zhihu.web"
	param["signature"] = signature
	param["username"] = username
	param["password"] = password
	if captcha != "nil" {
		param["captcha"] = captcha
	}
	param["lang"] = "en"
	data, err := json.Marshal(param)
	if err != nil {
		fmt.Println("json转换失败")
	}
	jsonStr := string(data)
	var r = strings.NewReader(jsonStr)
	client := http.Client{}
	request, err := http.NewRequest("POST", "https://www.zhihu.com/api/v3/oauth/sign_in", r)
	request.Header.Add("User-Agent", `Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Mobile Safari/537.36`)
	request.Header.Add("content-type", `application/x-www-form-urlencoded`)
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("请求失败")
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {

	}
	fmt.Println(string(b))

}

func openCaptchaFile(filename string) error {
	var args []string
	switch runtime.GOOS {
	case "linux":
		args = []string{"xdg-open", filename}
	case "darwin":
		args = []string{"open", filename}
	case "freebsd":
		args = []string{"open", filename}
	case "netbsd":
		args = []string{"open", filename}
	case "windows":
		var (
			cmd      = "url.dll,FileProtocolHandler"
			runDll32 = filepath.Join(os.Getenv("SYSTEMROOT"), "System32", "rundll32.exe")
		)
		args = []string{runDll32, cmd, filename}
	default:
		fmt.Printf("无法确定操作系统，请自行打开验证码 %s 文件，并输入验证码。", filename)
	}

	err := exec.Command(args[0], args[1:]...).Run()
	if err != nil {
		return err
	}

	return nil
}

func readCaptchaInput() string {
	var captcha string
	fmt.Print(color.CyanString("请输入验证码："))
	fmt.Scanf("%s", &captcha)
	return captcha
}
