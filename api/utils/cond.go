package utils

// Cond is a function to write a condition
func Cond(f1, eval string, f2 interface{}) M {
	return M{"type": "cond", "f1": f1, "eval": eval, "f2": f2}
}

// And is a function to "and" multiple conditions together
func And(conds ...M) M {
	return M{"type": "and", "conds": conds}
}

// Or is a function to "or" multiple conditions together
func Or(conds ...M) M {
	return M{"type": "or", "conds": conds}
}
