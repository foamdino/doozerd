include ../../Make.inc

TARG=doozer/web
GOFILES=\
	main.css.go\
	main.html.go\
	stats.html.go\
	main.js.go\
	web.go\

include $(GOROOT)/src/Make.pkg

%.css.go: %.css file2gostring
	./file2gostring $(TARG) $< < $< > $@.part
	mv $@.part $@

%.html.go: %.html file2gostring
	./file2gostring $(TARG) $< < $< > $@.part
	mv $@.part $@

%.js.go: %.js file2gostring
	./file2gostring $(TARG) $< < $< > $@.part
	mv $@.part $@

CLEANFILES+=*.css.go *.html.go *.js.go
