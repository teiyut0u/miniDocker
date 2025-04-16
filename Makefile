PROJECT=minidocker
GOC=go

BIN_DIR=$(abspath bin)

CLI_DIR=$(abspath cli)
CLI_SRC=$(shell find $(CLI_DIR) -name '*.go')
CLI_TARGET=$(BIN_DIR)/$(PROJECT)-cli

RUNTIME_DIR=$(abspath runtime)
RUNTIME_SRC=$(shell find $(RUNTIME_DIR) -name '*.go')
RUNTIME_TARGET=$(BIN_DIR)/$(PROJECT)-runtime

all: $(CLI_TARGET) $(RUNTIME_TARGET)

$(CLI_TARGET):$(CLI_SRC)|dir
	go build -o $@ $(CLI_DIR)

$(RUNTIME_TARGET):$(RUNTIME_SRC)|dir
	go build -o $@ $(RUNTIME_DIR)

dir:
	@mkdir -pv bin

clean:
	@rm -rf bin

.PHONY:
	clean dir
