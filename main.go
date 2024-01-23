package main

func main() {
	s := NewStubRepository()
	r := NewRouter(s)
	_ = r.Run()
}
