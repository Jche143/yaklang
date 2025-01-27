package yakgrpc

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/yaklang/yaklang/common/utils"
	_ "github.com/yaklang/yaklang/common/yak"
	"github.com/yaklang/yaklang/common/yakgrpc/ypb"
	"os"
	"testing"
)

func TestGRPCMUSTPASS_GetYakitCompletionRaw(t *testing.T) {
	test := assert.New(t)

	c, err := NewLocalClient()
	if err != nil {
		test.FailNow(err.Error())
	}

	rsp, err := c.GetYakitCompletionRaw(context.Background(), &ypb.Empty{})
	if err != nil {
		test.FailNow(err.Error())
	}

	if len(rsp.RawJson) <= 0 {
		test.FailNow("empty result")
	}

	if !utils.MatchAllOfSubString(
		string(rsp.RawJson),
		"QueryUrlsAll() chan string",
		"O_RDWR: int = 0x2",
	) {
		panic("generate completion failed")
	}
}

func TestServer_GetYakitCompletionRaw_Antlr4Yak(t *testing.T) {
	test := assert.New(t)

	os.Setenv("YAKMODE", "vm")

	c, err := NewLocalClient()
	if err != nil {
		test.FailNow(err.Error())
	}

	rsp2, err := c.GetYakVMBuildInMethodCompletion(context.Background(), &ypb.GetYakVMBuildInMethodCompletionRequest{})
	if err != nil {
		panic(err)
	}

	if len(rsp2.GetSuggestions()) <= 0 {
		panic(1)
	}

	rsp, err := c.GetYakitCompletionRaw(context.Background(), &ypb.Empty{})
	if err != nil {
		test.FailNow(err.Error())
	}

	if len(rsp.RawJson) <= 0 {
		test.FailNow("empty result")
	}
	spew.Dump(len(rsp.RawJson))

}
