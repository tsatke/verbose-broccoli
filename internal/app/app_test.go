package app

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/nettest"
)

func TestAppSuite(t *testing.T) {
	gin.SetMode(gin.TestMode)
	suite.Run(t, new(AppSuite))
}

type AppSuite struct {
	suite.Suite

	app     *App
	cookies *cookiejar.Jar
}

func (suite *AppSuite) SetupTest() {
	suite.cookies, _ = cookiejar.New(nil)

	lis, err := nettest.NewLocalListener("tcp")
	suite.NoError(err)

	suite.app = New(lis)
	suite.IsType(&MemObjectStorage{}, suite.app.objects)
	suite.IsType(&MemDocumentRepo{}, suite.app.documents)
	suite.IsType(&MemAuthService{}, suite.app.auth)

	go func() {
		if err := suite.app.Run(); err != nil {
			panic(err)
		}
	}()
}

func (suite *AppSuite) TearDownTest() {
	if suite.app != nil {
		suite.NoError(suite.app.Close())
	}
}

func (suite *AppSuite) Request(method, path string) TestRequest {
	return TestRequest{
		suite:  suite,
		path:   path,
		method: method,
	}
}

func (suite *AppSuite) login() string {
	user := uuid.New().String()
	pass := uuid.New().String()

	suite.createUser(user, pass)

	suite.
		Request("POST", "/auth/login").
		Body(M{
			"username": user,
			"password": pass,
		}).
		Expect(http.StatusOK, M{
			"success": true,
		})

	return user
}

func (suite *AppSuite) logout() {
	suite.
		Request("GET", "/auth/logout").
		Expect(http.StatusOK, M{
			"success": true,
		})
}

func (suite *AppSuite) createUser(user, pass string) {
	suite.app.auth.(*MemAuthService).data[user] = pass
}

func (suite *AppSuite) createContent(id string, content []byte) {
	suite.NoError(suite.app.objects.(*MemObjectStorage).Create(DocID(id), bytes.NewReader(content)))
}

func (suite *AppSuite) performTestRequest(r TestRequest, wantStatus int, wantResponse M) {
	var buf bytes.Buffer
	if r.rawBody != nil {
		_, _ = buf.Write(r.rawBody)
	} else {
		enc := json.NewEncoder(&buf)
		suite.NoError(enc.Encode(r.body))
	}

	req, err := http.NewRequest(r.method, "http://"+suite.app.listener.Addr().String()+"/rest"+r.path, &buf)
	suite.NoError(err)
	for _, header := range r.header {
		req.Header.Add(header[0], header[1])
	}

	client := &http.Client{
		Jar:     suite.cookies,
		Timeout: 5 * time.Second,
	}

	response, err := client.Do(req)
	suite.NoError(err)
	gotResp, err := io.ReadAll(response.Body)
	suite.NoError(err)
	defer func() {
		_ = response.Body.Close()
	}()

	wantResp, err := json.Marshal(wantResponse)
	suite.NoError(err)

	suite.Equal(wantStatus, response.StatusCode)

	suite.JSONEq(string(wantResp), string(gotResp))
}

func (suite *AppSuite) performTestRequestRaw(r TestRequest, wantStatus int, wantResponse []byte) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	suite.NoError(enc.Encode(r.body))

	req, err := http.NewRequest(r.method, "http://"+suite.app.listener.Addr().String()+"/rest"+r.path, &buf)
	suite.NoError(err)
	for _, header := range r.header {
		req.Header.Add(header[0], header[1])
	}

	client := &http.Client{
		Jar:     suite.cookies,
		Timeout: 5 * time.Second,
	}

	response, err := client.Do(req)
	suite.NoError(err)
	gotResp, err := io.ReadAll(response.Body)
	suite.NoError(err)
	defer func() {
		_ = response.Body.Close()
	}()

	suite.Equal(wantStatus, response.StatusCode)
	suite.Equal(wantResponse, gotResp)
}

type M map[string]interface{}
type Header [2]string

type TestRequest struct {
	suite   *AppSuite
	method  string
	path    string
	header  []Header
	body    M
	rawBody []byte
}

func (r TestRequest) Body(b M) TestRequest {
	r.body = b
	return r
}

func (r TestRequest) RawBody(b []byte) TestRequest {
	r.rawBody = b
	return r
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
		RawBody(buf.Bytes())
}

func (r TestRequest) Header(key, val string) TestRequest {
	r.header = append(r.header, Header{key, val})
	return r
}

func (r TestRequest) Expect(status int, response M) {
	r.suite.performTestRequest(r, status, response)
}

func (r TestRequest) ExpectRaw(status int, response []byte) {
	r.suite.performTestRequestRaw(r, status, response)
}
