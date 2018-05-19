GoGlobLS
====
_A tiny CLI tool to list files according to a yaml with a glob descriptor_

Usage
----

```sh
# Go to the test folder
$ cd test

# List all files according to the test.yaml descriptor
$ ../goglobls -config test.yaml .
a\a\a
a\b\a
a\b\b
a\c\c
a\d
b\b

# List all files according to the test.yaml descriptor filtered by type "source"
$ ../goglobls -config test.yaml -type source .
a\a\a
a\b\a
a\b\b
a\c\c
a\d

# List all files according to the test.yaml descriptor filtered by type "source" and "test"
$ ../goglobls -config test.yaml -type source -type test .
a\a\a
a\b\a
a\b\b
a\c\c
a\d
b\b
```

Check the descriptor [example](test/test.yaml).

License
----
This project is licensed under the [MIT](LICENSE.md) License.

