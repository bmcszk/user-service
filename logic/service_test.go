package logic

import (
	"testing"
)

func TestService_CreateUser(t *testing.T) {
	given, when, then := NewBlocks(t)
	given.aUser().and().
		dbCanCreateUser()

	when.serviceCreatesUser()

	then.noError().and().
		userIsReturned().and().
		returnedUserIsValid()
}

func TestService_CreateUser_Duplicated(t *testing.T) {
	given, when, then := NewBlocks(t)
	given.aUser().and().
		dbCannotCreateDuplicatedUser()

	when.serviceCreatesUser()

	then.returnedErrorIs(ErrUserAlreadyExists)
}

func TestService_CreateUser_InvalidUser(t *testing.T) {
	given, when, then := NewBlocks(t)
	given.aInvaliUser()

	when.serviceCreatesUser()

	then.returnedErrorIs(ErrUserNameEmpty)
}

func TestService_GetsUser(t *testing.T) {
	given, when, then := NewBlocks(t)
	given.aID().and().
		dbCanGetUser()

	when.serviceGetsUser()

	then.noError().and().
		userIsReturned().and().
		returnedUserIsValid()
}

func TestService_GetsUser_NotFound(t *testing.T) {
	given, when, then := NewBlocks(t)
	given.aID().and().
		dbCannotFindUser()

	when.serviceGetsUser()

	then.returnedErrorIs(ErrUserNotFound)
}

func TestService_UpdatesUser(t *testing.T) {
	given, when, then := NewBlocks(t)
	given.aUser().and().
		aID().and().
		dbCanUpdateUser()

	when.serviceUpdatesUser()

	then.noError().and().
		userIsReturned().and().
		returnedUserIsValid()
}

func TestService_UpdatesUser_NotFound(t *testing.T) {
	given, when, then := NewBlocks(t)
	given.aUser().and().
		aID().and().
		dbCannotFindUserForUpdate()

	when.serviceUpdatesUser()

	then.returnedErrorIs(ErrUserNotFound)
}

func TestService_UpdatesUser_ExistingUsername(t *testing.T) {
	given, when, then := NewBlocks(t)
	given.aUser().and().
		aID().and().
		dbCannotUpdateWithExistingUsername()

	when.serviceUpdatesUser()

	then.returnedErrorIs(ErrUserAlreadyExists)
}

func TestService_UpdatesUser_InvalidUser(t *testing.T) {
	given, when, then := NewBlocks(t)
	given.aInvaliUser().and().
		aID()

	when.serviceUpdatesUser()

	then.returnedErrorIs(ErrUserNameEmpty)
}

func TestService_DeletesUser(t *testing.T) {
	given, when, then := NewBlocks(t)
	given.aID().and().
		dbCanDeleteUser()

	when.serviceDeletesUser()

	then.noError()
}

func TestService_DeletesUser_NotFound(t *testing.T) {
	given, when, then := NewBlocks(t)
	given.aID().and().
		dbCannotFindUserForDelete()

	when.serviceDeletesUser()

	then.returnedErrorIs(ErrUserNotFound)
}

func TestService_ListsUsers(t *testing.T) {
	given, when, then := NewBlocks(t)
	given.dbCanListUsers()

	when.serviceListsUsers()

	then.noError().and().
		usersAreReturned().and().
		returnedUsersAreValid()
}
