package usecase

import (
	"context"
	"fmt"

	pb "authentification/internal/generated/user"
)

// CreateClient handles new client creation
func (s *AuthServiceServer) CreateClient(ctx context.Context, req *pb.ClientRequest) (*pb.ClientResponse, error) {
	s.log.Info("CreateClient called", "full_name", req.FullName)

	resp, err := s.repo.CreateClient(req)
	if err != nil {
		s.log.Error("Failed to create client", "error", err)
		return nil, fmt.Errorf("could not create client: %w", err)
	}

	s.log.Info("Client created successfully", "client_id", resp.Id)
	return resp, nil
}

// GetClient retrieves a client by their ID
func (s *AuthServiceServer) GetClient(ctx context.Context, req *pb.UserIDRequest) (*pb.ClientResponse, error) {
	s.log.Info("GetClient called", "client_id", req.Id)

	client, err := s.repo.GetClient(req)
	if err != nil {
		s.log.Error("Failed to retrieve client", "client_id", req.Id, "error", err)
		return nil, fmt.Errorf("could not retrieve client: %w", err)
	}

	s.log.Info("Client retrieved successfully", "client_id", client.Id)
	return client, nil
}

// GetClientList retrieves a list of clients based on filters
func (s *AuthServiceServer) GetListClient(ctx context.Context, req *pb.FilterClientRequest) (*pb.ClientListResponse, error) {
	s.log.Info("GetClientList called", "filters", req)

	clients, err := s.repo.GetListClient(req)
	if err != nil {
		s.log.Error("Failed to retrieve client list", "error", err)
		return nil, fmt.Errorf("could not retrieve client list: %w", err)
	}

	s.log.Info("Client list retrieved successfully", "count", len(clients.Clients))
	return clients, nil
}

// UpdateClient updates a client's information
func (s *AuthServiceServer) UpdateClient(ctx context.Context, req *pb.ClientUpdateRequest) (*pb.ClientResponse, error) {
	s.log.Info("UpdateClient called", "client_id", req.Id)

	updatedClient, err := s.repo.UpdateClient(req)
	if err != nil {
		s.log.Error("Failed to update client", "client_id", req.Id, "error", err)
		return nil, fmt.Errorf("could not update client: %w", err)
	}

	s.log.Info("Client updated successfully", "client_id", updatedClient.Id)
	return updatedClient, nil
}

// DeleteClient removes a client from the system
func (s *AuthServiceServer) DeleteClient(ctx context.Context, req *pb.UserIDRequest) (*pb.MessageResponse, error) {
	s.log.Info("DeleteClient called", "client_id", req.Id)

	resp, err := s.repo.DeleteClient(req)
	if err != nil {
		s.log.Error("Failed to delete client", "client_id", req.Id, "error", err)
		return nil, fmt.Errorf("could not delete client: %w", err)
	}

	s.log.Info("Client deleted successfully", "message", resp.Message)
	return resp, nil
}
