COMPILER=go build
OPTIONS=

BIN=distru
BINDIR=/usr/local/bin
CONF=/etc/distru.conf
WEBDIR=/etc/distru

SRCWEBDIR=ui

GENCONF=--genconf
SETWEBDIR=--webdir $(WEBDIR)

.DEFAULT_GOAL := distru

distru:
	$(COMPILER) $(OPTIONS)

.PHONY : clean
clean  :
	rm $(BIN)

.PHONY  : install
install : distru
	./$(BIN) $(GENCONF) $(SETWEBDIR) $(CONF)
	mv $(BIN) $(BINDIR)/
	mkdir $(WEBDIR)
	cp $(SRCWEBDIR)/* $(WEBDIR)/

.PHONY    : uninstall
uninstall :
	rm $(BINDIR)/$(BIN)
	rm $(CONF)
	rm -r $(WEBDIR)