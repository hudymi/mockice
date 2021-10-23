package endpoint_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/hudymi/mockice/pkg/endpoint"
)

func TestDefaultConfig(t *testing.T) {
	// Given
	g := NewGomegaWithT(t)

	// When
	cfg := endpoint.DefaultConfig()

	// Then
	g.Expect(cfg).ToNot(BeEmpty())
	g.Expect(cfg).To(HaveLen(1))
	g.Expect(cfg[0].Name).To(Equal("hello"))
	g.Expect(cfg[0].DefaultResponseContent).ToNot(BeZero())
}

func TestEndpoint_Handle_DefaultConfig(t *testing.T) {
	cfg := endpoint.DefaultConfig()[0]

	for _, method := range []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodTrace} {
		t.Run(fmt.Sprintf("Method %s", method), func(t *testing.T) {
			// Given
			g := NewGomegaWithT(t)

			end := endpoint.New(cfg)
			recorder := httptest.NewRecorder()
			handler := http.HandlerFunc(end.Handle)
			request := httptest.NewRequest(method, fmt.Sprintf("/%s", cfg.Name), strings.NewReader("test"))

			// When
			handler.ServeHTTP(recorder, request)
			defer recorder.Result().Body.Close()

			// Then
			g.Expect(recorder.Result().StatusCode).To(Equal(http.StatusOK))
			g.Expect(recorder.Result().Header.Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))

			buff := new(bytes.Buffer)
			buff.ReadFrom(recorder.Result().Body)
			g.Expect(buff.String()).To(Equal(cfg.DefaultResponseContent))
		})
	}
}

func TestEndpoint_Handle(t *testing.T) {
	for testName, testCase := range map[string]struct {
		method               string
		expectedResponseCode int
		config               endpoint.Config
	}{
		"Simple": {
			method:               http.MethodGet,
			expectedResponseCode: http.StatusOK,
			config: endpoint.Config{
				Name:                   "simple",
				DefaultResponseContent: "test",
			},
		},
		"With methods - OK": {
			method:               http.MethodTrace,
			expectedResponseCode: http.StatusOK,
			config: endpoint.Config{
				Name:                   "simple",
				DefaultResponseContent: "test",
				Methods:                []string{http.MethodGet, http.MethodPost, http.MethodTrace},
			},
		},
		"With methods - Invalid": {
			method:               http.MethodDelete,
			expectedResponseCode: http.StatusMethodNotAllowed,
			config: endpoint.Config{
				Name:                   "simple",
				DefaultResponseContent: "test",
				Methods:                []string{http.MethodGet, http.MethodPost, http.MethodTrace},
			},
		},
		"With content-type": {
			method:               http.MethodGet,
			expectedResponseCode: http.StatusOK,
			config: endpoint.Config{
				Name:                       "simple",
				DefaultResponseContent:     "test",
				DefaultResponseContentType: stringPointer("text/markdown; charset=utf-8"),
			},
		},
		"With response code": {
			method:               http.MethodGet,
			expectedResponseCode: http.StatusNotModified,
			config: endpoint.Config{
				Name:                   "simple",
				DefaultResponseContent: "test",
				DefaultResponseCode:    intPointer(http.StatusNotModified),
			},
		},
		"With file": {
			method:               http.MethodGet,
			expectedResponseCode: http.StatusOK,
			config: endpoint.Config{
				Name:                       "simple",
				DefaultResponseCode:        intPointer(http.StatusOK),
				DefaultResponseContentType: stringPointer("text/markdown; charset=utf-8"),
				DefaultResponseFile:        stringPointer("testdata/test.md"),
			},
		},
		"With file and content": {
			method:               http.MethodGet,
			expectedResponseCode: http.StatusOK,
			config: endpoint.Config{
				Name:                       "simple",
				DefaultResponseContent:     "test",
				DefaultResponseCode:        intPointer(http.StatusOK),
				DefaultResponseContentType: stringPointer("text/markdown; charset=utf-8"),
				DefaultResponseFile:        stringPointer("testdata/test.md"),
			},
		},
		"With not existing file": {
			method:               http.MethodGet,
			expectedResponseCode: http.StatusInternalServerError,
			config: endpoint.Config{
				Name:                       "simple",
				DefaultResponseCode:        intPointer(http.StatusOK),
				DefaultResponseContentType: stringPointer("text/markdown; charset=utf-8"),
				DefaultResponseFile:        stringPointer("testdata/test.md.fail"),
			},
		},
	} {
		t.Run(testName, func(t *testing.T) {
			// Given
			g := NewGomegaWithT(t)

			end := endpoint.New(testCase.config)
			recorder := httptest.NewRecorder()
			handler := http.HandlerFunc(end.Handle)
			request := httptest.NewRequest(testCase.method, fmt.Sprintf("/%s", testCase.config.Name), strings.NewReader("test"))

			// When
			handler.ServeHTTP(recorder, request)
			defer recorder.Result().Body.Close()

			// Then
			g.Expect(recorder.Result().StatusCode).To(Equal(testCase.expectedResponseCode))
			if recorder.Result().StatusCode != http.StatusOK && !(testCase.config.DefaultResponseCode != nil && recorder.Result().StatusCode == *testCase.config.DefaultResponseCode) {
				return
			}

			if testCase.config.DefaultResponseContentType != nil {
				g.Expect(recorder.Result().Header.Get("Content-Type")).To(Equal(*testCase.config.DefaultResponseContentType))
			} else {
				g.Expect(recorder.Result().Header.Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
			}

			buff := new(bytes.Buffer)
			buff.ReadFrom(recorder.Result().Body)
			if testCase.config.DefaultResponseFile == nil {
				g.Expect(buff.String()).To(Equal(testCase.config.DefaultResponseContent))
			} else {
				content, err := readFile(*testCase.config.DefaultResponseFile)
				g.Expect(err).To(Succeed())
				g.Expect(buff.String()).To(Equal(content))
			}
		})
	}
}

func readFile(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func stringPointer(value string) *string {
	return &value
}

func intPointer(value int) *int {
	return &value
}

func TestEndpoint_Name(t *testing.T) {
	// Given
	g := NewGomegaWithT(t)
	end := endpoint.New(endpoint.Config{Name: "test"})

	// When
	name := end.Name()

	// Then
	g.Expect(name).To(Equal("test"))
}
