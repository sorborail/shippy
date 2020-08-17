// shippy-service-consignment/main.go
package main

import (
	"context"
	"github.com/micro/go-micro/v2"
	"log"

	"github.com/pkg/errors"
	// Import the generated protobuf code
	pb "github.com/sorborail/shippy/shippy-service-consignment/proto/consignment"
	vesselProto "github.com/sorborail/shippy/shippy-service-vessel/proto/vessel"
)

type repository interface {
	Create(ctx context.Context, consignment *pb.Consignment) error
	GetAll(ctx context.Context) ([]*pb.Consignment, error)
}

// Repository - Dummy repository, this simulates the use of a datastore
// of some kind. We'll replace this with a real implementation later on.
type Repository struct {
	consignments []*pb.Consignment
}

// Service should implement all of the methods to satisfy the service
// we defined in our protobuf definition. You can check the interface
// in the generated code itself for the exact method signatures etc
// to give you a better idea.

type handler struct {
	repo repository
	vesselClient vesselProto.VesselService
}

// Create a new consignment
func (repo *Repository) Create(ctx context.Context, consignment *pb.Consignment) error {
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	return nil
}

// GetAll consignments
func (repo *Repository) GetAll(ctx context.Context) ([]*pb.Consignment, error) {
	return repo.consignments, nil
}

/*// CreateConsignment - we created just one method on our service,
// which is a create method, which takes a context and a request as an
// argument, these are handled by the gRPC server.
func (s *handler) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {

	// Save our consignment
	consignment, err := s.repository.Create(req)
	if err != nil {
		return err
	}

	// Return matching the `Response` message we created in our
	// protobuf definition.
	res.Created = true
	res.Consignment = consignment
	return nil
}*/

// CreateConsignment - we created just one method on our service,
// which is a create method, which takes a context and a request as an
// argument, these are handled by the gRPC server.
func (s *handler) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {

	// Here we call a client instance of our vessel service with our consignment weight,
	// and the amount of containers as the capacity value
	log.Println("Handler create....")
	vesselResponse, err := s.vesselClient.FindAvailable(ctx, &vesselProto.Specification{
		MaxWeight: req.Weight,
		Capacity:  int32(len(req.Containers)),
	})
	if vesselResponse == nil {
		return errors.New("error fetching vessel, returned nil")
	}

	if err != nil {
		return err
	}

	// We set the VesselId as the vessel we got back from our
	// vessel service
	req.VesselId = vesselResponse.Vessel.Id

	log.Println(req.VesselId)

	// Save our consignment
	err = s.repo.Create(ctx, req)
	if err != nil {
		return err
	}

	res.Created = true
	res.Consignment = req
	return nil
}


// GetConsignments -
func (s *handler) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	consignments, _ := s.repo.GetAll(ctx)
	res.Consignments = consignments
	return nil
}

func main() {

	repo := &Repository{}

	// Create a new service. Optionally include some options here.
		service := micro.NewService(

		// This name must match the package name given in your protobuf definition
		micro.Name("shippy.service.consignment"),
	)

	// Init will parse the command line flags.
	service.Init()

	log.Println("Main in.....")
	vesselClient := vesselProto.NewVesselService("shippy.service.vessel", service.Client())
	h := &handler{repo, vesselClient}

	// Register service
	if err := pb.RegisterShippingServiceHandler(service.Server(), h); err != nil {
		log.Panic(err)
	}

	// Run the server
	if err := service.Run(); err != nil {
		log.Panic(err)
	}

	/*// Set-up our gRPC server.
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	// Register our service with the gRPC server, this will tie our
	// implementation into the auto-generated interface code for our
	// protobuf definition.
	pb.RegisterShippingServiceServer(s, &consignmentService{repo})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	log.Println("Running on port:", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}*/
}