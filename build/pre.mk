pre:
	@echo "=== TSAK.clips === [Preinstallation of some stuff]"
	go get github.com/client9/misspell/cmd/misspell
	go get golang.org/x/tools/cmd/godoc
	go get github.com/llorllale/go-gitlint/cmd/go-gitlint
	go get github.com/psampaz/go-mod-outdated
	go get golang.org/x/tools/cmd/goimports
	go get gopkg.in/alecthomas/kingpin.v2
	go get gotest.tools/gotestsum
	go get github.com/stretchr/testify/assert@v1.6.1
	go get github.com/goombaio/namegenerator
	go get github.com/rs/xid
	go get github.com/glycerine/zygomys/zygo
	go get github.com/lrita/cmap
	go get github.com/edwingeng/deque
	go get github.com/gosnmp/gosnmp
	go mod download github.com/stretchr/testify
