package main

type StubAuth struct{}

func NewStubAuth() StubAuth {
	return StubAuth{}
}

func (auth StubAuth) Verify(accessToken string) (uint64, error) {
	if accessToken == "token-1" {
		return 1, nil
	}
	if accessToken == "token-2" {
		return 2, nil
	}
	if accessToken == "token-3" {
		return 3, nil
	}
	if accessToken == "token-4" {
		return 4, nil
	}
	return 0, ErrInvalidAccessToken{}
}
