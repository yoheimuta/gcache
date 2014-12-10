#!bin/sh -xe
# see #1

usage_exit() {
    echo "Usage: $0 [-v]" 1>&2
    exit 1
}

while getopts v OPT
do
    case $OPT in
        v)  FLAG_V=1
            ;;
        \?) usage_exit
            ;;
    esac
done

CWD=`dirname "${0}"`
expr "${0}" : "/.*" > /dev/null || CWD=`(cd "${CWD}" && pwd)`
cd $CWD
cd ..

OUT="${CWD}/gotest.out"
rm $OUT

find * -name '*_test.go' | while read file; do
n=$(dirname -- "$file")
echo "$n"
done | sort -u | while read d; do
    if [ "$FLAG_V" ]; then
        go test -v "$d" | tee -a "${OUT}"
    else
        go test "$d" | tee -a "${OUT}"
    fi
done

if ! grep -e "exit status" -e "FAIL" "$OUT" >/dev/null 2>&1; then
    rm $OUT
    echo "Succeeded gotest.sh"
    exit 0
else
    rm $OUT
    echo "Failed gotest.sh"
    exit 1
fi
