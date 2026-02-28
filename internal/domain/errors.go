package domain

import "fmt"

type DomainError struct {
	Code string
	Message string
}

func (e *DomainError) Error() string {
    return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

var (
	ErrSourceNotFound = &DomainError{Code: "SOURCE_NOT_FOUND", Message: "Source not found"}
	ErrAlbumNotFound  = &DomainError{Code: "ALBUM_NOT_FOUND", Message: "Album not found"}
	ErrPhotoNotFound  = &DomainError{Code: "PHOTO_NOT_FOUND", Message: "Photo not found"}
)
