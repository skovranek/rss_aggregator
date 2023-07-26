package auth

import (
	"fmt"
	"net/http"
    "strings"
	"testing"
)

func TestGetAPIKey(t *testing.T) {
	tests := []struct {
		key    string
		value  string
		expect string
        expectErr string
	}{
        {
            expectErr: "no authorization header",
        },
        {
            key: "Authorization",
            expectErr: "no authorization header",
        },
        {
            key: "Authorization",
            value: "-",
            expectErr: "incorrect formatting of authorization header",
        },
        {
            key: "Authorization",
            value: "ApiKey ",
            expectErr: "no api key in authorization header",
        },
        {
            key:    "Authorization",
            value:  "Bearer xxxxxx",
            expectErr: "incorrect formatting of authorization header",
        },
		{
			key:    "Authorization",
			value:  "ApiKey xxxxxx",
			expect: "xxxxxx",
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("Test Case #%v:", i), func(t *testing.T) {
			header := http.Header{}
			header.Add(test.key, test.value)

			output, err := GetAPIKey(header)
			if err != nil {
                if strings.Contains(err.Error(), test.expectErr) {
                    return
                }
				t.Errorf("Unexpected: %v\n", err)
				return
			}

			if output != test.expect {
				t.Errorf("Unexpected: %s", output)
				return
			}
		})
	}
}