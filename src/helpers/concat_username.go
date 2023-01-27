package helpers

func ConcatUsername(firstname string, lastname string) string {
	if lastname == "" {
		return firstname
	}
	return firstname + " " + lastname
}
