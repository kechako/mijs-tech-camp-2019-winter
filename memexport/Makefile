WAT2WASMCMD=wat2wasm
TARGET=memexport.wasm
RMCMD=rm
SERVERDIR=../server
SERVERCMD=$(SERVERDIR)/server

SRCS = memexport.wat

all: $(TARGET)

$(TARGET): $(SRCS)
	$(WAT2WASMCMD) -o $(TARGET) $(SRCS)

$(SERVERCMD):
	cd $(SERVERDIR) && make

serve: $(SERVERCMD)
	$(SERVERCMD)

clean:
	$(RMCMD) $(TARGET)

.PHONY: all serve clean

