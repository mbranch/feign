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
	out(p)
	// Output:
	// {
	//   "Name": "sVedgmJqWUdRjXZQpSbJeosBNHIgFoyZenboACWMGHbMLpzQAnthAGo",
	//   "Age": 1689,
	//   "Skills": {
	//     "MYkESUcXArFAGg osYWgjSrQYOqTVzK MrHZpXOdzxVHmRMcohl": false,
	//     "hlFqUvrJQ": false,
	//     "lqRffoWblyuKDyZDTkxgDZincDwMDpVLNLkWOo": true
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
	feign.MustFill(&c, func(path string) (bool, interface{}) {
		switch path {
		case ".Email":
			return true, fake.EmailAddress()
		default:
			return false, nil
		}
	})
	out(c)
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
	out(i)
	// Output:
	// 12282
}

func Example_string() {
	var s string
	feign.Seed(0)
	feign.MustFill(&s)
	out(s)
	// Output:
	// "sVedgmJqWUdRjXZQpSbJeosBNHIgFoyZenboACWMGHbMLpzQAnthAGo"
}

func Example_nested() {
	type OrderItem struct {
		ID         uuid.UUID
		Name       string
		PriceCents int
	}

	type Order struct {
		ID      uuid.UUID
		Items   []OrderItem
		Created time.Time
	}

	var o Order
	feign.Seed(1)
	feign.MustFill(&o, func(path string) (bool, interface{}) {
		switch path {
		case ".Items.PriceCents":
			return true, 100 + ((feign.Rand().Int63() % 400) * 25)
		default:
			return false, nil
		}
	})
	out(o)
	// Output:
	// {
	//   "ID": "210fc7bb-8186-39ac-48a4-c6afa2f1581a",
	//   "Items": [
	//     {
	//       "ID": "9525e20f-da68-927f-2b2f-f836f73578db",
	//       "Name": "ysAGsItGVGGRRDeTRPTNinYcyJwhrzeTW",
	//       "PriceCents": 7075
	//     }
	//   ],
	//   "Created": "0287-02-13T15:47:24.282286269Z"
	// }
}

func out(v interface{}) {
	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "  ")
	if err := e.Encode(v); err != nil {
		panic(err)
	}
}
