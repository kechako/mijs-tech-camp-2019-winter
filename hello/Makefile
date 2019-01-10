GOCMD=go
TARGET=main.wasm
RMCMD=rm
SERVERDIR=../server
SERVERCMD=$(SERVERDIR)/server

TESTSRCS = $(wildcard *_test.go)
SRCS = $(filter-out $(TESTSRCS), $(wildcard *.go))

all: $(TARGET)

$(TARGET): $(SRCS)
	GOOS=js GOARCH=wasm $(GOCMD) build -o $(TARGET)

$(SERVERCMD):
	cd $(SERVERDIR) && make

serve: $(SERVERCMD)
	$(SERVERCMD)

clean:
	$(RMCMD) $(TARGET)

.PHONY: all serve clean

