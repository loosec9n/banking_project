package utils

const (
	INR = "INR"
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
	AUD = "AUD"
)

// IsSupportedCurrency return a boolean value if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case INR, USD, EUR, CAD, AUD:
		return true
	}
	return false
}
