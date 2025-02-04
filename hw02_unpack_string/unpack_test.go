package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "🙃0", expected: ""},
		{input: "aaф0b", expected: "aab"},
		// uncomment if task with asterisk completed.
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\5`, expected: `qwe\5`},
		{input: `qwe\\\3`, expected: `qwe\3`},
		// мои тесты.
		{input: `\53w4e`, expected: `555wwwwe`},
		{input: `a9b9c9d9`, expected: `aaaaaaaaabbbbbbbbbcccccccccddddddddd`},
		{input: `\15t\4\5`, expected: `11111t45`},
		{input: `t0a\40\55g0`, expected: `a55555`},
		{input: "d\n5abc", expected: "d\n\n\n\n\nabc"},
		{input: "d\t3ab0c", expected: "d\t\t\tac"},
		{input: "日3本語2", expected: "日日日本語語"},
		{input: `😈2😅1😎😇0ен\5`, expected: `😈😈😅😎ен5`},
		{input: "🙉4🙊0🙋\r3", expected: "🙉🙉🙉🙉🙋\r\r\r"},
		{input: "ジⅥ✌1💬\v2", expected: "ジⅥ✌💬\v\v"},
		{input: " 4💬5\f3", expected: "    💬💬💬💬💬\f\f\f"},
		{input: "f2 3ds", expected: "ff   ds"},
		{input: " ", expected: " "},
		{input: "\"3f0g4", expected: "\"\"\"gggg"},
		{input: `\\4`, expected: `\\\\`},
		{input: `\4`, expected: `4`},
		{input: `d3-5`, expected: `ddd-----`},
		{input: "п2р3с4", expected: "ппрррсссс"},
		{input: "!@4#2", expected: "!@@@@##"},
		{input: `\4\5abc3`, expected: "45abccc"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackError(t *testing.T) {
	tests := []struct {
		input       string
		expectedErr error
	}{
		{input: "3abc", expectedErr: ErrStartsWithDigit},
		{input: "45", expectedErr: ErrStartsWithDigit},
		{input: "aaa10b", expectedErr: ErrConsecutiveDigits},
		{input: `qw\ne`, expectedErr: ErrInvalidString},
		{input: `dhg43d\`, expectedErr: ErrBackSlashEndString},
		{input: `"日3\本2"`, expectedErr: ErrInvalidString},
		{input: `2😅1😎😇0ен`, expectedErr: ErrStartsWithDigit},
		{input: `😅0ен\`, expectedErr: ErrBackSlashEndString},
		{input: `\`, expectedErr: ErrBackSlashEndString},
		{input: `3\`, expectedErr: ErrStartsWithDigit},
		{input: `d3\-5`, expectedErr: ErrInvalidString},
		{input: string([]byte{0xFF, 0xFF}), expectedErr: ErrStringIsNotUtf8},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			_, err := Unpack(tc.input)
			require.Truef(t, errors.Is(err, tc.expectedErr), "actual error %q", err)
		})
	}
}
