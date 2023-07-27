package yakgrpc

import (
	"fmt"
	"github.com/yaklang/yaklang/common/log"
	"github.com/yaklang/yaklang/common/utils"
	"github.com/yaklang/yaklang/common/utils/lowhttp"
	"github.com/yaklang/yaklang/common/yakgrpc/ypb"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestGRPCMUSTPASS_FuzzerSequence(t *testing.T) {
	c, err := NewLocalClient()
	if err != nil {
		panic(err)
	}

	var (
		redirect302done = false
		token           = utils.RandStringBytes(32)
		verified        = false
	)
	host, port := utils.DebugMockHTTPHandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.RequestURI {
		case "/verify":
			if request.Header.Get("Authorization") == "Bearer "+token {
				verified = true
			}

		case "/abc":
			redirect302done = true
			writer.Write([]byte(`{"key": "` + token + `"}`))
			return
		case "/":
			writer.Header().Set("Location", "/abc")
			writer.WriteHeader(302)
			writer.Write([]byte("HELLO HTTP2"))
			return
		}
		writer.Write([]byte("GREAT"))
		return
	})

	client, err := c.HTTPFuzzerSequence(
		utils.TimeoutContextSeconds(10),
		&ypb.FuzzerRequests{Requests: []*ypb.FuzzerRequest{
			{
				Request: string(lowhttp.ReplaceHTTPPacketHeader([]byte(`GET / HTTP/1.1
Host: www.example.com

abc`), "Host", utils.HostPort(host, port))),
				IsHTTPS:                  false,
				PerRequestTimeoutSeconds: 5,
				RedirectTimes:            3,
				Extractors: []*ypb.HTTPResponseExtractor{
					{
						Name:   "test",
						Type:   "json",
						Scope:  "body",
						Groups: []string{".key"},
					},
				},
			},
			{
				Request: string(lowhttp.ReplaceHTTPPacketHeader([]byte(`GET /verify HTTP/1.1
Host: www.example.com
Authorization: Bearer {{params(test)}}

abc`), "Host", utils.HostPort(host, port))),
				IsHTTPS:                  false,
				PerRequestTimeoutSeconds: 5,
				RedirectTimes:            3,
				InheritVariables:         true,
				ForceFuzz:                true,
			},
		}},
	)
	if err != nil {
		panic(err)
	}
	for {
		resp, err := client.Recv()
		if err != nil {
			log.Error(err)
			break
		}
		if resp == nil {
			break
		}
		println(string(resp.Response.RequestRaw))
	}

	if !redirect302done {
		t.Fatal("redirect302done")
	}

	if !verified {
		t.Fatal("verified extractor ")
	}
}

func TestGRPCMUSTPASS_FuzzerSequence_FuzzerWithTag(t *testing.T) {
	c, err := NewLocalClient()
	if err != nil {
		panic(err)
	}
	host, port := utils.DebugMockHTTPHandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(`{"path":` + strconv.Quote(request.URL.Path) + `}`))
		return
	})

	client, err := c.HTTPFuzzerSequence(
		utils.TimeoutContextSeconds(10),
		&ypb.FuzzerRequests{
			Concurrent: 1,
			Requests: []*ypb.FuzzerRequest{
				{
					Request: string(lowhttp.ReplaceHTTPPacketHeader([]byte(`GET /aa={{int(1-10)}} HTTP/1.1
Host: www.example.com

abc`), "Host", utils.HostPort(host, port))),
					IsHTTPS:                  false,
					PerRequestTimeoutSeconds: 5,
					RedirectTimes:            3,
					ForceFuzz:                true,
					Extractors: []*ypb.HTTPResponseExtractor{
						{
							Name:   "test",
							Type:   "json",
							Scope:  "body",
							Groups: []string{`.path`},
						},
					},
				},
				{
					Request: string(lowhttp.ReplaceHTTPPacketHeader([]byte(`GET /verify?a={{param(test)}}/{{int(1-10)}} HTTP/1.1
Host: www.example.com
Authorization: Bearer {{params(test)}}

abc`), "Host", utils.HostPort(host, port))),
					IsHTTPS:                  false,
					PerRequestTimeoutSeconds: 5,
					RedirectTimes:            3,
					InheritVariables:         true,
					ForceFuzz:                true,
				},
			}},
	)
	if err != nil {
		panic(err)
	}
	var count = 0
	for {
		resp, err := client.Recv()
		if err != nil {
			log.Error(err)
			break
		}
		if resp == nil {
			break
		}
		println(resp.Response.GetUrl())
		count++
	}
	if count != 100+10 {
		t.Fatal("Fuzztag COUNT: " + fmt.Sprint(count) + " failed")
	}
}

func TestGRPCMUSTPASS_FuzzerSequence_FuzzerTagWithConcurrent(t *testing.T) {
	c, err := NewLocalClient()
	if err != nil {
		panic(err)
	}
	host, port := utils.DebugMockHTTPHandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		time.Sleep(time.Millisecond * 500)
		writer.Write([]byte(`{"path":` + strconv.Quote(request.URL.Path) + `}`))
		return
	})

	start := time.Now()
	client, err := c.HTTPFuzzerSequence(
		utils.TimeoutContextSeconds(10),
		&ypb.FuzzerRequests{
			Concurrent: 1,
			Requests: []*ypb.FuzzerRequest{
				{
					Request: string(lowhttp.ReplaceHTTPPacketHeader([]byte(`GET /aa={{int(1-10)}} HTTP/1.1
Host: www.example.com

abc`), "Host", utils.HostPort(host, port))),
					IsHTTPS:                  false,
					PerRequestTimeoutSeconds: 5,
					RedirectTimes:            3,
					ForceFuzz:                true,
					Extractors: []*ypb.HTTPResponseExtractor{
						{
							Name:   "test",
							Type:   "json",
							Scope:  "body",
							Groups: []string{`.path`},
						},
					},
				},
				{
					Request: string(lowhttp.ReplaceHTTPPacketHeader([]byte(`GET /verify?a={{param(test)}}/{{int(1-10)}} HTTP/1.1
Host: www.example.com
Authorization: Bearer {{params(test)}}

abc`), "Host", utils.HostPort(host, port))),
					IsHTTPS:                  false,
					PerRequestTimeoutSeconds: 5,
					RedirectTimes:            3,
					InheritVariables:         true,
					ForceFuzz:                true,
				},
			}},
	)
	if err != nil {
		panic(err)
	}
	var count = 0
	for {
		resp, err := client.Recv()
		if err != nil {
			log.Error(err)
			break
		}
		if resp == nil {
			break
		}
		println(resp.Response.GetUrl())
		count++
	}
	if count != 100+10 {
		t.Fatal("Fuzztag COUNT: " + fmt.Sprint(count) + " failed")
	}
	if time.Now().Sub(start).Seconds() <= 5 {
		panic("concurrent(flowmax) is not working")
	}
}

func TestGRPCMUSTPASS_FuzzerSequence_InheritCookie(t *testing.T) {
	c, err := NewLocalClient()
	if err != nil {
		panic(err)
	}

	var (
		redirect302done = false
		token           = utils.RandStringBytes(32)
		verified        = false
	)

	var token2 = utils.RandStringBytes(100)
	host, port := utils.DebugMockHTTPHandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		raw, _ := utils.HttpDumpWithBody(request, true)

		switch request.URL.Path {
		case "/verify":
			if request.Header.Get("Authorization") == "Bearer "+token {
				if lowhttp.GetHTTPPacketCookie(raw, "test") == token2 {
					verified = true
				}
			}

		case "/abc":
			redirect302done = true
			if lowhttp.GetHTTPPacketCookie(raw, "test") == token2 {
				writer.Write([]byte(`{"key": "` + token + `"}`))
			}
			return
		case "/":
			writer.Header().Set("Location", "/abc")
			http.SetCookie(writer, &http.Cookie{
				Name:  "test",
				Value: token2,
			})
			writer.WriteHeader(302)
			writer.Write([]byte("HELLO HTTP2"))
			return
		}
		writer.Write([]byte("GREAT"))
		return
	})

	client, err := c.HTTPFuzzerSequence(
		utils.TimeoutContextSeconds(10),
		&ypb.FuzzerRequests{Requests: []*ypb.FuzzerRequest{
			{
				Request: string(lowhttp.ReplaceHTTPPacketHeader([]byte(`GET / HTTP/1.1
Host: www.example.com

abc`), "Host", utils.HostPort(host, port))),
				IsHTTPS:                  false,
				PerRequestTimeoutSeconds: 5,
				RedirectTimes:            3,
				Extractors: []*ypb.HTTPResponseExtractor{
					{
						Name:   "test",
						Type:   "json",
						Scope:  "body",
						Groups: []string{".key"},
					},
				},
			},
			{
				Request: string(lowhttp.ReplaceHTTPPacketHeader([]byte(`GET /verify HTTP/1.1
Host: www.example.com
Authorization: Bearer {{params(test)}}

abc`), "Host", utils.HostPort(host, port))),
				IsHTTPS:                  false,
				PerRequestTimeoutSeconds: 5,
				RedirectTimes:            3,
				InheritVariables:         true,
				InheritCookies:           true,
				ForceFuzz:                true,
			},
		}},
	)
	if err != nil {
		panic(err)
	}
	for {
		resp, err := client.Recv()
		if err != nil {
			break
		}
		if resp == nil {
			break
		}
		println(string(resp.Response.RequestRaw))
		println(string(resp.Response.ResponseRaw))
		println()
	}

	if !redirect302done {
		t.Fatal("redirect302done")
	}

	if !verified {
		t.Fatal("verified extractor ")
	}
}