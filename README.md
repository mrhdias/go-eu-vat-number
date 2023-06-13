# go-eu-vat-number
EU VAT Number Validation

This module is used to validate the EU VAT number of e-commerce customer orders.
In some cases, it uses lazy and minimal evaluation. Please check the rules applied to each country.
Fixes are welcome.

## Installation
```
go get github.com/mrhdias/go-eu-vat-number
```
## Example
```go
package main

import (
	"fmt"
	"log"

	eu_vat_number "github.com/mrhdias/go-eu-vat-number"
)

func numberWithContryCode() {
	query, _ := eu_vat_number.New()

	valid, err := query.IsValid("ES01023456M")
	if err != nil {
		log.Fatalln(err)
	}
	if valid {
		fmt.Println("The VAT number is valid.")
	} else {
		fmt.Println("The VAT number is not valid!")
	}
}

func numberSepFromContryCode() {

	query, _ := eu_vat_number.New()

	// The number separated from the country code
	valid, err := query.IsValid("01023456M", "ES")
	if err != nil {
		log.Fatalln(err)
	}
	if valid {
		fmt.Println("The VAT number is valid.")
	} else {
		fmt.Println("The VAT number is not valid!")
	}
}

func withDefaultCountryCode() {
	// with the default country code set
	query, err := eu_vat_number.New("ES")
	if err != nil {
		log.Fatalln(err)
	}

	valid, err := query.IsValid("01023456M")
	if err != nil {
		log.Fatalln(err)
	}
	if valid {
		fmt.Println("The VAT number is valid.")
	} else {
		fmt.Println("The VAT number is not valid!")
	}
}

func main() {
	numberWithContryCode()
	numberSepFromContryCode()
	withDefaultCountryCode()
}
```

