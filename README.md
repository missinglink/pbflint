
pbflint scans the contents of an openstreetmap PBF extract, checking for errors and optimizations.

### Builds

You can download binary the files in `./build` for your architecture:

| linux x64 | osx x64 | windows x64 | arm6 |
|:-:|:-:|:-:|:-:|
| build/pbflint.linux.bin | build/pbflint.osx.bin | build/pbflint.exe | build/pbflint.arm6.bin |

### Usage

Simply provide the path to the PBF extract as the first CLI argument:

```bash
$ pbflint.linux.bin /tmp/wellington_new-zealand.osm.pbf
```
```bash
error: way 93494080 invalid refcount 1
error: way 119134158 invalid refcount 1
... etc
```

### Run the go code from source

Make sure `Go` is installed and configured on your system, see: https://gist.github.com/missinglink/4212a81a7d9c125b68d9

**Note:** You should install the latest version of Golang, at least `1.5+`, last tested on `1.6.2`

```bash
go get;
go run pbflint.go /path/to/extract.pbf;
```

### Compile source for all supported architecture

If you are doing a release and would like to compile for all supported architectures:

```bash
./compile.sh;
```

```bash
[compile] linux arm
[compile] linux x64
[compile] darwin x64
[compile] windows x64
```

### License

```
This work ‘as-is’ we provide.
No warranty express or implied.
  We’ve done our best,
  to debug and test.
Liability for damages denied.

Permission is granted hereby,
to copy, share, and modify.
  Use as is fit,
  free or for profit.
These rights, on this notice, rely.
```
