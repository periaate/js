# js
Simple CLI utility tool for using JS expression in pipes.

## Help
```
Usage:
  js [options] <arrow function>

Application Options:
  -v, --verbose   Print verbose output.
  -e, --errors    Panic on errors.
```

## Semantics
- The argument(s) given to `js` need to create a valid, anonymous arrow function `a => ...`.
- Each piped value is individually evaluated with the expression given as argument.
- Values which evaluate to empty strings are not printed.
- If given multiple arguments, arguments are joined with " ".

## Examples
Lowercase every piped value.
```
ls | js "a => a.toLowerCase()"
```

Filter away strings which contain any of the following `.git`, `.go` `.mod` as substrings.
```
ls | js 'a => [".git", ".go", ".mod"].some(b => a.includes(b)) ? "" : a'
```