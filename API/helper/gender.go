package helper

import (
	"strconv"
	"strings"
	"time"
)

func Gender(gender int) string {
	switch gender {
	case 1:
		return "ប្រុស"
	case 2:
		return "ស្រី"
	default:
		return ""
	}
}

func GenderKhftEn(gender int) (genderKh string, genderEn string) {
	switch gender {
	case 1:
		return "ប្រុស", "Male"
	case 2:
		return "ស្រី", "Female"
	default:
		return "", ""
	}
}

func National(a string) string {
	switch a {
	case "ខ្មែរ":
		return "Khmer"
	default:
		return ""
	}
}

func ConvertToKhmerNumber(input string) string {
	var result strings.Builder
	//បង្កើតកន្លែងផ្ទុក String
	for _, char := range input {
		switch char {
		case '0':
			result.WriteRune('០')
			//បន្ថែមតួអក្សរមួយទៅក្នុង String។
		case '1':
			result.WriteRune('១')
		case '2':
			result.WriteRune('២')
		case '3':
			result.WriteRune('៣')
		case '4':
			result.WriteRune('៤')
		case '5':
			result.WriteRune('៥')
		case '6':
			result.WriteRune('៦')
		case '7':
			result.WriteRune('៧')
		case '8':
			result.WriteRune('៨')
		case '9':
			result.WriteRune('៩')
		default:
			result.WriteRune(char)
		}
	}
	return result.String()
}

func ConvertDateToKhmer(date string) string {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return ""
	}

	months := []string{
		"មករា", "កុម្ភៈ", "មីនា", "មេសា",
		"ឧសភា", "មិថុនា", "កក្កដា", "សីហា",
		"កញ្ញា", "តុលា", "វិច្ឆិកា", "ធ្នូ",
	}

	day := ConvertToKhmerNumber(strconv.Itoa(t.Day()))
	// strocnv.Itoa បម្លែងលេខប្រភេទ int ទៅជា string
	year := ConvertToKhmerNumber(strconv.Itoa(t.Year()))

	return day + " " + months[int(t.Month())-1] + " " + year
}
