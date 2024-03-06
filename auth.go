package main

type Auth interface {
	Verify(accessToken string) (uint64, error)
}
