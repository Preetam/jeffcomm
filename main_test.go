package main

import (
	"testing"
)

func TestParse(t *testing.T) {
	type testCase struct {
		line     string
		expected float32
	}

	cases := []testCase{
		testCase{
			`<a href="https://jeffersoncommons.residentportal.com/resident_portal/?module=ar_payments&action=create_ar_payment_transaction&kill_session=1" class="balance-adjusted">Your Balance: <b class="green-text bold">-$1,900.56</b><span> Pay Now<i class="arrow"></i></span></a>`,
			-1900.56,
		},
		testCase{
			`<a href="https://jeffersoncommons.residentportal.com/resident_portal/?module=ar_payments&action=create_ar_payment_transaction&kill_session=1" class="balance-adjusted">Your Balance: <b class="green-text bold">$1.00</b><span> Pay Now<i class="arrow"></i></span></a>`,
			1,
		},
		testCase{
			`<a href="https://jeffersoncommons.residentportal.com/resident_portal/?module=ar_payments&action=create_ar_payment_transaction&kill_session=1" class="balance-adjusted">Your Balance: <b class="green-text bold">-$1.14</b><span> Pay Now<i class="arrow"></i></span></a>`,
			-1.14,
		},
		testCase{
			`<a href="https://jeffersoncommons.residentportal.com/resident_portal/?module=ar_payments&action=create_ar_payment_transaction&kill_session=1" class="balance-adjusted">Your Balance: <b class="green-text bold">$0.00</b><span> Pay Now<i class="arrow"></i></span></a>`,
			0,
		},
		testCase{
			`<a href="https://jeffersoncommons.residentportal.com/resident_portal/?module=ar_payments&action=create_ar_payment_transaction&kill_session=1" class="balance-adjusted">Your Balance: <b class="green-text bold">$99.11</b><span> Pay Now<i class="arrow"></i></span></a>`,
			99.11,
		},
	}

	for i, c := range cases {
		if b := parseBalance(c.line); b != c.expected {
			t.Errorf("didn't parse the correct balance for case %d. Got %v, expected %v", i, b, c.expected)
		}
	}
}
