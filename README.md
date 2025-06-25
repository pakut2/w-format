# W Format

> [!WARNING]
> This software has no use.

Transpile Javascript to [Whitespace](https://esolangs.org/wiki/Whitespace). Format Javascript/Typescript files according to the generated Whitespace, combining two different programs into a single source file.

## Requirements

- [Go](https://go.dev/) ver. 1.24.4 or newer

## Quickstart

```shell
make run-example
```

After running, a new file will be generated under `./examples/output.ts`. The result takes the [source Javascript file](./examples/source.js) and transpiles it to Whitespace. Afterwards, the [target file](./examples/format.ts) is formatted to include the generated Whitespace instructions, creating the final output file.

This results in a single source file, running the same program while interpreted as different programming languages.

Typescript result can be executed with:

```shell
node --experimental-strip-types examples/output.ts
```

Whitespace can be executed using an online [Whitespace interpreter](https://naokikp.github.io/wsi/whitespace.html).

## Options

Display available options:

```shell
go run cmd/jsWhitespaceFormatter/main.go -h
```

Print transpiled Whitespace to standard output:

```shell
go run cmd/jsWhitespaceFormatter/main.go -source-file=<js-file-path>
```

Save transpiled Whitespace to file:

```shell
go run cmd/jsWhitespaceFormatter/main.go -source-file=<js-file-path> -output-file=<output-file-path>
```

Format file with transpiled Whitespace, print result to standard output:

```shell
go run cmd/jsWhitespaceFormatter/main.go -source-file=<js-file-path> -format-file<format-file-path>
```

Save file formatted with transpiled Whitespace:

```shell
go run cmd/jsWhitespaceFormatter/main.go -source-file=<js-file-path> -format-file<format-file-path> -output-file=<output-file-path>
```

## Supported syntax

Not all Javascript instructions are supported by the transpiler. The covered subset includes:

- Literal types:
  - **Number**: integers only
  - **String**
  - **Boolean**
  
- Variable declarations using **let**:

```javascript
let a = 42;
let b = "Hello World";
let c = true;
```

- **Number** and **Boolean** variable reassignments:

```javascript
let a = 42;
a = 69;
```

- Printing to standard output with **console.log**:

```javascript
console.log(42, "Hello World", true);
```

- Arithmetic/logic operations:
  - **+**
  - **-**
  - **\***
  - **/**
  - **%**
  - **++**
  - **--**
  - **&&**
  - **||**
  - **!**
  - **===**
  - **!==**
  - **<**
  - **<=**
  - **>**
  - **>=**

```javascript
console.log(2 + 2 * 2);
```

- String concatenation:

```javascript
console.log("Hello" + "World");
```

- **if** / **if/else** statements:

```javascript 
if (2 > 1) {
    console.log(true);
}
```

- C-style **for** loops:

```javascript
for (let i = 0; i < 10; i++) {
    if (i === 8) {
        break;
    }
  
    if (i % 2 === 0) {
        continue;
    }
    
    console.log(i);
}
```
