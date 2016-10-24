OBJ=uniq

all: ${OBJ}

${OBJ}: *.go
	go build .

start: ${OBJ}
	./${OBJ} -datadir data -days 30 -port 6532
