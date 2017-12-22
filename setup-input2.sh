#!/bin/bash

ROOTDIR="${PWD}"

BDIR="medfacjp2017"
IDIR="medfacjp-input"
ODIR="medfacjp-output"

#~/go/bin/trdsql -id '\t' -od ' ' "select c1||'\'||c2 from /home/igarashi/go/src/github.com/msnoigrs/medfacjp/medfacjp2017/20171220_利用するファイル_0.txt order by c1"|sed -e 's/時点//g' -e 's/\\\\/\\/g' -e 's#\\#/#g' -e 's/"//g' > flist.txt

uzip() {
    cd "${ROOTDIR}/${IDIR}/${2}"
    unzip "${ROOTDIR}/${BDIR}/${1}.zip"
    cd "${ROOTDIR}"
}

cpxls() {
    bname=$(basename "${1}")
    rmext=${bname%.*}
    ftype=$(file -b "${1}" | cut -d , -f 1)
    ext=""
    if [ "${ftype}" = "Composite Document File V2 Document" ]; then
        ext=".xls"
    elif [ "${ftype}" = "Zip archive data" -o "${ftype}" = "Microsoft Excel 2007+" ]; then
        ext=".xlsx"
    else
        echo ${ftype}
    fi
    if [ -n "$ext" ]; then
        cp "${1}" "${2}/${rmext}${ext}"
    fi
}

# ファイル名に空白が入ってると面倒
while read line
do
    fd=$(dirname "${line}")
    idir="${IDIR}/${fd}"
    if [ ! -d "${idir}" ]; then
        mkdir -p "${idir}"
        echo "${idir}"
    fi
    cpxls "${BDIR}/${line}" "${idir}"
done < <(cat flist.txt)

# find "${IDIR}" -name '*.xlsx' -exec ./xlsx2xlsx.py {} \;
# find "${IDIR}" -name '*.xls' -exec ./xls2xlsx.py {} \;
# find "${IDIR}" -name '*.xls' -delete
# find "${IDIR}" -name '*.xlsx' -exec ./medfacjp xlsx2csv {} \;

# find "${IDIR}" -name '*.xlsx' -print0 | xargs -0 -P 4 ./xlsx2xlsx.py {} \;
# find "${IDIR}" -name '*.xls' -print0 | xargs -0 -P 4 ./xls2xlsx.py {} \;
# find "${IDIR}" -name '*.xls' -delete
# find "${IDIR}" -name '*.xlsx' -print0 | xargs -0 -P 4 ./medfacjp xlsx2csv {} \;
