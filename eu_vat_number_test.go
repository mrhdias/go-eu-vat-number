//
// Copyright 2023 The Eu Vat Number Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.
// Last Modification: 2023-05-24 10:32:11
//

package eu_vat_number

import "testing"

func TestEuroVatNumber(t *testing.T) {

	query := New()

	if !query.CountryCodeAvailable("ES") {
		t.Fatal("country code not available")
	}

	want := true
	got, err := query.IsValid("ESA2345678J")
	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Errorf("got %t, want %t", got, want)
	}

	got, err = query.IsValid("A2345678J", "ES")
	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Errorf("got %t, want %t", got, want)
	}
}
