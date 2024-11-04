package printer

import (
	"flag"
	"testing"

	"github.com/lingwei0604/kitty/rule/dto"
)

var useEtcd bool

func init() {
	flag.BoolVar(&useEtcd, "etcd", false, "use local etcd for testing")
}

func TestPrinter_Sprintf(t *testing.T) {
	if !useEtcd {
		t.Skip("requires etcd")
	}
	p, _ := NewPrinter(dto.Payload{
		PackageName: "com.example.test",
	})
	s := p.Sprintf("kitty.example", 50)
	if s != "用户已获得50积分" {
		t.Fatalf("want %s, got %s", "用户已获得50积分", s)
	}
}
