package models

import (
	"fmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

func (s *SongCreateRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := s.Song.Validate(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (s *SongUpdateRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if s.Song != nil {
		if err := validation.Validate(*s.Song, validation.NilOrNotEmpty); err != nil {
			res = append(res, fmt.Errorf("song: %w", err))
		}
	}

	if s.Song.GroupName != "" {
		if err := validation.Validate(s.Song.GroupName, validation.NilOrNotEmpty); err != nil {
			res = append(res, fmt.Errorf("group: %w", err))
		}
	}

	if s.Song.Link != "" {
		if err := validation.Validate(s.Song.Link, validation.NilOrNotEmpty, is.URL); err != nil {
			res = append(res, fmt.Errorf("link: %w", err))
		}
	}

	if s.Song.ReleaseDate != 0 {
		if err := validation.Validate(s.Song.ReleaseDate, validation.Min(0)); err != nil {
			res = append(res, fmt.Errorf("releaseDate: %w", err))
		}
	}

	if s.Song.SongText != "" {
		if err := validation.Validate(s.Song.SongText, validation.NilOrNotEmpty); err != nil {
			res = append(res, fmt.Errorf("text: %w", err))
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}

	return nil
}
func (s *Song) Validate(formats strfmt.Registry) error {
	var res []error

	if err := validation.Validate(s.GroupName, validation.Required); err != nil {
		res = append(res, err)
	}

	if err := validation.Validate(s.SongTitle, validation.Required); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
