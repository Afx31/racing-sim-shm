BINARY_NAME_BRIDGE=./udp-bridge/bin/udp-bridge
BINARY_DIR_BRIDGE=./udp-bridge/
BINARY_NAME_CAN=./udp-to-can/bin/udp-to-can
BINARY_DIR_CAN=./udp-to-can/

compile-bridge:
	go build -o ${BINARY_NAME_BRIDGE} ${BINARY_DIR_BRIDGE}

compile-can:
	go build -o ${BINARY_NAME_CAN} ${BINARY_DIR_CAN}

rebuild-bridge: compile-bridge
	./${BINARY_NAME_BRIDGE}

rebuild-can: compile-can
	./${BINARY_NAME_CAN}

run-bridge:
	./${BINARY_NAME_BRIDGE}

run-can:
	./${BINARY_NAME_CAN}
