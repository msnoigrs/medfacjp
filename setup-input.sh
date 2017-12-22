#!/bin/bash

ROOTDIR="${PWD}"

BDIR="medfacjp2017"
IDIR="medfacjp-input"
ODIR="medfacjp-output"

mkdir -p medfacjp-input/{20170403/05_kinki,20170424/04_tokai,20170502/{01_hokkaido,07_shikoku,08_kyushu},20170523/{02_tohoku,03_kanto,06_chugoku,20170524_ngfile_manual_retry}}
#mkdir -p medfacjp-input/{20170424/04_tokai,20170502/{01_hokkaido,05_kinki,07_shikoku,08_kyushu},20170523/{02_tohoku,03_kanto,06_chugoku,20170524_ngfile_manual_retry}}

uzip() {
    cd "${ROOTDIR}/${IDIR}/${2}"
    unzip "${ROOTDIR}/${BDIR}/${1}.zip"
    cd "${ROOTDIR}"
}

cpxls() {
    of=$(file -b "${ROOTDIR}/${BDIR}/${1}" | cut -d , -f 1)
    if [ "${of}" = "Zip archive data" -o "${of}" = "Microsoft Excel 2007+" ]; then
        cp "${ROOTDIR}/${BDIR}/${1}" "${ROOTDIR}/${IDIR}/${1}x"
    else
        cp "${ROOTDIR}/${BDIR}/${1}" "${ROOTDIR}/${IDIR}/${2}"
    fi
}


D="20170403"
P="${D}/05_kinki"
cpxls "${P}/3fukui-sisetukijun-ika.xls" "${P}"
cpxls "${P}/3fukui-sisetukijun-sika.xls" "${P}"
cpxls "${P}/3fukui-sisetukijun-yakkyoku.xls" "${P}"
cpxls "${P}/3hyougo-sisetukijun-ika.xls" "${P}"
cpxls "${P}/3hyougo-sisetukijun-sika.xls" "${P}"
cpxls "${P}/3hyougo-sisetukijun-yakkyoku.xls" "${P}"
cpxls "${P}/3kyoto-sisetukijun-ika.xls" "${P}"
cpxls "${P}/3kyoto-sisetukijun-sika.xls" "${P}"
cpxls "${P}/3kyoto-sisetukijun-yakkyoku.xls" "${P}"
cpxls "${P}/3nara-sisetukijun-ika.xls" "${P}"
cpxls "${P}/3nara-sisetukijun-sika.xls" "${P}"
cpxls "${P}/3nara-sisetukijun-yakkyoku.xls" "${P}"
cpxls "${P}/3oosaka-sisetukijun-ika.xls" "${P}"
cpxls "${P}/3oosaka-sisetukijun-sika.xls" "${P}"
cpxls "${P}/3oosaka-sisetukijun-yakkyoku.xls" "${P}"
cpxls "${P}/3siga-sisetukijun-ika.xls" "${P}"
cpxls "${P}/3siga-sisetukijun-sika.xls" "${P}"
cpxls "${P}/3siga-sisetukijun-yakkyoku.xls" "${P}"
cpxls "${P}/3wakayama-sisetukijun-ika.xls" "${P}"
cpxls "${P}/3wakayama-sisetukijun-sika.xls" "${P}"
cpxls "${P}/3wakayama-sisetukijun-yakkyoku.xls" "${P}"

D="20170424"
P="${D}/04_tokai"
uzip "${P}/todokede_ika201704" "${P}"
uzip "${P}/todokede_yaku201704" "${P}"
uzip "${P}/todokede_shika201704" "${P}"

D="20170502"
P="${D}/01_hokkaido"
cpxls "${P}/hokkaido-sisetukijyun-ika.xls" "${P}"
cpxls "${P}/hokkaido-sisetukijyun-shika.xls" "${P}"
cpxls "${P}/hokkaido-sisetukijyun-yakkyoku.xls" "${P}"

# P="${D}/05_kinki"
# cpxls "${P}/4fukui-sisetukijun-ika.xls" "${P}"
# cpxls "${P}/4fukui-sisetukijun-sika.xls" "${P}"
# cpxls "${P}/4fukui-sisetukijun-yakkyoku.xls" "${P}"
# cpxls "${P}/4hyougo-sisetukijun-ika.xls" "${P}"
# cpxls "${P}/4hyougo-sisetukijun-sika.xls" "${P}"
# cpxls "${P}/4hyougo-sisetukijun-yakkyoku.xls" "${P}"
# cpxls "${P}/4kyoto-sisetukijun-ika.xls" "${P}"
# cpxls "${P}/4kyoto-sisetukijun-sika.xls" "${P}"
# cpxls "${P}/4kyoto-sisetukijun-yakkyoku.xls" "${P}"
# cpxls "${P}/4nara-sisetukijun-ika.xls" "${P}"
# cpxls "${P}/4nara-sisetukijun-sika.xls" "${P}"
# cpxls "${P}/4nara-sisetukijun-yakkyoku.xls" "${P}"
# cpxls "${P}/4oosaka-sisetukijun-ika.xls" "${P}"
# cpxls "${P}/4oosaka-sisetukijun-sika.xls" "${P}"
# cpxls "${P}/4oosaka-sisetukijun-yakkyoku.xls" "${P}"
# cpxls "${P}/4siga-sisetukijun-ika.xls" "${P}"
# cpxls "${P}/4siga-sisetukijun-sika.xls" "${P}"
# cpxls "${P}/4siga-sisetukijun-yakkyoku.xls" "${P}"
# cpxls "${P}/4wakayama-sisetukijun-ika.xls" "${P}"
# cpxls "${P}/4wakayama-sisetukijun-sika.xls" "${P}"
# cpxls "${P}/4wakayama-sisetukijun-yakkyoku.xls" "${P}"

P="${D}/07_shikoku"
cpxls "${P}/02_12.xls" "${P}"
cpxls "${P}/02_01.xls" "${P}"
cpxls "${P}/02_02.xls" "${P}"
cpxls "${P}/02_03.xls" "${P}"
cpxls "${P}/02_04.xls" "${P}"
cpxls "${P}/02_05.xls" "${P}"
cpxls "${P}/02_06.xls" "${P}"
cpxls "${P}/02_07.xls" "${P}"
cpxls "${P}/02_08.xls" "${P}"
cpxls "${P}/02_09.xls" "${P}"
cpxls "${P}/02_10.xls" "${P}"
cpxls "${P}/02_11.xls" "${P}"

P="${D}/08_kyushu"
uzip "${P}/shisetsu_fukuoka_02" "${P}"
uzip "${P}/shisetsu_kagoshima_02" "${P}"
uzip "${P}/shisetsu_kumamoto_02" "${P}"
uzip "${P}/shisetsu_miyazaki_02" "${P}"
uzip "${P}/shisetsu_nagasaki_02" "${P}"
uzip "${P}/shisetsu_okinawa_02" "${P}"
uzip "${P}/shisetsu_ooita_02" "${P}"
uzip "${P}/shisetsu_saga_02" "${P}"

D="20170523"
P="${D}/02_tohoku"
uzip "${P}/todokede_zentai_ika_h2904" "${P}"
uzip "${P}/todokede_zentai_shika_h2904" "${P}"
uzip "${P}/todokede_zentai_yakkyoku_h2904" "${P}"

P="${D}/03_kanto"
uzip "${P}/shisetsu_ika1_h2905" "${P}"
uzip "${P}/shisetsu_ika2_h2905" "${P}"
uzip "${P}/shisetsu_yakkyoku_h2905" "${P}"

P="${D}/06_chugoku"
cpxls "${P}/31de2904.xls" "${P}"
cpxls "${P}/32de2904.xls" "${P}"
cpxls "${P}/33de2904.xls" "${P}"
cpxls "${P}/34de2904.xls" "${P}"
cpxls "${P}/35de2904-1.xls" "${P}"

P="${D}/20170524_ngfile_manual_retry"
uzip "${P}/shisetsu_shika_h2905" "${P}"

find "${IDIR}" -name '*.xlsx' -exec ./xlsx2xlsx.py {} \;
find "${IDIR}" -name '*.xls' -exec ./xls2xlsx.py {} \;
find "${IDIR}" -name '*.xls' -delete
#find "${IDIR}" -name '*.xlsx' -exec ./medfacjp data {} "${ROOTDIR}/${ODIR}" \;
