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

#${TRDSQL} -id '\t' -oat "select c1 as no, c2 as date, c3 as kubun, c4 as code, c5 as name, c6 as postalcode, c7 as address, c8 as tel, c9 as fax from ${LALL} limit 1"

if [ -d check ]; then
    rm check/*
    rmdir check
fi
mkdir check

######### all
find "${TDIR}" -name data1.txt -print > check/data1list.txt
find "${TDIR}" -name data2.txt -print > check/data2list.txt
data1ika_num=$(grep ika check/data1list.txt| grep -v sika | wc -l)
data1sika_num=$(grep sika check/data1list.txt | wc -l)
data1yaku_num=$(grep yaku check/data1list.txt | wc -l)
data2ika_num=$(grep ika check/data2list.txt| grep -v sika | wc -l)
data2sika_num=$(grep sika check/data2list.txt | wc -l)
data2yaku_num=$(grep yaku check/data2list.txt | wc -l)

if [ "${data1ika_num}" -eq "47" -a "${data1sika_num}" -eq "47" -a "${data1yaku_num}" -eq "47" ]; then
    echo "PASS: number of data1.txt is valid."
else
    i=1
    while [ $i -le 47 ]
    do
        ii=$(printf "%02d" $i)
        ikadata1="${TDIR}/${ii}/ika/data1.txt"
        sikadata1="${TDIR}/${ii}/sika/data1.txt"
        yakudata1="${TDIR}/${ii}/sika/data1.txt"
        if [ ! -f "${ikadata1}" ]; then
            echo "NG: ${ikadata1} not found"
        fi
        if [ ! -f "${sikadata1}" ]; then
            echo "NG: ${sikadata1} not found"
        fi
        if [ ! -f "${yakudata1}" ]; then
            echo "NG: ${yakudata1} not found"
        fi
        i=$(expr ${i} + 1)
    done
fi
if [ "${data2ika_num}" -eq "47" -a "${data2sika_num}" -eq "47" -a "${data2yaku_num}" -eq "47" ]; then
    echo "PASS: number of data2.txt is valid."
else
    i=1
    while [ $i -le 47 ]
    do
        ii=$(printf "%02d" $i)
        ikadata2="${TDIR}/${ii}/ika/data2.txt"
        sikadata2="${TDIR}/${ii}/sika/data2.txt"
        yakudata2="${TDIR}/${ii}/sika/data2.txt"
        if [ ! -f "${ikadata2}" ]; then
            echo "NG: ${ikadata2} not found"
        fi
        if [ ! -f "${sikadata2}" ]; then
            echo "NG: ${sikadata2} not found"
        fi
        if [ ! -f "${yakudata2}" ]; then
            echo "NG: ${yakudata2} not found"
        fi
        i=$(expr ${i} + 1)
    done
fi

l_all_count=$(${TRDSQL} -id '\t' -od ' ' "select count(c4) from ${LALL}")
t_all_count=$(${TRDSQL} -id '\t' -od ' ' "select count(c4) from ${TALL}")

l_all_count_d=$(${TRDSQL} -id '\t' -od ' ' "select count(distinct c4) from ${LALL}")
t_all_count_d=$(${TRDSQL} -id '\t' -od ' ' "select count(distinct c4) from ${TALL}")

echo "${LASTYEAR}: ${l_all_count} ${l_all_count_d} records"
echo "${THISYEAR}: ${t_all_count} ${t_all_count_d} records"

${TRDSQL} -id '\t' -od ' ' "select c4, count from (select c4, count(c4) as count from ${LALL} group by c4) where count > 1" > check/double-last.txt
${TRDSQL} -id '\t' -od ' ' "select c4, count from (select c4, count(c4) as count from ${TALL} group by c4) where count > 1" > check/double-this.txt

# echo "year${LASTYEAR}:${l_all_count} ${l_all_count}" > check/allrecords.txt
# echo "year${THISYEAR}:${t_all_count} ${t_all_count}" >> check/allrecords.txt

# ${TRDSQL} -od '\t' -id ' ' "select * from check/allrecords.txt" | ${CHART} bar

########## kubun

${TRDSQL} -id '\t' -od ' ' "select '${LASTYEAR}', c3, count(c4), count(distinct c4) from ${LALL} group by c3" > check/kubunrecords.txt
${TRDSQL} -id '\t' -od ' ' "select '${THISYEAR}', c3, count(c4), count(distinct c4) from ${TALL} group by c3" >> check/kubunrecords.txt

cat check/kubunrecords.txt

${TRDSQL} -id '\t' "select count(*) from ${LALL} group by c3" > check/kubun-last.txt
${TRDSQL} -id '\t' "select count(*) from ${TALL} group by c3" > check/kubun-this.txt

sed -i -r -e ':loop;N;$!b loop;s/\n/ /g' -e 's/ +/\t/g' check/kubun-last.txt
sed -i -r -e ':loop;N;$!b loop;s/\n/ /g' -e 's/ +/\t/g' check/kubun-this.txt

printf "year${LASTYEAR}\t" > check/kubunb.txt
cat check/kubun-last.txt >> check/kubunb.txt
printf "year${THISYEAR}\t" >> check/kubunb.txt
cat check/kubun-this.txt >> check/kubunb.txt

${TRDSQL} -od '\t' -id '\t' "select * from check/kubunb.txt" | ${CHART} bar

# sed -i -e 's/医科/ika/g' \
#     -e 's/歯科/sika/g' \
#     -e 's/薬局/yaku/g' \
#     check/kubunrecords.txt

# grep ika check/kubunrecords.txt | grep -v sika > check/ika.txt
# grep sika check/kubunrecords.txt > check/sika.txt
# grep yaku check/kubunrecords.txt > check/yaku.txt

# ${TRDSQL} -od '\t' -id ' ' "select c3 from ${ROOTDIR}/check/ika.txt" > check/ikayoko.txt
# ${TRDSQL} -od '\t' -id ' ' "select c3 from ${ROOTDIR}/check/sika.txt" > check/sikayoko.txt
# ${TRDSQL} -od '\t' -id ' ' "select c3 from ${ROOTDIR}/check/yaku.txt" > check/yakuyoko.txt

# sed -i -r -e ':loop;N;$!b loop;s/\n/ /g' -e 's/ +/ /g' check/ikayoko.txt
# sed -i -r -e ':loop;N;$!b loop;s/\n/ /g' -e 's/ +/ /g' check/sikayoko.txt
# sed -i -r -e ':loop;N;$!b loop;s/\n/ /g' -e 's/ +/ /g' check/yakuyoko.txt

# echo -n "ika " > check/kubun.txt
# cat check/ikayoko.txt >> check/kubun.txt
# echo -n "sika " >> check/kubun.txt
# cat check/sikayoko.txt >> check/kubun.txt
# echo -n "yaku " >> check/kubun.txt
# cat check/yakuyoko.txt >> check/kubun.txt

# # ${TRDSQL} -od '\t' -id ' ' "select * from check/kubun.txt" | ${CHART} line --date-format 2006

######## pref

${TRDSQL} -id '\t' -od '\t' "select 'pref' || pref, count(c4) from (select substr(c4, 1, 2) as pref, c4 from ${LALL} where c3 = '医科') group by pref order by pref" > check/pref-ika-last.txt
${TRDSQL} -id '\t' -od '\t' "select count(c4) from (select substr(c4, 1, 2) as pref, c4 from ${TALL} where c3 = '医科') group by pref order by pref" > check/pref-ika-this.txt

paste check/pref-ika-last.txt check/pref-ika-this.txt > check/pref-ika.txt

${CHART} bar --title "ika" < check/pref-ika.txt


${TRDSQL} -id '\t' -od '\t' "select 'pref' || pref, count(c4) from (select substr(c4, 1, 2) as pref, c4 from ${LALL} where c3 = '歯科') group by pref order by pref" > check/pref-sika-last.txt
${TRDSQL} -id '\t' -od '\t' "select count(c4) from (select substr(c4, 1, 2) as pref, c4 from ${TALL} where c3 = '歯科') group by pref order by pref" > check/pref-sika-this.txt

paste check/pref-sika-last.txt check/pref-sika-this.txt > check/pref-sika.txt

${CHART} bar --title "sika" < check/pref-sika.txt


${TRDSQL} -id '\t' -od '\t' "select 'pref' || pref, count(c4) from (select substr(c4, 1, 2) as pref, c4 from ${LALL} where c3 = '薬局') group by pref order by pref" > check/pref-yaku-last.txt
${TRDSQL} -id '\t' -od '\t' "select count(c4) from (select substr(c4, 1, 2) as pref, c4 from ${TALL} where c3 = '薬局') group by pref order by pref" > check/pref-yaku-this.txt

paste check/pref-yaku-last.txt check/pref-yaku-this.txt > check/pref-yaku.txt

${CHART} bar --title "yaku" < check/pref-yaku.txt

${TRDSQL} -id '\t' -od '\t' "select c1, c3 - c2 from ${ROOTDIR}/check/pref-ika.txt" > check/pref-diff-ika.txt
${TRDSQL} -id '\t' -od '\t' "select c1, c3 - c2 from ${ROOTDIR}/check/pref-sika.txt" > check/pref-diff-sika.txt
${TRDSQL} -id '\t' -od '\t' "select c1, c3 - c2 from ${ROOTDIR}/check/pref-yaku.txt" > check/pref-diff-yaku.txt

paste check/pref-diff-ika.txt check/pref-diff-sika.txt check/pref-diff-yaku.txt > check/pref-diff.txt

${CHART} bar --title "diff" < check/pref-diff.txt

i=1
while [ $i -le 47 ]
do
    ii=$(printf "%02d" $i)
    tail -n 1 ${LDIR}/${ii}/ika/data1.txt >> check/tail-last-ika.txt
    tail -n 1 ${TDIR}/${ii}/ika/data1.txt >> check/tail-this-ika.txt
    tail -n 1 ${LDIR}/${ii}/sika/data1.txt >> check/tail-last-sika.txt
    tail -n 1 ${TDIR}/${ii}/sika/data1.txt >> check/tail-this-sika.txt
    tail -n 1 ${LDIR}/${ii}/yaku/data1.txt >> check/tail-last-yaku.txt
    tail -n 1 ${TDIR}/${ii}/yaku/data1.txt >> check/tail-this-yaku.txt
    i=$(expr ${i} + 1)
done
${TRDSQL} -id '\t' -od '\t' "select c4, c5, c6, c7 from ${ROOTDIR}/check/tail-last-ika.txt" > check/tail-last-ika-m.txt
${TRDSQL} -id '\t' -od '\t' "select c4, c5, c6, c7 from ${ROOTDIR}/check/tail-this-ika.txt" > check/tail-this-ika-m.txt
${TRDSQL} -id '\t' -od '\t' "select c4, c5, c6, c7 from ${ROOTDIR}/check/tail-last-sika.txt" > check/tail-last-sika-m.txt
${TRDSQL} -id '\t' -od '\t' "select c4, c5, c6, c7 from ${ROOTDIR}/check/tail-this-sika.txt" > check/tail-this-sika-m.txt
${TRDSQL} -id '\t' -od '\t' "select c4, c5, c6, c7 from ${ROOTDIR}/check/tail-last-yaku.txt" > check/tail-last-yaku-m.txt
${TRDSQL} -id '\t' -od '\t' "select c4, c5, c6, c7 from ${ROOTDIR}/check/tail-this-yaku.txt" > check/tail-this-yaku-m.txt

echo "last record ika diff"
diff -u check/tail-last-ika-m.txt check/tail-this-ika-m.txt
echo "last record sika diff"
diff -u check/tail-last-sika-m.txt check/tail-this-sika-m.txt
echo "last record yaku diff"
diff -u check/tail-last-yaku-m.txt check/tail-this-yaku-m.txt
