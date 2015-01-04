#!/bin/bash
# Counts non-generated lines by type.

git rm -rf Godeps >/dev/null 2>&1
./format.sh
SQL=$(find . -name "*.sql" | xargs cat | wc -l;)
SH=$(find . -name "*.sh" | xargs cat | wc -l;)
GO=$(find . -name "*.go" | xargs cat | wc -l;)
echo "Go: $GO"
echo "Shell: $SH"
echo "SQL: $SQL"
TOTAL=`expr $GO + $SH + $SQL`
echo "Total: $TOTAL"
godep save
git add .
