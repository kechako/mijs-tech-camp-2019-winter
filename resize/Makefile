GOCMD=go
TARGET=resize.wasm
RMCMD=rm
SERVERDIR=../server
SERVERCMD=$(SERVERDIR)/server

SRCS = resize.go

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

