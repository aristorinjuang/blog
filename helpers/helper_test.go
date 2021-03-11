package helpers

import "testing"

var samples = []struct {
	email   string
	correct bool
}{
	{"florian@carrere.cc", true},
	{"support@g2email.com", true},
	{" florian@carrere.cc", false},
	{"florian@carrere.cc ", false},
	{"test@912-wrong-domain902.com", true},
	{"0932910-qsdcqozuioqkdmqpeidj8793@gemail.com", true},
	{"@gemail.com", false},
	{"test@gemail@gemail.com", false},
	{"test test@gemail.com", false},
	{" test@gemail.com", false},
	{"test@wrong domain.com", false},
	{"é&ààà@gemail.com", false},
	{"admin@busyboo.com", true},
	{"a@gemail.fi", true},
}

func TestValidateEmail(t *testing.T) {
	for _, s := range samples {
		err := ValidateEmail(s.email)
		switch {
		case err != nil && s.correct == true:
			t.Errorf(`"%s" => unexpected error: "%v"`, s.email, err)
		case err == nil && s.correct == false:
			t.Errorf(`"%s" => expected error`, s.email)
		}
	}
}

func BenchmarkValidateEmail(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, s := range samples {
			err := ValidateEmail(s.email)
			switch {
			case err != nil && s.correct == true:
				b.Errorf(`"%s" => unexpected error: "%v"`, s.email, err)
			case err == nil && s.correct == false:
				b.Errorf(`"%s" => expected error`, s.email)
			}
		}
	}
}
