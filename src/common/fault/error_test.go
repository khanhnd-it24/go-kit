package fault

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getUserRepo(username string, found bool, err error) (string, error) {
	if err != nil {
		return "", Wrap(err).Tag(TagInternal)
	}
	if !found {
		return "", Wrap(fmt.Errorf("not found user")).Tag(TagNotFound)
	}

	return username, nil
}

func TestHasTag(t *testing.T) {
	_, errInternal := getUserRepo("test", false, errors.New("db crash"))
	_, errNotFound := getUserRepo("test", false, nil)

	assert.True(t, IsTag(errInternal, TagInternal))
	assert.True(t, IsTag(errNotFound, TagNotFound))
}

func TestWrap(t *testing.T) {
	original := errors.New("db crash")

	repoErr := Wrap(fmt.Errorf("failed in repos: %w", original)).Tag(TagAlreadyExists)

	bizErr := Wrap(fmt.Errorf("failed in biz: %w", repoErr)).Tag(TagUnauthenticated)

	assert.Equal(t, "failed in biz: failed in repos: db crash", bizErr.Error())
}
