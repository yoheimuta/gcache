DEBUG_FLAG = $(if $(DEBUG),-v)
GOPATH_ENV="$(PWD)/.godeps:$(PWD)"
GOBIN_ENV="$(PWD)/.godeps/bin"

deps:
	wget -qO- https://raw.githubusercontent.com/pote/gpm/v1.2.3/bin/gpm | GOPATH=$(GOPATH_ENV) bash

test: deps
	sh src/test/setup_subnet.sh
	GOPATH=$(GOPATH_ENV) sh src/test/gotest.sh $(DEBUG_FLAG)

install: deps
	GOBIN=$(GOBIN_ENV) GOPATH=$(GOPATH_ENV) go install gcache
