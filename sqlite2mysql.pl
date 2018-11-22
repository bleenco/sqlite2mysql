#!/usr/bin/env perl

use strict;
use warnings;

my $file = "abstruse.sqlite-dump";
open my $info, $file or die "Could not open $file: $!";

# print "SET @@global.sql_mode= 'NO_BACKSLASH_ESCAPES,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION';\nSET FOREIGN_KEY_CHECKS=0;\n"

while (my $line = <$info>) {
    if (($line !~  /BEGIN TRANSACTION/) && ($line !~ /COMMIT/) && ($line !~ /sqlite_sequence/) && ($line !~ /CREATE UNIQUE INDEX/) && ($line !~ /PRAGMA foreign_keys=OFF/)){

        $line =~ s/\"/\`/g;
        $line =~ s/CREATE TABLE IF NOT EXISTS \`(.*?)\` (.*)\;/DROP TABLE IF EXISTS $1;\nCREATE TABLE IF NOT EXISTS $1 $2;\n/;
        $line =~ s/(CREATE TABLE.*)(primary key) (autoincrement)(.*)()\);/$1AUTO_INCREMENT$4, PRIMARY KEY(id))$5;/;

        while ($line =~ /(.*)([0-9]{13})(.*)/g) {
            my $n = substr($2, 0, 10);
            $line = "$1$n$3";
        }

        # $line =~ s/\b([0-9]{13})\b/FROM_UNIXTIME($1)/g;
        $line =~ s/\b([0-9]{10})\b/FROM_UNIXTIME($1)/g;

        # if ($line =~ /CREATE TABLE IF NOT EXISTS \`(.*?)\` (.*)\;/) {
        #         $line = "DROP TABLE IF EXISTS $1;\nCREATE TABLE IF NOT EXISTS $1 $2;\n";
        # }
        if ($line =~ /INSERT INTO \"(\w*)\"(.*)/){
                $line = "INSERT INTO $1$2\n";
                $line =~ s/\"/\\\"/g;
                $line =~ s/\"/\'/g;
        }
        # if ($line =~ /(CREATE TABLE.*)(primary key) (autoincrement)(.*)()\);/) {
        #         $line = "$1AUTO_INCREMENT$4, PRIMARY KEY(id))$5;\n"
        # }
        $line =~ s/\'\\n\'/\'\'/g;
        # }

        $line =~ s/\'t\'/1/g;
        $line =~ s/\'f\'/0/g;
        $line =~ s/AUTOINCREMENT/AUTO_INCREMENT/g;
        $line =~ s/(.*)varchar([^\(0-9\)].*)/$1varchar\(255\)$2/g;
        print $line;
        print "\n";
    }
}

close $info;
