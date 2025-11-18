#!/bin/sh

check_code() {
        EXCODE=$?
        if [ "$EXCODE" != "0" ]; then
                echo "build fail."
                exit $EXCODE
        fi
}

out="runtime"
echo "build file to ./$out"

mkdir -p "$out/conf"

GO111MODULE=on go build -o ./$out/cronnode ./cmd/cronnode
check_code
GO111MODULE=on go build -o ./$out/cronweb ./cmd/cronweb
check_code

sources=`find ./conf/files -name "*.json.sample"`
check_code
for source in $sources;do
        yes | echo $source|sed "s/.*\/\(.*\.json\).*/cp -f & .\/$out\/conf\/\1/"|bash
        check_code
done

echo "build success."
