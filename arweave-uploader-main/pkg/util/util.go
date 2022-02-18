package util

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

func RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

func HTTPGet(url string) (data map[string]interface{}, err error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/json")
	transport := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Printf("[error] client.Do err:%s", err)
		return nil, err
	}
	defer resp.Body.Close()

	if http.StatusOK != resp.StatusCode {
		return nil, errors.Wrap(err, "StatusCode error")
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		log.Printf("[error] ioutil.ReadAll err:%s", err)
		return nil, errors.Wrap(err, "ioutil.ReadAll error")
	}

	var v interface{}
	err = jsoniter.Unmarshal(buf, &v)
	if nil != err {
		log.Printf("[error] jsoniter.Unmarshal err:%s", err)
		return nil, errors.Wrap(err, "Unmarshal error")
	}

	data = make(map[string]interface{}, 0)
	data["metadata"] = v

	return data, nil
}

func IsExistFromSlice(source []string, value string) bool {
	for _, item := range source {
		if item == value {
			return true
		}
	}
	return false
}
