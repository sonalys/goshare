package domain

type Form FieldErrorList

func (e *Form) Close() error {
	if len(*e) == 0 {
		return nil
	}
	return FieldErrorList(*e)
}

func (e *Form) Append(err FieldError) {
	*e = append(*e, err)
}
