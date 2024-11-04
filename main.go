//go:generate trs proto/app.proto
//go:generate trs proto/share.proto
//go:generate trs proto/dmp.proto
//go:generate trs proto/preload.proto
//go:generate protoc -I. --go_out=. --go_opt=paths=source_relative proto/event_hub.proto
package main

import (
	"github.com/lingwei0604/kitty/cmd"
)

func main() {
	cmd.Execute()
}

//main
