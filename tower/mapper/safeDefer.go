package mapper

func safeDefer(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
