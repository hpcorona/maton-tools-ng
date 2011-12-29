include $(GOROOT)/src/Make.inc

TARG = ng
GOFILES = \
				ng.go \
				func.go


include $(GOROOT)/src/Make.cmd

run: all
	./ng -file=project/ng.js -work=project

