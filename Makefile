BUILD_TIME := $(shell date "+%F %T")
COMMIT_SHA1 := $(shell git rev-parse HEAD )

# 版本手动指定，由于合约名称在链上有规则，去掉了版本号中 .
VERSION=v100
LDFLAGS := "-s -w -X 'chainweaver.org.cn/chainweaver/trade-guard/notification-contract-go/version.BuildTime=${BUILD_TIME}'  -X 'chainweaver.org.cn/chainweaver/trade-guard/notification-contract-go/version.CommitID=${COMMIT_SHA1}'  -X 'chainweaver.org.cn/chainweaver/trade-guard/notification-contract-go/version.Version=${VERSION}'"

BUILD_NAME := trade-guard-nft-${VERSION}

.PHONY:build

build:
	./build.sh ${BUILD_NAME} ${LDFLAGS}

ut:
	#cd scripts && ./ut_cover.sh
	go test -coverprofile cover.out ./nft.go
	@echo "\n"
	@echo "综合UT覆盖率：" `go tool cover -func=cover.out | tail -1  | grep -P '\d+\.\d+(?=\%)' -o`
	@echo "\n"

lint:
	golangci-lint run ./...

comment:
	gocloc --include-lang=Go --output-type=json --not-match=".*_test.go" . | jq '(.total.comment-.total.files*5)/(.total.code+.total.comment)*100'

pre-commit: lint ut comment

update-mod:
	go get github.com/ethereum/go-ethereum@v1.14.11
	go get chainmaker.org/chainmaker/sdk-go/v2@v2.3.8
	go get chainmaker.org/chainmaker/contract-sdk-go/v2@v2.3.9
	go mod tidy
