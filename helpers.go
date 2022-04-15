package main

func check(e *error) bool {
	if *e != nil {
		// TODO probably don't panic so easily
		panic(e)
	}
	return false
}
