package fs

func (s *Plugin) MethodSignatures() []string {
	return []string{
		"Chdir(path string) string",
		"Dir() (list []string)",
	}
}