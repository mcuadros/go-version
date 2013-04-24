go-version
==========

Version normalizer and comparison library for go













go-version [![Build Status](https://travis-ci.org/mcuadros/go-version.png?branch=master)](https://travis-ci.org/mcuadros/go-version)
==============================

Version normalizer and comparison library for go, heavy based on PHP version_compare function and Version comparsion libs from [Composer](https://github.com/composer/composer) PHP project

Installation
------------

The recommended way to install go-version

```
go get github.com/mcuadros/go-version
}
```


Examples
--------

How import the package

```go
import (
    "fmt"
    "github.com/mcuadros/go-version"
)
```

Fuction Normalize(): Normalizes a version string to be able to perform comparisons on it

```go
Normalize("10.4.13-b")
\\Returns: 10.4.13.0-beta
```


License
-------

MIT, see [LICENSE](LICENSE)