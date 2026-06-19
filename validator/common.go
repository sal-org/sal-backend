package validator

func Required(value string, field string) (string, bool) {

	if value == "" {
		return field + " is required", false
	}

	if len(value) > 40 {
		return field + " is too long", false
	}

	return value, true
}
