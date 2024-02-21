muffet-filter
=============

[![CI build](https://github.com/bhamail/muffet-filter/actions/workflows/test.yaml/badge.svg)](https://github.com/bhamail/muffet-filter/actions)

Uses [muffet](https://github.com/raviqqe/muffet) to check a website for broken links and ignore known failures.

`muffet-filter` allows you to create a file (`.muffet-filter/ignores.json`) containing link errors to be ignored.

Typical usage is to run `muffet-filter`, and copy and paste selected URL's into the `.muffet-filter/ignores.json` file.
Wash, rinse, repeat until all links are either fixed, or ignored.

Then setup `muffet-filter` to run as part of your nightly CI job.

TODO:
* Investigate use of [lychee](https://github.com/lycheeverse/lychee)

Dev Notes:
---------
Local test command:

```shell
./muffet-filter -i testdata/urlErrorIgnore.json https://bhamail.github.io/picapsule/
```
