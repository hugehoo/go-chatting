package network

type api struct {
	server *Server
}

func registerServer(server *Server) {
	a := &api{server: server}

	server.engine.Get()

}
