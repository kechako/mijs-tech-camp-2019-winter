GOCMD=go
TARGET=server
RMCMD=rm

TESTSRCS = $(wildcard *_test.go)
SRCS = $(filter-out $(TESTSRCS), $(wildcard *.go))

all: $(TARGET)

$(TARGET): $(SRCS)
	$(GOCMD) build -o $(TARGET)

clean:
	$(RMCMD) $(TARGET)

.PHONY: all clean

