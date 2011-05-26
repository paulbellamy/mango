include $(GOROOT)/src/Make.inc

TARG=mango
GOFMT=gofmt

GOFILES=\
				mango.go\
				logger.go\
				show_errors.go\
				sessions.go\
				routing.go\

include $(GOROOT)/src/Make.pkg

format:
	${GOFMT} -w ${GOFILES}
	${GOFMT} -w mango_test.go
	${GOFMT} -w examples/*.go
