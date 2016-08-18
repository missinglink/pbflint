
pbflint scans the contents of an openstreetmap PBF extract, checking for errors and optimizations.

### Builds

You can download binary the files in `./build` for your architecture:

| linux x64 | osx x64 | windows x64 | arm6 |
|:-:|:-:|:-:|:-:|
| [pbflint.linux.bin](https://github.com/missinglink/pbflint/blob/master/build/pbflint.linux.bin?raw=true) | [pbflint.osx.bin](https://github.com/missinglink/pbflint/blob/master/build/pbflint.osx.bin?raw=true) | [pbflint.exe](https://github.com/missinglink/pbflint/blob/master/build/pbflint.exe?raw=true) | [pbflint.arm6.bin](https://github.com/missinglink/pbflint/blob/master/build/pbflint.arm6.bin?raw=true) |

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

All lint output is sent to `stdout`, debug info and statistics are sent to `stderr`:

```
$ ./pbflint.linux.bin /data/extract/new-york_new-york.osm.pbf 1>/dev/null
ErrorCount: 219
WarningCount: 175311
TotalNodes: 9191661
TotalWays: 1512199
TotalRelations: 8145
exit status 1
```

The linter will `exit(1)` if the ErrorCount is greater than 1, it will `exit(2)` on file errors and `exit(0)` otherwise.

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
