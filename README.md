# cmdlog
[![Go Reference](https://pkg.go.dev/badge/github.com/TheGejr/cmdlog.svg)](https://pkg.go.dev/github.com/TheGejr/cmdlog) ![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/TheGejr/cmdlog/release.yml) ![GitHub](https://img.shields.io/github/license/TheGejr/cmdlog) ![GitHub release (with filter)](https://img.shields.io/github/v/release/TheGejr/cmdlog)





## Installation and Usage
Install this program:
```
$ go install github.com/TheGejr/cmdlog@latest
```

```
USAGE: cmdlog [OPTION...] CMD [CMD OPTION...]

$ cmdlog
```

## Example usage
```sh
$ cmdlog ls -la
total 12
drwxr-xr-x 2 user user 4096 Nov 29  2024 .
drwxr-xr-x 6 user user 4096 Nov 29  2024 ..
-rw-r--r-- 1 user user  304 Nov 29  2024 cmdlog_ls_-la.log
-rw-r--r-- 1 user user    0 Nov 29  2024 file1
-rw-r--r-- 1 user user    0 Nov 29  2024 file2
-rw-r--r-- 1 user user    0 Nov 29  2024 file3

$ cat cmdlog_ls_-la.log
ls -la

total 12
drwxr-xr-x 2 plov plov 4096 Nov 29 07:55 .
drwxr-xr-x 6 plov plov 4096 Nov 29 07:55 ..
-rw-r--r-- 1 plov plov    8 Nov 29 07:55 cmdlog_ls_-la.log
-rw-r--r-- 1 plov plov    0 Nov 29 07:55 file1
-rw-r--r-- 1 plov plov    0 Nov 29 07:55 file2
-rw-r--r-- 1 plov plov    0 Nov 29 07:55 file3
```
