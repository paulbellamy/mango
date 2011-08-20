include $(GOROOT)/src/Make.inc

TARG=mango
GOFMT=gofmt

GOFILES=\
				mango.go\
				logger.go\
				show_errors.go\
				sessions.go\
				routing.go\
				static.go\
				jsonp.go\
				mime.go\
				redirect.go\

include $(GOROOT)/src/Make.pkg

format:
	${GOFMT} -w ${GOFILES}
	${GOFMT} -w *_test.go
	${GOFMT} -w examples/*.go
