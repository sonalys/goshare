package domain

type Form FieldErrors

func (e *Form) Close() error {
	if len(*e) == 0 {
		return nil
	}

	return FieldErrors(*e)
}

func (e *Form) Append(err FieldError) {
	*e = append(*e, err)
}
