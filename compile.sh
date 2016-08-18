#!/bin/bash -e

# create build dir
[ -d "build" ] || mkdir "build"

echo "[compile] linux arm";
env GOOS=linux GOARCH=arm go build -o build/pbflint.arm6.bin;
chmod +x build/pbflint.arm6.bin;

echo "[compile] linux x64";
env GOOS=linux GOARCH=amd64 go build -o build/pbflint.linux.bin;
chmod +x build/pbflint.linux.bin;

echo "[compile] darwin x64";
env GOOS=darwin GOARCH=amd64 go build -o build/pbflint.osx.bin;
chmod +x build/pbflint.osx.bin;

echo "[compile] windows x64";
env GOOS=windows GOARCH=386 go build -o build/pbflint.exe;

# ensure the files were compiled to the correct architecture
declare -A matrix
matrix["build/pbflint.osx.bin"]="Mach-O 64-bit x86_64 executable"
matrix["build/pbflint.arm6.bin"]="ELF 32-bit LSB executable, ARM, EABI5 version 1 (SYSV), statically linked, not stripped"
matrix["build/pbflint.linux.bin"]="ELF 64-bit LSB executable, x86-64, version 1 (SYSV), statically linked, not stripped"
matrix["build/pbflint.exe"]="PE32 executable (console) Intel 80386 (stripped to external PDB), for MS Windows"

function checkFiles() {
  for path in "${!matrix[@]}"
  do
    expected="$path: ${matrix[$path]}";
    actual=$(file $path);
    if [ "$actual" != "$expected" ]; then
      echo "invalid file architecture: $path"
      echo "expected: $expected"
      echo "actual: $actual"
      exit 1
    fi
  done
}

checkFiles
