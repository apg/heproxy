# heproxy

An HTTP -> HTTPS proxy, in the spirit of [HTTPS Everywhere](https://www.eff.org/https-everywhere), which actually just uses the same rulesets as HTTPS Everywhere.

*NOTE: This is pretty experimental. Expect there to be bugs*

## Installation

`heproxy` is currently in development, and I haven't thought about
things like "releases" yet. For the brave, a simple:

```shell
$ go get github.com/apg/heproxy/cmd/heproxy
```

will work.

## Usage

There's not much to it. If you're already using HTTPS Everywhere in
your browser, but just want to make sure your `lynx`, `curl` and
`wget` connections are secured when possible, do the following:

```shell
$ heproxy -listen :9999 &
```

Then set the `http_proxy` environment variable:

```shell
$ export http_proxy=http://localhost:9999
```

Magic will now happen whenever you use an HTTP client that supports
the `http_proxy` environment variable.

## Contributing and Feedback

It's possible that heproxy has bugs, and/or does something it
shouldn't be, and/or has lots of room for improvement. If you'd like
to fix or contribute something, please fork and submit a pull request,
or open an issue.

If you have any other feedback, feel free to email me at the below
address.

## Authors

Andrew Gwozdziewycz <web@apgwoz.com>

## Copyleft

Copyright 2014, Andrew Gwozdziewycz, web@apgwoz.com
Licensed under the GNU AGPLv3. See LICENSE for more details.

The rulesets are Copyright HTTPS Everywhere Authors under the GPLv2.
