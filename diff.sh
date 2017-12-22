#!/bin/bash

ROOTDIR="${PWD}"

THISYEAR="2017"
LASTYEAR=$(expr ${THISYEAR} - 1)

LDIR="medfacjp-${LASTYEAR}"
TDIR="medfacjp-output"

LALL="${LDIR}/all_data1.txt"
TALL="${TDIR}/all_data1.txt"

TRDSQL="/home/igarashi/go/bin/trdsql"
CHART="/home/igarashi/go/bin/chart"

if [ -d diff ]; then
    rm diff/*
    rmdir diff
fi
mkdir diff


i=1
while [ $i -le 47 ]
do
    ii=$(printf "%02d" $i)
    ${TRDSQL} -id '\t' -od ' ' "select c3, c4, c5 from ${LDIR}/${ii}/ika/data1.txt" > diff/"${ii}-${LASTYEAR}-data1.txt"

    ${TRDSQL} -id '\t' -od ' ' "select c3, c4, c5 from ${TDIR}/${ii}/ika/data1.txt" > diff/"${ii}-${THISYEAR}-data1.txt"
    diff -u diff/"${ii}-${LASTYEAR}-data1.txt" diff/"${ii}-${THISYEAR}-data1.txt" > diff/"${ii}-${LASTYEAR}to${THISYEAR}-data1.txt"

    i=$(expr ${i} + 1)
done
