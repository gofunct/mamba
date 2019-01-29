# Fmap
 
## Overview
 
 ```text
String Functions: trim, wrap, randAlpha, plural, etc.
String List Functions: splitList, sortAlpha, etc.
Math Functions: add, max, mul, etc.
Integer Slice Functions: until, untilStep
Date Functions: now, date, etc.
Defaults Functions: default, empty, coalesce, toJson, toPrettyJson
Encoding Functions: b64enc, b64dec, etc.
Lists and List Functions: list, first, uniq, etc.
Dictionaries and Dict Functions: dict, hasKey, pluck, etc.
Type Conversion Functions: atoi, int64, toString, etc.
File Path Functions: base, dir, ext, clean, isAbs
Flow Control Functions: fail
Advanced Functions
UUID Functions: uuidv4
OS Functions: env, expandenv
Version Comparison Functions: semver, semverCompare
Reflection: typeOf, kindIs, typeIsLike, etc.
Cryptographic and Security Functions: derivePassword, sha256sum, genPrivateKey
 
 ```
 
 
## Examples

### Strings
```text
{{trim "   hello    "}}:                                                            hello
{{trimAll "$" "$5.00"}}:                                                            5.00
{{trimSuffix "-" "hello-"}}:                                                        hello
{{upper "hello"}}:                                                                  HELLO
{{lower "HELLO"}}:                                                                  hello
{{title "hello world"}}:                                                            Hello World
{{untitle "Hello World"}}:                                                          hello world
{{repeat 3 "hello"}}:                                                               hellohellohello
{{substr 0 5 "hello world"}}:                                                       hello
{{nospace "hello w o r l d"}}:                                                      helloworld
{{trunc 5 "hello world"}}:                                                          hello
{{abbrev 5 "hello world"}}:                                                         he...
{{abbrevboth 5 10 "1234 5678 9123"}}:                                               ...5678...
{{initials "First Try"}}:                                                           FT
{{randNumeric 3}}:                                                                  528
{{- /*{{wrap 80 $someText}}*/}}:
{{wrapWith 5 "\t" "Hello World"}}:                                                  Hello	World
{{contains "cat" "catch"}}:                                                         true
{{hasPrefix "cat" "catch"}}:                                                        true
{{cat "hello" "beautiful" "world"}}:                                                hello beautiful world
{{- /*{{indent 4 $lots_of_text}}*/}}:
{{- /*{{indent 4 $lots_of_text}}*/}}:
{{"I Am Henry VIII" | replace " " "-"}}:                                            I-Am-Henry-VIII
{{len .Service.Method | plural "one anchovy" "many anchovies"}}:                    many anchovies
{{snakecase "FirstName"}}:                                                          first_name
{{camelcase "http_server"}}:                                                        HttpServer
{{shuffle "hello"}}:                                                                holle
{{regexMatch "[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,}" "test@acme.com"}}:   true
{{- /*{{regexFindAll "[2,4,6,8]" "123456789"}}*/}}:
{{regexFind "[a-zA-Z][1-9]" "abcd1234"}}:                                           d1
{{regexReplaceAll "a(x*)b" "-ab-axxb-" "${1}W"}}:                                   -W-xxW-
{{regexReplaceAllLiteral "a(x*)b" "-ab-axxb-" "${1}"}}:                             -${1}-${1}-
{{regexSplit "z+" "pizza" -1}}:                                                     [pi a]

# Get one specific method on array method using index
{{ index .Service.Method 1 }}:                                                      name:"Iii" input_type:".dummy.Dummy2" output_type:".dummy.Dummy1" options:<> 

# Sprig: advanced
{{if contains "cat" "catch"}}yes{{else}}no{{end}}:   yes
{{1 | plural "one anchovy" "many anchovies"}}:       one anchovy
{{2 | plural "one anchovy" "many anchovies"}}:       many anchovies
{{3 | plural "one anchovy" "many anchovies"}}:       many anchovies

```

### Protoc

```text
# Common variables
{{.File.Name}}:                                                                           helpers.proto
{{.File.Name | upper}}:                                                                   HELPERS.PROTO
{{.File.Package | base | replace "." "-"}}                                                dummy
{{$packageDir := .File.Name | dir}}{{$packageDir}}                                        .
{{$packageName := .File.Name | base | replace ".proto" ""}}{{$packageName}}               helpers
{{$packageImport := .File.Package | replace "." "_"}}{{$packageImport}}                   dummy
{{$namespacedPackage := .File.Package}}{{$namespacedPackage}}                             dummy
{{$currentFile := .File.Name | getProtoFile}}{{$currentFile}}                             <nil>
{{- /*{{- $currentPackageName := $currentFile.GoPkg.Name}}{{$currentPackageName}}*/}}
```