GOCMD=go
TARGET=hello.wasm
RMCMD=rm
SERVERDIR=../server
SERVERCMD=$(SERVERDIR)/server

SRCS = hello.go

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

