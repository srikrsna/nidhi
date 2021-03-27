package nidhi

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	InvalidCursor         Error = "invalid cursor"
	NotFound              Error = "not found"
	DuplicateMetadataKeys Error = "two or more metadata providers using same keys"	
)
