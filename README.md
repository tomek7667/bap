# bap

bap or bump-and-push, is a small utility tool written in go that simiplifies git tag versioning

## Installation

```bash
go install github.com/tomek7667/bap@latest
```

## Usage

To see the usage:

```bash
bap --help
```

To git tag and push patch version:

```bash
bap
```

To not actuall call any git changing commands:

```bash
bap -dry
```

To bump major version:

```bash
bap -b major
```

To bump minor version:

```bash
bap -b minor
```

To bump patch version _(Although this one is default)_:

```bash
bap -b patch # or `bap`
```
