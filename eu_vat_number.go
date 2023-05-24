//
// Copyright 2023 The Eu Vat Number Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.
// Last Modification: 2023-05-24 10:32:17
//
// References:
// https://gist.github.com/gaiqus/4e3ee860e61f9667e911d01d816c020e#file-eu_vat_number_validation
// https://www.fonoa.com/blog/eu-vat-number-formats
// https://www.oreilly.com/library/view/regular-expressions-cookbook/9781449327453/ch04s21.html
//

package eu_vat_number

import (
	"errors"
	"fmt"
	"strings"
)

type CountryVatNumber struct {
	CountryName  string
	VATLocalName string
	Format       string
	Example      string
	EUMember     bool
}

type EuroVatNumber struct {
	ISOCountryCodes    map[string]CountryVatNumber
	DefaultCountryCode string
}

func stringIsNumeric(str string) bool {
	for _, char := range str {
		if char < '0' || char > '9' {
			return false
		}
	}

	return len(str) > 0
}

func (evt EuroVatNumber) CountryCodeAvailable(countryCode string) bool {
	if _, ok := evt.ISOCountryCodes[countryCode]; ok {
		return true
	}
	return false
}

func (evt EuroVatNumber) IsValid(id string, options ...string) (bool, error) {

	if len(id) == 0 {
		return false, errors.New("empty vat number")
	}

	if len(options) > 1 {
		return false, errors.New("wrong number of parameters")
	}

	isoCountryCode, number := func() (string, string) {
		if len(options) == 1 {
			return strings.ToUpper(options[0]), strings.ToUpper(id)
		}

		if evt.DefaultCountryCode != "" {
			return evt.DefaultCountryCode, strings.ToUpper(id)
		}

		return strings.ToUpper(id[:2]), strings.ToUpper(id[2:])
	}()

	if len(isoCountryCode) == 0 {
		return false, errors.New("empty country code")
	}

	if len(isoCountryCode) < 2 || len(isoCountryCode) > 3 {
		return false, errors.New("incorrect country code")
	}

	if _, ok := evt.ISOCountryCodes[isoCountryCode]; !ok {
		return false, errors.New("vat number not available")
	}

	switch isoCountryCode {
	case "AT": // Austria
		// /^(AT)U(\d{8})$/
		if len(number) == 9 &&
			number[0] == 'U' &&
			stringIsNumeric(number[1:]) {
			return true, nil
		}

	case "BE": // Belgium
		// /(BE)(0?\d{9})$/
		if len(number) == 9 &&
			stringIsNumeric(number) {
			return true, nil
		}
		if len(number) == 10 &&
			number[0] == '0' &&
			stringIsNumeric(number[1:]) {
			return true, nil
		}

	case "BG": // Bulgaria
		// /(BG)(\d{9,10})$/
		if len(number) >= 9 && len(number) <= 10 &&
			stringIsNumeric(number) {
			return true, nil
		}

	case "CHE": // Switzerland
		// /(CHE)(\d{9})(MWST)?$/
		if len(number) == 9 &&
			stringIsNumeric(number) {
			return true, nil
		}
		if len(number) == 13 &&
			stringIsNumeric(number[:9]) &&
			number[9:13] == "MWST" {
			return true, nil
		}

	case "CY": // Cyprus
		// /^(CY)([0-5|9]\d{7}[A-Z])$/
		if len(number) == 9 &&
			strings.ContainsAny(string(number[0]), "0123459") &&
			stringIsNumeric(number[1:8]) && !(number[8] < 'A' || number[8] > 'Z') {
			return true, nil
		}

	case "CZ": // Czech Republic
		// /^(CZ)(\d{8,10})(\d{3})?$/
		if len(number) >= 8 && len(number) <= 13 &&
			stringIsNumeric(number) {
			return true, nil
		}

	case "DE": // Germany
		// /^(DE)([1-9]\d{8})/
		if len(number) == 9 &&
			!(number[0] < '1' || number[0] > '9') &&
			stringIsNumeric(number[1:]) {
			return true, nil
		}

	case "DK": // Denmark
		// /^(DK)(\d{8})$/
		if len(number) == 8 &&
			stringIsNumeric(number) {
			return true, nil
		}

	case "EE": // Estonia
		// /^(EE)(10\d{7})$/
		if len(number) == 9 &&
			number[0] == '1' &&
			number[1] == '0' &&
			stringIsNumeric(number[2:]) {
			return true, nil
		}

	case "EL": // Greece
		// /^(EL)(\d{9})$/
		if len(number) == 9 &&
			stringIsNumeric(number) {
			return true, nil
		}

	case "ES": // Spain
		if len(number) != 9 {
			return false, nil
		}

		if !(number[0] < 'A' || number[0] > 'Z') &&
			stringIsNumeric(number[1:]) { // // /^(ES)([A-Z]\d{8})$/
			return true, nil

		}
		if ((number[0] >= 'A' || number[0] <= 'H') ||
			(number[0] >= 'N' || number[0] <= 'S') ||
			number[0] == 'W') &&
			stringIsNumeric(number[1:8]) &&
			!(number[8] < 'A' || number[8] > 'J') { // /^(ES)([A-H|N-S|W]\d{7}[A-J])$/
			return true, nil

		}
		if ((number[0] >= 'A' || number[0] <= 'H') ||
			(number[0] >= 'N' || number[0] <= 'S') ||
			number[0] == 'W') &&
			stringIsNumeric(number[1:8]) &&
			!(number[8] < 'A' || number[8] > 'Z') { // /^(ES)([0-9|Y|Z]\d{7}[A-Z])$/
			return true, nil

		}
		if strings.ContainsAny(string(number[0]), "KLMX") &&
			stringIsNumeric(number[1:8]) &&
			!(number[8] < 'A' || number[8] > 'Z') { // /^(ES)([K|L|M|X]\d{7}[A-Z])$/
			return true, nil
		}

	case "EU": // EU type
		// /^(EU)(\d{9})$/
		if len(number) == 9 &&
			stringIsNumeric(number) {
			return true, nil
		}

	case "FI": // Finland
		// /^(FI)(\d{8})$/
		if len(number) == 8 &&
			stringIsNumeric(number) {
			return true, nil
		}

	case "GR": // Greece
		// /^(GR)(\d{8,9})$/
		if len(number) >= 8 && len(number) <= 9 &&
			stringIsNumeric(number) {
			return true, nil
		}
	case "HR": // Croatia
		// /^(HR)(\d{11})$/
		if len(number) == 11 &&
			stringIsNumeric(number) {
			return true, nil
		}

	case "HU": // Hungary
		// /^(HU)(\d{8})$/
		if len(number) == 8 &&
			stringIsNumeric(number) {
			return true, nil
		}

	case "IE": // Ireland
		if len(number) == 8 {
			return false, nil
		}

		if stringIsNumeric(number[:7]) &&
			!(number[7] < 'A' || number[7] > 'W') { // /^(IE)(\d{7}[A-W])$/
			return true, nil
		}
		if !(number[0] < '7' || number[0] > '9') &&
			!(number[1] < 'A' || number[1] > 'Z') &&
			stringIsNumeric(number[2:7]) &&
			!(number[7] < 'A' || number[7] > 'W') { // /^(IE)([7-9][A-Z\*\+)]\d{5}[A-W])$/
			return true, nil
		}
		if stringIsNumeric(number[:7]) &&
			!(number[7] < 'A' || number[7] > 'W') &&
			(number[8] == 'A' || number[8] == 'H') { // /^(IE)(\d{7}[A-W][AH])$/
			return true, nil
		}

	case "IT": // Italy
		// /^(IT)(\d{11})$/
		if len(number) == 11 &&
			stringIsNumeric(number) {
			return true, nil
		}
	case "LV": // Latvia
		// /^(LV)(\d{11})$/
		if len(number) == 11 &&
			stringIsNumeric(number) {
			return true, nil
		}

	case "LT": // Lithunia
		// /^(LT)(\d{9}|\d{12})$/
		if (len(number) == 9 || len(number) == 12) &&
			stringIsNumeric(number) {
			return true, nil
		}

	case "LU": // Luxembourg
		// /^(LU)(\d{8})$/
		if len(number) == 8 &&
			stringIsNumeric(number) {
			return true, nil
		}

	case "MT": // Malta
		// /^(MT)([1-9]\d{7})$/
		if len(number) == 8 &&
			!(number[0] < '1' || number[0] > '9') &&
			stringIsNumeric(number[1:]) {
			return true, nil
		}

	case "NL": // Netherland
		// /^(NL)(\d{9})B\d{2}$/
		if len(number) == 12 &&
			stringIsNumeric(number[:9]) &&
			number[9] == 'B' &&
			stringIsNumeric(number[10:12]) {
			return true, nil
		}

	case "NO": // Norway
		// /^(NO)(\d{9})$/
		if len(number) == 9 &&
			stringIsNumeric(number) {
			return true, nil
		}

	case "PL": // Poland
		// /^(PL)(\d{10})$/
		if len(number) == 10 &&
			stringIsNumeric(number) {
			return true, nil
		}

	case "PT": // Portugal
		// /^(PT)(\d{9})$/
		if len(number) == 9 &&
			stringIsNumeric(number) {
			return true, nil
		}

	case "RO": // Romania
		// /^(RO)([1-9]\d{1,9})$/
		if len(number) >= 2 && len(number) <= 10 &&
			!(number[0] < '1' || number[0] > '9') &&
			stringIsNumeric(number[1:]) {
			return true, nil
		}

	case "RS": // Serbia
		// /^(RS)(\d{9})$/
		if len(number) == 9 &&
			stringIsNumeric(number) {
			return true, nil
		}

	case "SI": // Slovenia
		// /^(SI)([1-9]\d{7})$/
		if len(number) == 8 &&
			!(number[0] < '1' || number[0] > '9') &&
			stringIsNumeric(number[1:]) {
			return true, nil
		}

	case "SK": // Slovak Republic
		// /^(SK)([1-9]\d[(2-4)|(6-9)]\d{7})$/
		if len(number) == 10 &&
			!(number[0] < '1' || number[0] > '9') &&
			!(number[1] < '0' || number[1] > '9') &&
			(!(number[2] < '2' || number[2] > '4') || !(number[2] < '6' || number[2] > '9')) &&
			stringIsNumeric(number[3:]) {
			return true, nil
		}

	case "SE": // Sweden
		// /^(SE)(\d{10}01)$/
		if len(number) == 12 &&
			stringIsNumeric(number) {
			return true, nil
		}

	default:
		return false, nil
	}

	return false, nil
}

func New(options ...string) (*EuroVatNumber, error) {

	if len(options) > 1 {
		return nil, errors.New("wrong number of parameters")
	}

	euroVatNumber := new(EuroVatNumber)

	euroVatNumber.ISOCountryCodes = map[string]CountryVatNumber{
		"AT": {
			CountryName:  "Austria",
			Format:       "8 digits and the prefix 'U'",
			VATLocalName: "USt",
			Example:      "ATU12345678",
			EUMember:     true,
		},
		"BE": {
			CountryName:  "Belgium",
			Format:       "8 digits and the prefix 'U'",
			VATLocalName: "USt",
			Example:      "BE1234567890",
			EUMember:     true,
		},
		"BG": {
			CountryName:  "Bulgaria",
			Format:       "9 or 10 digits",
			VATLocalName: "ДДС",
			Example:      "BG123456789, BG1234567890",
			EUMember:     true,
		},
		"CHE": {
			CountryName: "Switzerland",
			EUMember:    false,
		},
		"CY": {
			CountryName:  "Cyprus",
			Format:       "8 digits and 1 letter",
			VATLocalName: "ΦΠΑ",
			Example:      "CY12345678A",
			EUMember:     true,
		},
		"CZ": {
			CountryName:  "Czech Republic",
			Format:       "8, 9 or 10 digits",
			VATLocalName: "DPH",
			Example:      "CZ12345678, CZ123456789, CZ1234567890",
			EUMember:     true,
		},
		"DE": {
			CountryName:  "Germany",
			Format:       "9 digits",
			VATLocalName: "MwSt",
			Example:      "DE123456789",
			EUMember:     true,
		},
		"DK": {
			CountryName:  "Denmark",
			Format:       "8 digits",
			VATLocalName: "Momsloven",
			Example:      "DK12345678",
			EUMember:     true,
		},
		"EE": {
			CountryName:  "Estonia",
			Format:       "9 digits",
			VATLocalName: "Käibemaks",
			Example:      "EE123456789",
			EUMember:     true,
		},
		"EL": {
			CountryName:  "Greece",
			Format:       "9 digits",
			VATLocalName: "ΦΠΑ",
			Example:      "EL123456789",
			EUMember:     true,
		},
		"ES": {
			CountryName:  "Spain",
			Format:       "9 characters - first and/or last character can be a letter",
			VATLocalName: "ESX12345678,ES12345678X, ESX1234567X",
			Example:      "IVA",
			EUMember:     true,
		},
		"EU": {
			CountryName: "EU type",
		},
		"FI": {
			CountryName:  "Finland",
			Format:       "8 digits",
			VATLocalName: "Moms",
			Example:      "FI12345678",
			EUMember:     true,
		},
		"FR": {
			CountryName:  "France",
			Format:       "11 characters - letters can be included as first and/or second character (except letters O and I)",
			VATLocalName: "TVA",
			Example:      "FR12345678901, FRX1234567890, FR1X234567890, FRXX123456789",
			EUMember:     true,
		},
		"GB": {
			CountryName: "Great Britain",
		},
		"GR": {
			CountryName: "Greece",
			EUMember:    false,
		},
		"HR": {
			CountryName:  "Croatia",
			Format:       "11 digits",
			VATLocalName: "PDV",
			Example:      "HR12345678901",
			EUMember:     true,
		},
		"HU": {
			CountryName:  "Hungary",
			Format:       "8 digits",
			VATLocalName: "ÁFA",
			Example:      "HU12345678",
			EUMember:     true,
		},
		"IE": {
			CountryName:  "Ireland",
			Format:       "8 or 9 characters - letters can be included as last, or second and last, or last 2 characters",
			VATLocalName: "VAT",
			Example:      "IE1234567X, IE1X23456X, IE1234567XX",
			EUMember:     true,
		},
		"IT": {
			CountryName:  "Italy",
			Format:       "11 digits",
			VATLocalName: "IVA",
			Example:      "IT12345678901",
			EUMember:     true,
		},
		"LV": {
			CountryName:  "Latvia",
			Format:       "11 digits",
			VATLocalName: "PVN",
			Example:      "LV12345678901",
			EUMember:     true,
		},
		"LT": {
			CountryName:  "Lithunia",
			Format:       "9 or 12 digits",
			VATLocalName: "PVM",
			Example:      "LT123456789, LT123456789012",
			EUMember:     true,
		},
		"LU": {
			CountryName:  "Luxembourg",
			Format:       "8 digits",
			VATLocalName: "TVA",
			Example:      "LU12345678",
			EUMember:     true,
		},
		"MT": {
			CountryName:  "Malta",
			Format:       "8 digits",
			VATLocalName: "VAT",
			Example:      "MT12345678",
			EUMember:     true,
		},
		"NL": {
			CountryName:  "Netherland",
			Format:       "12 characters - the 10th character is always B",
			VATLocalName: "BTW",
			Example:      "NL123456789B01",
			EUMember:     true,
		},
		"NO": {
			CountryName: "Norway",
			EUMember:    false,
		},
		"PL": {
			CountryName:  "Poland",
			Format:       "10 digits",
			VATLocalName: "PTU",
			Example:      "PL1234567890",
			EUMember:     true,
		},
		"PT": {
			CountryName:  "Portugal",
			Format:       "9 digits",
			VATLocalName: "IVA",
			Example:      "PT123456789",
			EUMember:     true,
		},
		"RO": {
			CountryName:  "Romania",
			Format:       "From 2 to 10 digits",
			VATLocalName: "TVA",
			Example:      "RO12, RO123, RO1234, RO12345, RO123456, RO123456, RO1234567, RO12345678, RO123456789, RO1234567890",
			EUMember:     true,
		},
		"RS": {
			CountryName: "Serbia",
			EUMember:    false,
		},
		"SI": {
			CountryName:  "Slovenia",
			Format:       "8 digits",
			VATLocalName: "DDV",
			Example:      "SI12345678",
			EUMember:     true,
		},
		"SK": {
			CountryName:  "Slovak Republic",
			Format:       "10 digits",
			VATLocalName: "DPH",
			Example:      "SK1234567890",
			EUMember:     true,
		},
		"SE": {
			CountryName:  "Sweden",
			Format:       "12 digits",
			VATLocalName: "Moms",
			Example:      "SE123456789012",
			EUMember:     true,
		},
	}

	if len(options) == 1 {
		if options[0] == "" {
			return nil, errors.New("the country code is empty!")
		}

		if len(options[0]) < 2 || len(options[0]) > 3 {
			return nil, errors.New("wrong country code!")
		}

		cc := strings.ToUpper(options[0])
		if _, ok := euroVatNumber.ISOCountryCodes[cc]; !ok {
			return nil, fmt.Errorf("the caountry code %s is not available!", options[0])
		}

		euroVatNumber.DefaultCountryCode = cc
	}

	return euroVatNumber, nil
}
