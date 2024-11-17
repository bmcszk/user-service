package e2e

import (
	"net/http"
	"testing"
)

func TestPost(t *testing.T) {
	given, when, then := NewBlocks(t)

	given.aValidUserData()

	when.postRequest().sending()

	then.noError().and().
		statusCodeIs(http.StatusCreated).and().
		userIsReturned().and().
		returnedUserIsValid().and().
		userIsStoredInDB()
}

func TestPost_Duplicate(t *testing.T) {
	given, when, then := NewBlocks(t)

	given.aValidUserData().alreadyStoredInDB()

	when.postRequest().sending()

	then.noError().and().
		statusCodeIs(http.StatusConflict)
	// TODO check error response
}

func TestPost_Invalid(t *testing.T) {
	given, when, then := NewBlocks(t)

	given.aInvalidUserData()

	when.postRequest().sending()

	then.noError().and().
		statusCodeIs(http.StatusBadRequest)
	// TODO check error response
}

func TestGet(t *testing.T) {
	given, when, then := NewBlocks(t)

	given.aValidUserData().alreadyStoredInDB()

	when.getRequest().sending()

	then.noError().and().
		statusCodeIs(http.StatusOK).and().
		userIsReturned().and().
		returnedUserIsValid()
}

func TestGet_NotFound(t *testing.T) {
	given, when, then := NewBlocks(t)

	given.aID()

	when.getRequest().sending()

	then.noError().and().
		statusCodeIs(http.StatusNotFound)
}

func TestPut(t *testing.T) {
	given, when, then := NewBlocks(t)

	given.aValidUserData().alreadyStoredInDB().and().
		userIsChanged()

	when.putRequest().sending()

	then.noError().and().
		statusCodeIs(http.StatusOK).and().
		userIsReturned().and().
		returnedUserIsValid().and().
		userIsStoredInDB()
}

// TODO test PUT invalid, not found

func TestDelete(t *testing.T) {
	given, when, then := NewBlocks(t)

	given.aValidUserData().
		alreadyStoredInDB()

	when.deleteRequest().sending()

	then.noError().and().
		statusCodeIs(http.StatusNoContent)
}

func TestList(t *testing.T) {
	given, when, then := NewBlocks(t)

	given.aValidUserData().alreadyStoredInDB().and().
		userIsChanged().alreadyStoredInDB().and().
		userIsChanged().alreadyStoredInDB()

	when.listRequest().sending()

	then.noError().and().
		statusCodeIs(http.StatusOK).and().
		usersAreReturned().and().
		returnedUsersAreValid()
	// TODO check every user returned
	// TODO check pagination
}
