package currencyutil

type Amounts []Amount

func (amts Amounts) Sum() (Amount, error) {
	sum := Amount{}
	for _, amt := range amts {
		err := sum.Add(amt)
		if err != nil {
			return sum, err
		}
	}
	return sum, nil
}
