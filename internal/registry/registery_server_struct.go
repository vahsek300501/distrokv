package registry

import pb "github.com/Vahsek/distrokv/pkg/registry"

type server struct {
	pb.UnimplementedRegistryServiceServer
}
