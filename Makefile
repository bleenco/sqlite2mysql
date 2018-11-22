all:
	@make install

install:
	@chmod +x sqlite2mysql.pl
	@cp -rf sqlite2mysql.pl /usr/local/bin/sqlite2mysql

.PHONY: install
