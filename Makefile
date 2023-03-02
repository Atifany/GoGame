#name
NAME := GoGame
BUILD_DIR := builds/linux
BIN := $(NAME:%=$(BUILD_DIR)/%)

#sources
_SRC :=	button.go	\
		cell.go		\
		draw.go		\
		main.go		\
		parser.go

SRC_DIR := sources
SRC := $(_SRC:%=$(SRC_DIR)/%)

start: compile

compile: $(SRC)
	go build -o $(BIN) $(SRC)