package tools

import (
	"fmt"
	"github.com/yaklang/yaklang/common/fp"
	"testing"
)

func Test_scanFingerprint(t *testing.T) {
	//target := "150.129.109.26"
	target := "127.0.0.1"
	//target := "118.171.54.61"
	//target := "192.168.3.113"
	//target := "117.212.17.42"
	//target = "120.96.38.58"

	//port := "3307"
	//port := "21"
	port := "59380"

	//protoList := []interface{}{"tcp", "udp"}
	//protoList := []interface{}{"tcp"}
	protoList := []interface{}{"udp"}

	pp := func(proto ...interface{}) fp.ConfigOption {
		return fp.WithTransportProtos(fp.ParseStringToProto(proto...)...)
	}

	ch, err := scanFingerprint(target, port, pp(protoList...), fp.WithProbeTimeoutHumanRead(5))
	//ch, err := scanFingerprint(target, "162", pp(protoList...), fp.WithProbeTimeoutHumanRead(5))

	if err != nil {
		t.Error(err)
	}

	for v := range ch {
		fmt.Println(v.String())
	}
}
