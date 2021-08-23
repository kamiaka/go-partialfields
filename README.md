# go-partialfields


partialfields is a program for the Go language that checks that the fields of structural literals are partially defined.

#  🙅‍♀️ Don't use this.<br>Use [exhaustivestruct](https://github.com/mbilski/exhaustivestruct)

## Usage

```sh
partialfields [-flag] [package]

```

Just give the package path.

```sh
partialfields github.com/kamiaka/partialfields/testdata/src/a
```

## Control errors

## Comment

Skip if the struct literal comment starts with `// partial`.

```go
type Value struct {
  Foo, Bar int
}

// skip check.
var OK = Value{ // partial
  Foo: 1,
}

// requires Bar.
var NG = Value{
  Foo: 1,
}
```
