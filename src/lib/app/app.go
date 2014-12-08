package app

import (
	"fmt"
	"strings"

	"lib/index"

	"github.com/golang/groupcache"
)

func Handle(ctx groupcache.Context, key string, dst groupcache.Sink) error {
	if ctx == nil {
		return fmt.Errorf("nil Context is invalid")
	}

	// ex) [mtime]-[returntype]-[command]-[key]-[field] like 1417475105-str-HGET-ADINFO-1
	parts := strings.SplitN(key, "-", 5)
	if len(parts) < 4 {
		return fmt.Errorf("given key is invalid")
	}
	rettype := parts[1]
	command := parts[2]
	rkey := parts[3]
	var rfield string
	if len(parts) == 5 {
		rfield = parts[4]
	}

	c := ctx.(*index.Index)
	if got, err := c.Query(command, rettype, rkey, rfield); err != nil {
		return err
	} else {
		dst.SetString(got)
		return nil
	}
}
