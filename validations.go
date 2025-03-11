//
// Copyright 2023 The Eu Vat Number Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.
// Last Modification: 2025-03-11 18:22:01
//
// References:
// https://pt.wikipedia.org/wiki/Número_de_identificação_fiscal#Exemplo_de_validação_em_Go_[6]
//

package eu_vat_number

import (
	"strings"
)

func isValidESVatNumber(number string) bool {
	if len(number) != 9 {
		return false
	}

	if !(number[0] < 'A' || number[0] > 'Z') &&
		stringIsNumeric(number[1:]) { // // /^(ES)([A-Z]\d{8})$/
		return true

	}
	if ((number[0] >= 'A' || number[0] <= 'H') ||
		(number[0] >= 'N' || number[0] <= 'S') ||
		number[0] == 'W') &&
		stringIsNumeric(number[1:8]) &&
		!(number[8] < 'A' || number[8] > 'J') { // /^(ES)([A-H|N-S|W]\d{7}[A-J])$/
		return true

	}
	if ((number[0] >= 'A' || number[0] <= 'H') ||
		(number[0] >= 'N' || number[0] <= 'S') ||
		number[0] == 'W') &&
		stringIsNumeric(number[1:8]) &&
		!(number[8] < 'A' || number[8] > 'Z') { // /^(ES)([0-9|Y|Z]\d{7}[A-Z])$/
		return true

	}
	if strings.ContainsAny(string(number[0]), "KLMX") &&
		stringIsNumeric(number[1:8]) &&
		!(number[8] < 'A' || number[8] > 'Z') { // /^(ES)([K|L|M|X]\d{7}[A-Z])$/
		return true
	}

	return false
}

func isValidIEVatNumber(number string) bool {
	if len(number) < 8 || len(number) > 9 {
		return false
	}

	if stringIsNumeric(number[:7]) &&
		!(number[7] < 'A' || number[7] > 'W') { // /^(IE)(\d{7}[A-W])$/
		return true
	}
	if !(number[0] < '7' || number[0] > '9') &&
		!(number[1] < 'A' || number[1] > 'Z') &&
		stringIsNumeric(number[2:7]) &&
		!(number[7] < 'A' || number[7] > 'W') { // /^(IE)([7-9][A-Z\*\+)]\d{5}[A-W])$/
		return true
	}
	if stringIsNumeric(number[:7]) &&
		!(number[7] < 'A' || number[7] > 'W') &&
		(number[8] == 'A' || number[8] == 'H') { // /^(IE)(\d{7}[A-W][AH])$/
		return true
	}
	return false
}

// isValidPTVatNumber validates a Portuguese VAT number (NIF).
func isValidPTVatNumber(number string) bool {
	// Check if the number is exactly 9 digits
	if len(number) != 9 || !stringIsNumeric(number) {
		return false
	}

	// Tax Identification Number (NIF "123") is a service that allows the registration of a
	// individual citizen. The initial numbers "45" correspond to non-resident citizens
	// who only obtain income subject to definitive withholding tax in Portuguese territory;

	// Validate prefixes
	validPrefixes := map[string]bool{
		"45": true, "70": true, "71": true, "72": true, "74": true,
		"75": true, "77": true, "78": true, "79": true, "90": true,
		"91": true, "98": true, "99": true,
	}
	if !strings.ContainsAny(number[:1], "123568") && !validPrefixes[number[:2]] {
		return false
	}

	// Calculate check digit
	sum := 0
	for i, char := range number[:8] {
		value := int(char - '0') // Convert digit rune to integer
		sum += value * (9 - i)   // Apply weights 9 to 2
	}

	rmd := sum % 11
	var ckd int
	if rmd == 0 || rmd == 1 {
		ckd = 0
	} else {
		ckd = 11 - rmd
	}

	// Compare calculated check digit with the provided one
	compare := int(number[8] - '0')
	return compare == ckd
}
