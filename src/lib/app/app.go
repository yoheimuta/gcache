package app

import (
	"fmt"
	"strconv"
	"strings"

	"lib/index"

	"github.com/golang/groupcache"
)

func Handle(ctx groupcache.Context, key string, dst groupcache.Sink) error {
	if ctx == nil {
		return fmt.Errorf("nil Context is invalid")
	}

	rettype, command, commandArgs, err := parseKeyString(key)
	if err != nil {
		return err
	}

	c := ctx.(*index.Index)
	if got, err := c.Query(rettype, command, commandArgs); err != nil {
		return err
	} else {
		dst.SetString(got)
		return nil
	}
}

func parseKeyString(key string) (rettype, command string, commandArgs []interface{}, err error) {
	// key is [mtime]-[argc]-[rettype]-[command]-[key]-[field]
	// ex. 1417475105-4-str-HGET-ADINFO-1
	keys := strings.SplitN(key, "-", 3)
	if len(keys) != 3 {
		return "", "", nil, fmt.Errorf("given key is invalid :keys=%v", keys)
	}

	argc, err := strconv.Atoi(keys[1])
	if err != nil {
		return "", "", nil, fmt.Errorf("converting argc to int is failed")
	}
	argv := keys[2]

	parts := strings.SplitN(argv, "-", argc)
	if len(parts) != argc {
		return "", "", nil, fmt.Errorf("given argv is invalid :parts=%v :argc=%v", parts, argc)
	}

	rettype = parts[0]
	command = parts[1]
	commandArgs = convertStrSliceToInterfaceSlice(parts[2:])
	return rettype, command, commandArgs, nil
}

func convertStrSliceToInterfaceSlice(src []string) (dst []interface{}) {
	dst = make([]interface{}, len(src))
	for i, v := range src {
		dst[i] = interface{}(v)
	}
	return dst
}
