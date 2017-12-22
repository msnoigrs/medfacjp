#!/bin/bash

ROOTDIR="${PWD}"
IDIR="medfacjp-input"
ODIR="medfacjp-output"

#find "${IDIR}" -name '*.xlsx' -exec ./medfacjp data {} "${ROOTDIR}/${ODIR}" \;
find "${IDIR}" -name '*.txt' -exec ./medfacjp code {} "${ROOTDIR}/${ODIR}" \;
