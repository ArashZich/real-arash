package service

func discount(discountType string, price int, discountingAmount int) (int, error) {
	var discountedAmount int

	switch discountType {
	case "percent":
		// Calculate percentage of the price. Convert to float for accurate calculation.
		discountValue := float64(price) * (float64(discountingAmount) / 100.0)
		// Convert back to int after calculating the discount.
		discountedAmount = price - int(discountValue)
	case "fixed_amount":
		discountedAmount = price - discountingAmount
	default:
		// Handle unknown discount type.
		return 0, nil
	}

	return discountedAmount, nil
}
