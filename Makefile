NODES=./nodes/main.wasm
SQUARES=./marching_squares/main.wasm

.PHONY: all
all: nodes squares

nodes: ${NODES}

${NODES}:
	@GOOS=js GOARCH=wasm go build -o ${NODES} ./nodes

squares: ${SQUARES}

${SQUARES}:
	@GOOS=js GOARCH=wasm go build -o ${SQUARES} ./marching_squares

clear:
	@rm ${NODES}
	@rm ${SQUARES}
