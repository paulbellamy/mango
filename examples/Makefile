include $(GOROOT)/src/Make.inc

ALL=hello logger session cats_middleware silence_middleware

all: $(ALL)

clean:
	rm -rf *.[68] $(ALL)

%: %.go
	$(GC) $*.go
	$(LD) -o $@ $*.$O
