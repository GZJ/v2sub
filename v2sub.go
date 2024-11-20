package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"flag"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

type Vmess struct {
	Ps   string `json:"ps"`
	Port string `json:"port"`
	Add  string `json:"add"`
	Id   string `json:"id"`
	Aid  int `json:"aid"`
	net  string `json:"net"`
	Type string `json:"type"`
	Tls  string `json:"tls"`
}

func main() {
	var (
		name             string
		subUrl           string
		proxyPath        string
		tmplFile         string
		configOutputPath string
	)

	flag.StringVar(&name, "name", "", "Name of the subscription address.")
	flag.StringVar(&subUrl, "sub-url", "", "URL of the subscription address.")
	flag.StringVar(&proxyPath, "proxy-path", "", "Path to the proxy for accessing the subscription address.")
	flag.StringVar(&tmplFile, "tmpl-file", "config.tmpl", "Path to the template file.")
	flag.StringVar(&configOutputPath, "config-output-path", ".", "Base path for generated configuration files.")
	flag.Parse()

	enableProxy(proxyPath)
	generate(name, tmplFile, configOutputPath, subStr(subUrl))
}

func enableProxy(proxyPath string) {
	if proxyPath != "" {
		proxyUrl, err := url.Parse(proxyPath)
		if err != nil {
			log.Fatalf("Error parsing proxy URL %s: %v", proxyPath, err)
		}
		http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
		log.Println("Proxy is enabled.")
	} else {
		http.DefaultTransport = &http.Transport{Proxy: nil} 
		log.Println("Proxy is not enabled.")
	}
}

func subStr(subUrl string) []string {
	resp, err := http.Get(subUrl)
	if err != nil {
		log.Fatalf("Error getting subscription URL %s: %v", subUrl, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	strRaw := preprocessB64Str(string(body))
	strDecoded := strings.Split(strRaw, "\n")
	return strDecoded
}

func preprocessB64Str(strB64 string) string {
	var (
		bstr       []byte
		strDecoded string
		err        error
	)

	bstr, err = b64.URLEncoding.DecodeString(strB64)
	if err != nil {
		bstr, err = b64.StdEncoding.WithPadding(b64.NoPadding).DecodeString(strB64)
		if err != nil {
			bstr, err = b64.StdEncoding.DecodeString(strB64)
			if err != nil {
				log.Printf("Error decoding Base64: %v", err)
			}
		}
	}

	strDecoded = string(bstr)
	return strDecoded
}

func generate(name string, tmplFile string, configOutputPath string, subUrls []string) {
	var (
		strB64   string
		vmess    Vmess
		fileJson *os.File
	)

	err := os.MkdirAll(configOutputPath, os.ModePerm)
	if err != nil {
		log.Println("Error creating config output directory:", err)
		return
	}

	pathEncoded := path.Join(configOutputPath, name+"_encoded.txt")
	fdEncoded, err := os.OpenFile(pathEncoded, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Println("Error opening encoded file:", err)
	}
	defer fdEncoded.Close()

	pathDecoded := path.Join(configOutputPath, name+"_decoded.txt")
	fdDecoded, err := os.OpenFile(pathDecoded, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Println("Error opening decoded file:", err)
	}
	defer fdDecoded.Close()

	tmpl, err := template.ParseFiles(tmplFile)
	if err != nil {
		log.Println("Error parsing template file:", err)
	}
	re, err := regexp.Compile("vmess://*.")
	if err != nil {
		log.Println("Error compiling regular expression:", err)
	}

	for i, line := range subUrls {
		if re.MatchString(line) {
			strB64 = strings.Replace(line, "vmess://", "", -1)

			if _, err = fdEncoded.WriteString(strB64 + "\n"); err != nil {
				log.Println("Error writing encoded data:", err)
			}
			str := preprocessB64Str(strB64)
			if _, err = fdDecoded.WriteString(str + "\n"); err != nil {
				log.Println("Error writing decoded data:", err)
			}

			vmess = Vmess{}
			err = json.Unmarshal([]byte(str), &vmess)
			if err != nil {
				log.Println("Error unmarshaling JSON:", err)
			}

			pathConfig := path.Join(configOutputPath, name+strconv.Itoa(i+1)+".json")

			fileJson, err = os.Create(pathConfig)
			if err != nil {
				log.Println("Error creating JSON file:", err)
			}
			defer fileJson.Close()
			tmpl.Execute(fileJson, vmess)
		}
	}
}
