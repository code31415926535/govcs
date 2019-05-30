package govcs

func Init(path string) error {
	vcs, err := NewDefaultVcs(path)
	if err != nil {
		return err
	}
	return vcs.Init()
}

func Stat(path string) (Status, error) {
	vcs, err := LoadDefaultVcs(path)
	if err != nil {
		return Status{}, err
	}

	return vcs.Stat()
}

func AddFile(path string) error {
	vcs, err := LoadDefaultVcs(path)
	if err != nil {
		return err
	}

	return vcs.AddFile(path)
}

func RemoveFile(path string) error {
	vcs, err := LoadDefaultVcs(path)
	if err != nil {
		return err
	}

	return vcs.RemoveFile(path)
}

func CommitChanges(path string, message string) error {
	vcs, err := LoadDefaultVcs(path)
	if err != nil {
		return err
	}

	return vcs.CommitChanges(message)
}
