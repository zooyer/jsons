# jsons- golang library

![](https://travis-ci.org/boennemann/badges.svg?branch=master)  ![](https://img.shields.io/badge/license-MIT-blue.svg)  ![](https://img.shields.io/badge/godoc-reference-blue.svg)

Golang like scripting languages using the json library.

![](https://github.com/golang/go/blob/master/doc/gopher/fiveyears.jpg?raw=true)



#### Download and Install

```shell
go get github.com/zooyer/jsons
```



#### Features

- type: Raw/Bool/Number/String/Array/Object/Value.
- compatible with standard json library.
- support orm model mapping.
- chain calls.



#### Example

```go
package main

import (
	"fmt"
	
	"github.com/zooyer/jsons"
)

func main() {
	var json = `
		{
			"a": {
				"b": {
					"c": {
						"name": "test"
					}
				}
			}
		}
	`
	value, err := jsons.Unmarshal([]byte(json))
	if err != nil {
		panic(err)
	}

	c := value.Get("a", "b", "c")
	fmt.Println("c:", c.JSONString())

	name := c.Get("name")
	fmt.Println("name:", name.String())

	c.Set("name", "abc")
	fmt.Println("json:", value.JSONString())

	json = `
		[
			{
				"name": "z1",
				"age": 15
			},
			{
				"name": "z2",
				"age": 13
			},
			{
				"name": "z3",
				"age": 19
			},
			{
				"name": "z4",
				"age": 14
			}
		]
	`

	value, err = jsons.Unmarshal([]byte(json))
	if err != nil {
		panic(err)
	}

	value.Sort(func(i, j int) bool {
		return value.Get(i).Int("age") < value.Get(j).Int("age")
	})
	fmt.Println("sort:", value.JSONString())
	value.Reverse()
	fmt.Println("reverse:", value.JSONString())
}
```

ouput:

```shell
c: {"name":"test"}
name: test
json: {"a":{"b":{"c":{"name":"abc"}}}}
sort: [{"age":13,"name":"z2"},{"age":14,"name":"z4"},{"age":15,"name":"z1"},{"age":19,"name":"z3"}]
reverse: [{"age":19,"name":"z3"},{"age":15,"name":"z1"},{"age":14,"name":"z4"},{"age":13,"name":"z2"}]
```
