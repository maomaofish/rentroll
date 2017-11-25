#!/bin/bash
MAJVER=1
MINVER=2
if [ "${BUILD_NUMBER}x" = "x" ]; then
	BLDNOFILE="./buildno"
	BUILDNO=$(cat ${BLDNOFILE})
	BUILDNO=$((BUILDNO + 1))
else
	BUILDNO=${BUILD_NUMBER}
fi
FNAME="ver.go"
VER=$(printf "%d.%d.%06d" ${MAJVER} ${MINVER} ${BUILDNO})
BLD="${HOSTNAME}"
DAT=$(date)
PKG="package main"

if [ "x${1}" != "x" ]; then 
    PKG="package ${1}"
fi

cat >${FNAME} <<ZZ1EOF
${PKG}
// THIS FILE IS AUTOGENERATED
// DO NOT EDIT

// GetVersionNo returns the current version number
func GetVersionNo() string { return "${VER}" }
// GetBuildMachine returns the machine on which this executable was built
func GetBuildMachine() string { return "${BLD}" }
// GetBuildTime returns the time this executable was built
func GetBuildTime() string { return "${DAT}" }
ZZ1EOF

echo "${BUILDNO}" >${BLDNOFILE}
