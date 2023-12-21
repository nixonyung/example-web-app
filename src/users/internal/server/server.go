package server

import "source.local/users/internal/proto"

type Server struct {
	proto.UnimplementedUsersServiceServer
}
