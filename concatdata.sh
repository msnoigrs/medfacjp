#!/bin/bash

ROOTDIR="${PWD}"
IDIR="medfacjp-input"
ODIR="medfacjp-output"

ika_data1_list=$(find "${ODIR}" -name data1.txt | grep ika | grep -v sika)
sika_data1_list=$(find "${ODIR}" -name data1.txt | grep sika)
yaku_data1_list=$(find "${ODIR}" -name data1.txt | grep yaku)

ika_data2_list=$(find "${ODIR}" -name data2.txt | grep ika | grep -v sika)
sika_data2_list=$(find "${ODIR}" -name data2.txt | grep sika)
yaku_data2_list=$(find "${ODIR}" -name data2.txt | grep yaku)

# ika_data3_list=$(find "${ODIR}" -name data3.txt | grep ika | grep -v sika)
# sika_data3_list=$(find "${ODIR}" -name data3.txt | grep sika)
# yaku_data3_list=$(find "${ODIR}" -name data3.txt | grep yaku)

cat ${ika_data1_list} > "${ODIR}"/ika_data1.txt
cat ${sika_data1_list} > "${ODIR}"/sika_data1.txt
cat ${yaku_data1_list} > "${ODIR}"/yaku_data1.txt

cat ${ika_data2_list} > "${ODIR}"/ika_data2.txt
cat ${sika_data2_list} > "${ODIR}"/sika_data2.txt
cat ${yaku_data2_list} > "${ODIR}"/yaku_data2.txt

# cat ${ika_data3_list} > "${ODIR}"/ika_data3.txt
# cat ${sika_data3_list} > "${ODIR}"/sika_data3.txt
# cat ${yaku_data3_list} > "${ODIR}"/yaku_data3.txt

cat "${ODIR}"/ika_data1.txt \
    "${ODIR}"/sika_data1.txt \
    "${ODIR}"/yaku_data1.txt > "${ODIR}"/all_data1.txt
cat "${ODIR}"/ika_data2.txt \
    "${ODIR}"/sika_data2.txt \
    "${ODIR}"/yaku_data2.txt > "${ODIR}"/all_data2.txt
# cat "${ODIR}"/ika_data3.txt \
#     "${ODIR}"/sika_data3.txt \
#     "${ODIR}"/yaku_data3.txt > "${ODIR}"/all_data3.txt
