# Use the default go compiler
GO_BUILDFLAGS=-compiler gc
# Or uncomment the line below to use the gccgo compiler, which may
# or may not be faster than gc and which may or may not compile...
# GO_BUILDFLAGS=-compiler gccgo -gccgoflags '-static-libgcc -O4 -Ofast -march=native'

CGO_CFLAGS_ALLOW='(-fno-schedule-insns|-malign-double|-ffast-math)'

BUILD_TARGETS = all install 6g gccgo test 6gtest gccgotest bench 6gtest gccgotest clean stubs

SRCFILES := cgoflags.go cl.go cl_test.go context.go device.go event.go kernel.go memory.go platform.go program.go queue.go vkfft.go
OSFLAG := 

ifeq ($(OS), Windows_NT)
	OSFLAG +=  -tag windows
	SRCFILES += " context_win.go"
endif

.PHONY: $(BUILD_TARGETS)


all: install


install: stubs
	go install -v $(GO_BUILDFLAGS)


6g:
	go install -v
	go vet $(SRCFILES)
	gofmt -w $(SRCFILES)


GCCGO=gccgo -gccgoflags '-static-libgcc -O3'


gccgo:
	go build -v -compiler $(GCCGO)


test: 6gtest gccgotest


6gtest:
	go test


gccgotest:
	go test -compiler $(GCCGO)


bench: 6gbench gccgobench


6gbench:
	go test -bench=.


gccgobench:
	go test -bench=. -compiler $(GCCGO)


stubs:
	$(MAKE) -C ./stubs all


clean:
	make -C ./stubs clean
	go clean

