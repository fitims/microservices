package infrastructure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

// ServiceClient is a struct that contains details about the client that uses the service
type ServiceClient struct {
	httpClient *http.Client
	url        string
}

type requestBuilder func() (*http.Request, error)

type reqHeader struct {
	Key   string
	Value string
}

func newHeader(key, value string) reqHeader {
	return reqHeader{
		Key:   key,
		Value: value,
	}
}

// NewServiceClient return a new instance of a service client used to access services
func NewServiceClient(rootURL string) *ServiceClient {
	return &ServiceClient{
		url: rootURL,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

// GET is a method of ServiceClient and wraps HTTP GET functionality
func (svc *ServiceClient) GET(route string, queryParams ...string) ([]byte, error) {
	fmt.Println("Doing a client GET to : ", buildURL(svc.url, route, queryParams))

	builder := func() (*http.Request, error) {
		return http.NewRequest("GET", buildURL(svc.url, route, queryParams), nil)
	}

	return svc.doJSON(builder)
}

// POST is a method of ServiceClient and wraps HTTP POST functionality
func (svc *ServiceClient) POST(route string, data interface{}) ([]byte, error) {
	builder := func() (*http.Request, error) {
		buf, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		url := buildURL(svc.url, route, []string{})

		fmt.Println("\n\nDoing Post URL  : ", url)
		fmt.Println("\n\nMarshalled data : ", string(buf))

		reader := bytes.NewReader(buf)
		return http.NewRequest("POST", buildURL(svc.url, route, []string{}), reader)
	}

	return svc.doJSON(builder)
}

// PUT is a method of ServiceClient and wraps HTTP PUT functionality
func (svc *ServiceClient) PUT(route string, data interface{}) ([]byte, error) {
	builder := func() (*http.Request, error) {
		buf, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		reader := bytes.NewReader(buf)
		return http.NewRequest("PUT", buildURL(svc.url, route, []string{}), reader)
	}

	return svc.doJSON(builder)
}

// DELETE is a method of ServiceClient and wraps HTTP DELETE functionality
func (svc *ServiceClient) DELETE(route string, queryParams ...string) ([]byte, error) {
	builder := func() (*http.Request, error) {
		log.Println("Doing a client DELETE to : ", svc.url+route)
		return http.NewRequest("DELETE", buildURL(svc.url, route, queryParams), nil)
	}

	return svc.doJSON(builder)
}

// FILE is a method of ServiceClient and wraps HTTP multipart POST functionality
// to upload a file to the client
func (svc *ServiceClient) FILE(route string, filename string, rdr io.Reader) ([]byte, error) {
	builder := func() (*http.Request, error) {
		reader, header, err := getFileUploadHTTPBody(filename, rdr)
		if err != nil {
			return nil, err
		}
		req, err := http.NewRequest("POST", buildURL(svc.url, route, []string{}), reader)
		if err != nil {
			return nil, err
		}
		req.Header.Add("Content-Type", header)
		return req, nil
	}
	return svc.doFile(builder)
}

func (svc *ServiceClient) doFile(builder requestBuilder) ([]byte, error) {
	return svc.do(builder)
}

func (svc *ServiceClient) doJSON(builder requestBuilder) ([]byte, error) {
	return svc.do(builder, newHeader("Content-Type", "application/json"))
}

func (svc *ServiceClient) do(builder requestBuilder, headers ...reqHeader) ([]byte, error) {
	req, err := builder()
	if err != nil {
		return nil, err
	}

	if len(headers) > 0 {
		for _, v := range headers {
			req.Header.Add(v.Key, v.Value)
		}
	}

	resp, err := svc.httpClient.Do(req)
	if err != nil {
		log.Println("Error passing data to service : ", err)
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func getFileUploadHTTPBody(filename string, rdr io.Reader) (*bytes.Buffer, string, error) {
	fileContents, err := ioutil.ReadAll(rdr)
	if err != nil {
		return nil, "", err
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, "", err
	}
	part.Write(fileContents)

	err = writer.Close()
	if err != nil {
		return nil, "", err
	}

	return body, writer.FormDataContentType(), nil
}

func normalizeURL(url string) string {
	if strings.HasSuffix(url, "/") {
		return url
	}
	return url + "/"
}

func normalizeRoute(route string) string {
	normalized := route
	if strings.HasSuffix(route, "/") {
		normalized = strings.TrimSuffix(route, "/")
	}
	if strings.HasPrefix(route, "/") {
		normalized = strings.TrimPrefix(route, "/")
	}
	return normalized
}

func buildURL(baseURL string, route string, queryParams []string) string {
	if len(route) < 1 && len(queryParams) < 1 {
		return baseURL
	}

	q := strings.Join(queryParams, "&")
	u := normalizeURL(baseURL)
	r := normalizeRoute(route)

	if len(q) > 0 {
		q = "?" + q
	}
	return u + r + q
}
