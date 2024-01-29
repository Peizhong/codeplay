package evaluator

import (
	"io"
	"os"

	"github.com/swaggo/swag"
)

type swagInfo struct{}

func (*swagInfo) ReadDoc() string {
	fs, _ := os.Open("./gen/openapiv2/evaluator/evaluator.swagger.json")
	defer fs.Close()
	bs, _ := io.ReadAll(fs)
	return string(bs)
}

func init() {
	swag.Register("evaluator", &swagInfo{})
}
