package cql

func ListenAndServe(addr string) error {
	server := NewServer(addr)
	return server.ListenAndServe()
}
