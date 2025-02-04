package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	ErrInvalidString      = errors.New("invalid string")
	ErrStringIsNotUtf8    = errors.New("string is not UTF-8 encoded")
	ErrBackSlashEndString = errors.New("string ends with backslash")
	ErrStartsWithDigit    = errors.New("string starts with digit")
	ErrConsecutiveDigits  = errors.New("consecutive digits not allowed")
)

// Проврека, что rune число.
func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}

// Проврека, что rune обратный слэш.
func isBackSlash(r rune) bool {
	return r == '\\'
}

// Функция Unpack.
func Unpack(in string) (string, error) {
	// Проверка строки на соответвие UTF8.
	if !utf8.ValidString(in) {
		return "", ErrStringIsNotUtf8
	}

	// Перевод строки в массив rune для последующей работы с ней.
	runes := []rune(in)
	maxLen := len(runes)

	// Проверка если пустая строка.
	if maxLen == 0 {
		return "", nil
	}
	// Проверка если первый символ число.
	if isDigit(runes[0]) {
		return "", ErrStartsWithDigit
	}

	// Инициализация в памяти переменной.
	var result strings.Builder
	result.Grow(maxLen * 2)

	for i := 0; i < maxLen; i++ {
		// Читаем текущий символ.
		curr := runes[i]

		// Выход с ошибкой если сначала идет число.
		if isDigit(curr) {
			return "", ErrConsecutiveDigits
		}

		// Проверка на обратный слэш,
		// если это так то читаем следующий символ и проверяем его,
		// что он цифра или обратный слэш, иначе выход с ошибкой.
		if isBackSlash(curr) {
			// Проверка если последний символ обратный слэш.
			if i >= maxLen-1 {
				return "", ErrBackSlashEndString
			}
			i++
			curr = runes[i]
			if !isDigit(curr) && !isBackSlash(curr) {
				return "", ErrInvalidString
			}
		}

		// Проверка, что следующий символ число,
		// Если это так, то делаем повтор символов на данное число,
		// Записываем результат и идем в начало цикла.
		if maxLen > i+1 && isDigit(runes[i+1]) {
			i++
			count, err := strconv.Atoi(string(runes[i]))
			if err != nil {
				return "", err
			}

			if _, err := result.WriteString(strings.Repeat(string(curr), count)); err != nil {
				return "", err
			}
			continue
		}

		//  Если следующий символ не число,
		// то записываем текущий символ и идем в начало цикла.
		if _, err := result.WriteRune(curr); err != nil {
			return "", err
		}
	}

	// Выход с возвратом результата.
	return result.String(), nil
}
