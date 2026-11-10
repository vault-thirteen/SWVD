package m

const (
	Err_NotEnoughArguments     = "not enough arguments"
	Err_AuthorFolderIsNotFound = "author folder is not found"
	Err_NoPages                = "no pages"
)

func MustBeNoError(err error) {
	if err != nil {
		panic(err)
	}
}
