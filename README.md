Nutid prints out the current time in Local, UTC and Unix time format.

You can also pass a text file to nutid over stdin, and it will replace all timestamps it finds to local time (assuming the file contains UTC timestamps) e.g.

    cat application.log | nutid

You can convert a unix timestamp

    nutid -unix 1528088376

You can add a time offset

    nutid -add 24h

or substract time using a negative add

    nutid -add -24h

It colors the output by default, to disable coloring add no-colors

    nutid -no-colors