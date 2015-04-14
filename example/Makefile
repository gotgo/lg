# Makefile

# Build Data
IS_DIRTY=$(test -n "`git status --porcelain`" && echo "+CHANGES" || true)
COMMIT_HASH=`git rev-parse --short HEAD 2>/dev/null`
BUILD_DATE=`date +%FT%T%z`

# Compliler Flags
LDFLAGS=-ldflags "-X main.CommitHash ${COMMIT_HASH}${IS_DIRTY} -X main.BuildDate ${BUILD_DATE}"

all: build

build-info:
	echo ${COMMIT_HASH}${IS_DIRTY}
	echo ${BUILD_DATE}

build:
	${GOROOT}/bin/go build ${LDFLAGS} 

install:
	${GOROOT}/bin/go install -race ${LDFLAGS} 
