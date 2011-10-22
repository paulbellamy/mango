include $(GOROOT)/src/Make.inc

BINS=hello logger session routing cats
LIBS=cats_middleware.a silence_middleware.a

all: $(LIBS) $(BINS)

clean:
	rm -rf *.[68] $(BINS)

%: %.go
	$(GC) -I _obj $*.go
	$(LD) -L _obj -o $@ $*.$O

%.a: %.go
	$(GC) $*.go
	@mkdir -p _obj
	rm -f _obj/$@ $(OFILES)
	gopack grc _obj/$@ $*.$O $(OFILES)	
