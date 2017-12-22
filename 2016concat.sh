#!/bin/bash

ROOTDIR="${PWD}"
IDIR="20160401iryohokenlist"
ODIR="medfacjp-2016"

i=1
while [ $i -le 47 ]
do
    ii=$(printf "%02d" $i)
    mkdir -p ${ODIR}/${ii}/{ika,sika,yaku}
    cp ${IDIR}/outika/${ii}data1.txt ${ODIR}/${ii}/ika/data1.txt
    cp ${IDIR}/outsika/${ii}data1.txt ${ODIR}/${ii}/sika/data1.txt
    cp ${IDIR}/outyaku/${ii}data1.txt ${ODIR}/${ii}/yaku/data1.txt

    i=$(expr ${i} + 1)
done

ika_data1_list=$(find "${ODIR}" -name data1.txt | grep ika | grep -v sika)
sika_data1_list=$(find "${ODIR}" -name data1.txt | grep sika)
yaku_data1_list=$(find "${ODIR}" -name data1.txt | grep yaku)

cat ${ika_data1_list} > "${ODIR}"/ika_data1.txt
cat ${sika_data1_list} > "${ODIR}"/sika_data1.txt
cat ${yaku_data1_list} > "${ODIR}"/yaku_data1.txt

cat "${ODIR}"/ika_data1.txt \
    "${ODIR}"/sika_data1.txt \
    "${ODIR}"/yaku_data1.txt > "${ODIR}"/all_data1.txt
