# sqlite2mysql

SQLite3 dump to MySQL dump converter.

### Exporting from SQLite3 and importing to MySQL

First, make a dump of SQLite3 database:

```sh
sqlite3 db.sqlite .dump > db.sql
```

Then, run `sqlite2mysql.pl` script to convert dump to MySQL format:

```sh
./sqlite2mysql.pl db.sql > mysql-db.sql
```

Hopefully, everything went smooth, now:

Open additional terminal instance, run below two commands and leave it open during import:

```sh
mysql -u username -p
set global net_buffer_length=1000000;
set global max_allowed_packet=1000000000;
```

and now on first terminal instance:

```sh
mysql -u username -p database_name < mysql-db.sql
```
