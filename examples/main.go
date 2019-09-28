package main

import (
	"fmt"
	"net"
	"net/http"
	"net/http/cookiejar"
	"sync"
	"time"

	"github.com/winterssy/reqwest"
	"golang.org/x/net/publicsuffix"
)

func main() {
	// setParams()
	// setHeaders()
	// setCookies()
	// setFormPayload()
	// setJSONPayload()
	// setFilesPayload()
	// setBasicAuth()
	// setBearerToken()
	// customizeHTTPClient()
	// setProxy()
	// concurrentSafe()
}

func setParams() {
	data, err := reqwest.
		Get("http://httpbin.org/get").
		Params(reqwest.Value{
			"key1": "value1",
			"key2": "value2",
		}).
		Send().
		Text()
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}

func setHeaders() {
	data, err := reqwest.
		Get("http://httpbin.org/get").
		Headers(reqwest.Value{
			"Origin":  "http://httpbin.org",
			"Referer": "http://httpbin.org",
		}).
		Send().
		Text()
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}

func setCookies() {
	data, err := reqwest.
		Get("http://httpbin.org/cookies/set").
		Cookies(
			&http.Cookie{
				Name:  "name1",
				Value: "value1",
			},
			&http.Cookie{
				Name:  "name2",
				Value: "value2",
			},
		).
		Send().
		Text()
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}

func setFormPayload() {
	data, err := reqwest.
		Post("http://httpbin.org/post").
		Form(reqwest.Value{
			"key1": "value1",
			"key2": "value2",
		}).
		Send().
		Text()
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}

func setJSONPayload() {
	data, err := reqwest.
		Post("http://httpbin.org/post").
		JSON(reqwest.Data{
			"msg": "hello world",
			"num": 2019,
		}).
		Send().
		Text()
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}

func setFilesPayload() {
	data, err := reqwest.
		Post("http://httpbin.org/post").
		Files(
			&reqwest.File{
				FieldName: "testimage1",
				FileName:  "testimage1.jpg",
				FilePath:  "./testdata/testimage1.jpg",
			},
			&reqwest.File{
				FieldName: "testimage2",
				FileName:  "testimage2.jpg",
				FilePath:  "./testdata/testimage2.jpg",
			},
		).
		Send().
		Text()
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}

func setBasicAuth() {
	data, err := reqwest.
		Get("http://httpbin.org/basic-auth/admin/pass").
		BasicAuth("admin", "pass").
		Send().
		Text()
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}

func setBearerToken() {
	data, err := reqwest.
		Get("http://httpbin.org/bearer").
		BearerToken("grequests").
		Send().
		Text()
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}

func customizeHTTPClient() {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	redirectPolicy := func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	jar, _ := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	timeout := 120 * time.Second

	req := reqwest.
		WithTransport(transport).
		WithRedirectPolicy(redirectPolicy).
		WithCookieJar(jar).
		WithTimeout(timeout)

	data, err := req.
		Get("http://httpbin.org/get").
		Send().
		Text()
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}

func setProxy() {
	data, err := reqwest.
		WithProxy("http://127.0.0.1:1081").
		Get("http://httpbin.org/get").
		Send().
		Text()
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}

func concurrentSafe() {
	const MaxWorker = 1000
	wg := new(sync.WaitGroup)

	for i := 0; i < MaxWorker; i += 1 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			params := reqwest.Value{}
			params.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))

			data, err := reqwest.
				AcquireLock().
				Get("http://httpbin.org/get").
				Params(params).
				Send().
				Text()
			if err != nil {
				return
			}

			fmt.Println(data)
		}(i)
	}

	wg.Wait()
}
