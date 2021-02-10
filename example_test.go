package feign_test

import (
	"encoding/json"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/icrowley/fake"
	"github.com/mbranch/feign"
)

func Example_struct() {
	type person struct {
		Name   string
		Age    int
		Skills map[string]bool
	}
	var p person
	feign.Seed(0)
	feign.MustFill(&p)

	output(p)
	// Output:
	// {
	//   "Name": "sVedgmJqWUdRj",
	//   "Age": 22264,
	//   "Skills": {
	//     " DMYkESUcXArFAGg": true,
	//     "HbMLpzQAnthAG": false,
	//     "IgFoyZenboACW": true,
	//     "osB": false,
	//     "pSb": true
	//   }
	// }
}

func Example_fillers() {
	type customer struct {
		ID       uuid.UUID
		Email    string
		Disabled bool
	}
	var c customer
	fake.Seed(0)
	feign.Seed(0)
	feign.MustFill(&c, func(path string) (interface{}, bool) {
		switch path {
		case ".Email":
			return fake.EmailAddress(), true
		default:
			return nil, false
		}
	})

	output(c)
	// Output:
	// {
	//   "ID": "fa12f92a-fbe0-0f85-08d0-e83bab9cf8ce",
	//   "Email": "TeresaMiller@Zazio.edu",
	//   "Disabled": true
	// }
}

func Example_int() {
	var i int
	feign.Seed(0)
	feign.MustFill(&i)

	output(i)
	// Output:
	// 12282
}

func Example_string() {
	var s string
	feign.Seed(0)
	feign.MustFill(&s)

	output(s)
	// Output:
	// "sVedgmJqWUdRj"
}

func Example_nested() {
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
}

// output prints the json value to stdout.
func output(v interface{}) {
	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "  ")
	_ = e.Encode(v)
}
