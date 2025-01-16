muffet-filter
=============

[![CI build](https://github.com/bhamail/muffet-filter/actions/workflows/test.yaml/badge.svg)](https://github.com/bhamail/muffet-filter/actions)

Uses [muffet](https://github.com/raviqqe/muffet) to check a website for broken links and ignore known failures.

`muffet-filter` allows you to create a file (`.muffet-filter/ignores.json`) containing link errors to be ignored.

Typical usage is to run `muffet-filter`, and copy and paste selected URL's into the `.muffet-filter/ignores.json` file.
Wash, rinse, repeat until all links are either fixed, or ignored.

Then setup `muffet-filter` to run as part of your nightly CI job.

The easiest way to get started is to copy the muffet-filter wrapper directory and script: [.muffet-filter/mfw](.muffet-filter/mfw) 
into your project, and execute that script in your CI build. See this [example](https://github.com/bhamail/picapsule/blob/c032e40186ee3c7a679d78deb83f88932d689aef/.github/workflows/link-check.yaml#L13-L13).

[CLI Help](.snapshots/TestHelp)

ignores.json syntax
-------------------
The `.muffet-filter/ignores.json` file is a JSON file containing a map of URL's and error messages to ignore. Both the
URL and the error message support regular expression matching. Some example content is shown below. Another example is the
testdata file [testdata/ci-link-check-ignores.json](testdata/ci-link-check-ignores.json).

```json
[
  {
    "url": "https://contribute.sonatype.com/js/ga.js",
    "error": "404"
  },
  {
    "url": "https://www.docker.com/products/docker-desktop",
    "error": "403"
  },
  {
    "url": "https://github.githubassets.com/",
    "error": "404"
  }
]
```

TODO:
* Investigate use of [lychee](https://github.com/lycheeverse/lychee)

Dev Notes:
---------
Local test command:

```shell
./muffet-filter -i testdata/urlErrorIgnore.json https://bhamail.github.io/picapsule/
```

Release Process
---------------
To release a new version, create a new tag with a sematic version and push it to the repo. 
CI will automatically build and publish the new version.

```shell
git tag -a v0.0.1 -m "Release 0.0.1"
git push origin v0.0.1
```

See [GoReleaser](https://goreleaser.com/quick-start/) for more details.
