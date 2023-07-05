// Package crawlerx
// @Author bcy2007  2023/7/14 16:52
package crawlerx

import "testing"

func TestHeaderRawDataTransfer(t *testing.T) {
	headersData := `
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7
Accept-Encoding: gzip, deflate, br
Accept-Language: zh-CN,zh;q=0.9
Sec-Fetch-Dest: document
Sec-Fetch-Mode: navigate
Sec-Fetch-Site: none
Sec-Fetch-User: ?1
Upgrade-Insecure-Requests: 1
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36
sec-ch-ua: "Not.A/Brand";v="8", "Chromium";v="114", "Google Chrome";v="114"
sec-ch-ua-mobile: ?0
sec-ch-ua-platform: "macOS"
`
	result := headerRawDataTransfer(headersData)
	for _, item := range result {
		t.Logf(`%s %s`, item.Key, item.Value)
	}
}

func TestCookieRawDataTransfer(t *testing.T) {
	anoCookieData := `__jda=76161171.1689152731646522073078.1689152731.1689152731.1689152732.1; unpl=JF8EAKBnNSttURhWAU4AHREUSwhUW1hcQh4HbzUHVlUIS1EFSAFLExB7XlVdXxRKER9vZBRUW1NKUA4ZAisSEXtdU11UC3sSBW9nAVVaXXtUAhgLGCITS21Vbl0PQh8Da2QDVl1fTlMBGAEaFBJKW11uXDhLHwRfVzVTWF9NXQweBisTIEptHzBcRUsQCmdnAVdbWktTABwGGBERTV9VWFQ4SicA; __jdb=76161171.1.1689152731646522073078|1.1689152732; __jdc=76161171; __jdv=76161171|c.duomai.com|t_16282_47115064|jingfen|8b35d37251d144e8851c339a141b2a01|1689152732084; __jdu=1689152731646522073078; areaId=1; ipLoc-djd=1-2800-0-0; PCSYCityID=CN_110000_110100_0; shshshfpa=cdf005e8-5e36-6dfc-3702-70f0cd4e01de-1689152732; shshshfpx=cdf005e8-5e36-6dfc-3702-70f0cd4e01de-1689152732; 3AB9D23F7A4B3CSS=jdd03X6DMZ53N3MBT462GTWXP665CINHUIJFUMV4FHZIPGPAUV3W2EM4EKSCRX5VRO6GZDMTHOEHHFDI6WNIFZF34IEVAXIAAAAMJJFMT5UQAAAAADETJUEYPTCD57IX; _gia_d=1; shshshfpb=dvxly0h_de3L9xNKL8Gjw9Q; 3AB9D23F7A4B3C9B=X6DMZ53N3MBT462GTWXP665CINHUIJFUMV4FHZIPGPAUV3W2EM4EKSCRX5VRO6GZDMTHOEHHFDI6WNIFZF34IEVAXI`
	cookiesData := `_zap=e8c8cf21-2806-49fd-9b06-57636db04ba0; _xsrf=e1680bf4-e2aa-471b-be57-fb4984fa4afe; Hm_lvt_98beee57fd2ef70ccdd5ca52b9740c49=1689326086; Hm_lpvt_98beee57fd2ef70ccdd5ca52b9740c49=1689326086; d_c0=AABaj0X_FBePTib8QuNZhrDGnu4aZbt5Qq0=|1689326086; KLBRSID=57358d62405ef24305120316801fd92a|1689326086|1689326085; captcha_session_v2=2|1:0|10:1689326086|18:captcha_session_v2|88:VDN3K09RcWJkWVNTZ1ZvWnUzNkVnUGUzVHRidkRjTGRPRk9aQ1NkNnpBYW82REdhdGVNVDdRRUF1MnU3aGZjOQ==|d585037e5e28f463d8b7930771725d13e53bc7932634aaa49ac9a37408902c04; SESSIONID=PsQV1hvO3Pp55osLXohfXt9bzTXAtONfLN3ZGE0e3v3`
	result := cookieRawDataTransfer(cookiesData)
	for _, item := range result {
		t.Logf(`%s %s`, item.Name, item.Value)
	}
	anoResult := cookieRawDataTransfer(anoCookieData)
	for _, item := range anoResult {
		t.Logf(`%s %s`, item.Name, item.Value)
	}
}