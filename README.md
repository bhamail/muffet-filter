muffet-filter
=============

[![build](https://github.com/bhamail/muffet-filter/workflows/test.yml/badge.svg)](https://github.com/raviqqe/muffet/actions)

Uses [muffet](https://github.com/raviqqe/muffet) to check a web site for broken links and ignore known failures.

`muffet-filter` allows you to create a file (`.muffet-filter/ignores.json`) containing link errors to be ignored.

Typical usage is to run `muffet-filter`, and copy and paste selected URL's into the `.muffet-filter/ignores.json` file.
Wash, rinse, repeat until all links are either fixed, or ignored.

Then setup `muffet-filter` to run as part of your nightly CI job.
