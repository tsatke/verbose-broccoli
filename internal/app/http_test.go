package app

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type M map[string]interface{}

type TestRequest struct {
	suite    *AppSuite
	method   string
	endpoint string
	header   [][2]string
	body     io.Reader
}

func (suite *AppSuite) Get(endpoint string) TestRequest {
	return suite.Request("GET", endpoint)
}

func (suite *AppSuite) Post(endpoint string) TestRequest {
	return suite.Request("POST", endpoint)
}

func (suite *AppSuite) Request(method, endpoint string) TestRequest {
	return TestRequest{
		suite:    suite,
		method:   method,
		endpoint: endpoint,
	}
}

func (r TestRequest) Header(key, value string) TestRequest {
	r.header = append(r.header, [2]string{key, value})
	return r
}

func (r TestRequest) BodyReader(rd io.Reader) TestRequest {
	r.body = rd
	return r
}

func (r TestRequest) Body(data []byte) TestRequest {
	return r.BodyReader(bytes.NewReader(data))
}

func (r TestRequest) BodyJSON(m M) TestRequest {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	r.suite.NoError(enc.Encode(m))
	return r.Body(buf.Bytes())
}

func (r TestRequest) File(field, name string, data []byte) TestRequest {
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	w, err := mw.CreateFormFile(field, name)
	r.suite.NoError(err)

	_, err = w.Write(data)
	r.suite.NoError(err)
	r.suite.NoError(mw.Close())

	return r.
		Header("Content-Type", "multipart/form-data;boundary="+mw.Boundary()).
		Body(buf.Bytes())
}

func (r TestRequest) ExpectRaw(status int, data []byte) {
	r.ExpectCustom(func(res *http.Response) {
		got, err := io.ReadAll(res.Body)
		r.suite.NoError(err)
		r.suite.NoError(res.Body.Close())

		r.suite.Equal(status, res.StatusCode)
		r.suite.Equal(data, got)
	})
}

func (r TestRequest) ExpectJSON(status int, m M) {
	r.ExpectCustom(func(res *http.Response) {
		data, err := json.Marshal(m)
		r.suite.NoError(err)

		got, err := io.ReadAll(res.Body)
		r.suite.NoError(err)
		r.suite.NoError(res.Body.Close())

		r.suite.Equal(status, res.StatusCode)
		r.suite.JSONEq(string(data), string(got))
	})
}

func (r TestRequest) ExpectCustom(validate func(*http.Response)) {
	c := &http.Client{
		Jar:     r.suite.cookies,
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest(r.method, "http://"+r.suite.app.listener.Addr().String()+"/rest"+r.endpoint, r.body)
	r.suite.NoError(err)
	for _, header := range r.header {
		req.Header.Add(header[0], header[1])
	}

	res, err := c.Do(req)
	r.suite.NoError(err)

	validate(res)
}
