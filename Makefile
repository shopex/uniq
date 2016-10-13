OBJ=uniq

all: ${OBJ}

${OBJ}: *.go
	go build .

start: ${OBJ}
	./${OBJ} -datadir data -days 3 -port 6532
