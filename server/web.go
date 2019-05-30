package server

type Web struct {
	domain string
	key    string
	cert   string
}

//TODO: Watch key+cert for changes and automatically reload.
//TODO: Let's Encrypt support.
func (s *Server) AddWebServer(domain, cert, key string) {

}
