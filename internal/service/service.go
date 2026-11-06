package service

import (
	"errors"
	"fmt"
	"strings"
)

// переменная morseEncode хранит таблицу для кодирования текста в морзе
var morseEncode = map[string]string{
	"А": ".-", "Б": "-...", "В": ".--", "Г": "--.",
	"Д": "-..", "Е": ".", "Ж": "...-", "З": "--..",
	"И": "..", "Й": ".---", "К": "-.-", "Л": ".-..",
	"М": "--", "Н": "-.", "О": "---", "П": ".--.",
	"Р": ".-.", "С": "...", "Т": "-", "У": "..-",
	"Ф": "..-.", "Х": "....", "Ц": "-.-.", "Ч": "---.",
	"Ш": "----", "Щ": "--.-", "Ъ": "--.--", "Ы": "-.--",
	"Ь": "-..-", "Э": "..-..", "Ю": "..--", "Я": ".-.-",
	"Ё": ".",

	"0": "-----", "1": ".----", "2": "..---", "3": "...--",
	"4": "....-", "5": ".....", "6": "-....", "7": "--...",
	"8": "---..", "9": "----.",

	".": ".-.-.-", ",": "--..--", "?": "..--..",
	"!": "-.-.--", "-": "-....-", "/": "-..-.",
	"(": "-.--.", ")": "-.--.-",
}
// переменная morseDecode хранит обратную таблицу для декодирования морзе в текст
var morseDecode = buildDecodeTable(morseEncode)

// обратная таблица для декодирования морзе в текст
func buildDecodeTable(encode map[string]string) map[string]string {
	decode := make(map[string]string, len(encode))
	for k, v := range encode {
		decode[v] = k
	}
	return decode
}

// функция isMorse проверяет, является ли входная строка кодом Морзе
func isMorse(input string) bool {
	for _, r := range input {
		if !strings.ContainsRune(".- /", r) {
			return false
		}
	}
	return true
}

// данная функция textToMorse преобразует текст в код Морзе
func textToMorse(input string) (string, error) {
	input = strings.ToUpper(input)

	words := strings.Fields(input)
	textToMorse := make([]string, 0, len(words))
	for _, word := range words {
		var builder strings.Builder
		for i, r := range word {
			morse, ok := morseEncode[string(r)]
			if !ok {
				return "", fmt.Errorf("unknown symbol: %q", r)
			}
			if i > 0 {
				builder.WriteByte(' ')
			}
			builder.WriteString(morse)
		}

		textToMorse = append(textToMorse, builder.String())
	}
	return strings.Join(textToMorse, " "), nil
}

// данная функция morseToText преобразует код Морзе в текст
func morseToText(input string) (string, error) {
	words := strings.Split(input, " / ")
	morseToText := make([]string, 0, len(words))

	for _, word := range words {
		symbols := strings.Fields(word)
		if len(symbols) == 0 {
			return "", errors.New("invalid morse code: empty word")
		}
		var builder strings.Builder

		for _, symbol := range symbols {
			s, ok := morseDecode[symbol]
			if !ok {
				return "", fmt.Errorf("unknown symbol Morse: %q", symbol)
			}
			builder.WriteString(s)
		}
		morseToText = append(morseToText, builder.String())
	}
	return strings.Join(morseToText, " "), nil
}

// функция определяет строку как текст или как код Морзе и реализует соответствующее преобразование
func Convert(input string) (string, error) {
	input = strings.TrimSpace(input)

	if input == "" {
		return "", errors.New("input string is empty")
	}

	if isMorse(input) {
		return morseToText(input)
	}

	return textToMorse(input)
}
