include $(GOROOT)/src/Make.inc

TARG=mango
GOFMT=gofmt

GOFILES=\
				mango.go\

include $(GOROOT)/src/Make.pkg

format:
	${GOFMT} -w ${GOFILES}
	${GOFMT} -w mango_test.go
	${GOFMT} -w examples/hello.go
	${GOFMT} -w examples/logger.go
