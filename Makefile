
TARG = ng
all:
	go build -o $(TARG)

install:
	cp $(TARG) $(GOROOT)/bin

run: all
	./ng -file=project/ng.js -work=project test

