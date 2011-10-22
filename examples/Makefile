include $(GOROOT)/src/Make.inc

BINS=hello logger session routing
LIBS=cats_middleware.a silence_middleware.a

all: $(BINS) $(LIBS)

clean:
	rm -rf *.[68] $(BINS)

%: %.go
	$(GC) $*.go
	$(LD) -o $@ $*.$O

%.a: %.go
	$(GC) $*.go
	@mkdir -p _obj
	rm -f _obj/$@ $(OFILES)
	gopack grc _obj/$@ $*.$O $(OFILES)	
