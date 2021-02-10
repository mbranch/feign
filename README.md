# feign

[![go.dev](https://img.shields.io/badge/go.dev-pkg-007d9c.svg?style=flat)](https://pkg.go.dev/github.com/mbranch/feign)
[![CircleCI](https://img.shields.io/circleci/build/github/mbranch/feign)](https://circleci.com/gh/mbranch/feign/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/mbranch/feign)](https://goreportcard.com/report/github.com/mbranch/feign)

feign automatically fills types with random data.

It can be useful when testing, where you only need to verify persistence and
retrieval and aren't concerned with valid data values.

Based on https://github.com/bxcodec/faker.

## Install

```
go get -u github.com/mbranch/feign
```

## Examples

```go
// output prints the json value to stdout.
func output(v interface{}) {
  e := json.NewEncoder(os.Stdout)
  e.SetIndent("", "  ")
  _ = e.Encode(v)
}

// String example.

feign.Seed(0)
feign.MustFill(&s)

output(s)
// Output:
// "sVedgmJqWUdRj"

// Example with nested structs.

type OrderItem struct {
  ProductID       int
  Name            string
  PriceFractional int
  Attributes      map[string]string
}

type Order struct {
  ID      uuid.UUID
  Items   []OrderItem
  Created time.Time
}

var o Order
feign.Seed(1)
feign.MustFill(&o, func(path string) (interface{}, bool) {
  switch path {
  case ".Items.PriceFractional":
    return 100 + ((feign.Rand().Int63() % 400) * 25), true
  default:
    return nil, false
  }
})

output(o)
// Output:
// {
//   "ID": "210fc7bb-8186-39ac-48a4-c6afa2f1581a",
//   "Items": [
//     {
//       "ProductID": 23701,
//       "Name": "SCX",
//       "PriceFractional": 2875,
//       "Attributes": {
//         "AmTgVjiMDy": "AGsItGVGGRRDeTRPTNinYcyJ",
//         "EYAua wti": "NPFhIvN",
//         "IN NY": "XAnZHdKrMfWYLFocFYszCG eZj",
//         "TKvBqWJBscgSE": "IsJqDvttR",
//         "gTivDxUcOYVZwJCZbf": "PHPGGhoQQQoCFcgJCLF",
//         "hrzeTWVmkCrTDsmwcpWKwcxnzDyOyqx": "e",
//         "kdflCVbJoFXdsTMGBEdXryjTFQrd": "QW"
//       }
//     }
//   ],
//   "Created": "0073-06-09T21:30:44.650537874Z"
// }

```
