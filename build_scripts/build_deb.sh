#!/bin/bash

if [ -z "$1" ]
  then
    echo "No argument supplied"
fi

#Set variables
SGE_INSTALLER_VERSION="$2"
#SGE_INSTALLER_VERSION=0.0.4
#PLATFORM=amd64
PLATFORM="$3"
echo "Sged Installer Version is: $SGE_INSTALLER_VERSION"
echo "Sged Installer Platform is: $PLATFORM"
export SGE_INSTALLER_VERSION
export PLATFORM

echo "Create dist/"
#rm -rf dist
mkdir -p dist

echo "Create executables"
#Build binary file
#make build
echo "Prepare directory"
# Builds directory heirarhy
DEB_BASE="sge_$SGE_INSTALLER_VERSION-1_$PLATFORM"
mkdir -p "dist/$DEB_BASE/usr/bin"
mkdir -p "dist/$DEB_BASE/usr/share/sged"
mkdir -p "dist/$DEB_BASE/DEBIAN"

echo "Copy assets"
cp build_scripts/assets/deb/control.j2 "dist/$DEB_BASE/DEBIAN/control"
cp build_scripts/assets/deb/sged.service "dist/$DEB_BASE/usr/share/sged/sged.service"
#cp build/sged "dist/$DEB_BASE/usr/bin/sged"
#cp ~/go/bin/sged "dist/$DEB_BASE/usr/bin/sged"
cp "$1" "dist/$DEB_BASE/usr/bin/sged"
/usr/bin/sha256sum -b "dist/$DEB_BASE/usr/bin/sged" > "dist/$DEB_BASE/usr/share/sged/sged.summ"

# Prepare control file
sed -i "s/{{ SGE_INSTALLER_VERSION }}/${SGE_INSTALLER_VERSION}/" "dist/$DEB_BASE/DEBIAN/control"
sed -i "s/{{ PLATFORM }}/$PLATFORM/" "dist/$DEB_BASE/DEBIAN/control"

echo "Build package"
dpkg-deb --build --root-owner-group "dist/$DEB_BASE"

rm -Rf "dist/$DEB_BASE"
