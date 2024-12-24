package repo

import (
	"authentification/internal/entity"
	pb "authentification/internal/generated/user"
	"fmt"
	"strconv"
	"strings"
)

func (c *UserRepo) CreateClient(in *pb.ClientRequest) (*pb.ClientResponse, error) {
	query := `
		INSERT INTO clients (full_name, address, phone, type)
		VALUES ($1, $2, $3, $4)
		RETURNING id, full_name, address, phone, type
	`
	var client pb.ClientResponse
	err := c.db.QueryRowx(query, in.FullName, in.Address, in.Phone, in.Type).Scan(
		&client.Id,
		&client.FullName,
		&client.Address,
		&client.Phone,
		&client.Type,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	return &client, nil
}

func (c *UserRepo) GetClient(in *pb.UserIDRequest) (*pb.ClientResponse, error) {
	query := `
		SELECT id, full_name, address, phone, type
		FROM clients
		WHERE id = $1
	`
	var client pb.ClientResponse
	err := c.db.QueryRowx(query, in.Id).Scan(
		&client.Id,
		&client.FullName,
		&client.Address,
		&client.Phone,
		&client.Type,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve client: %w", err)
	}
	return &client, nil
}

func (c *UserRepo) GetListClient(in *pb.FilterClientRequest) (*pb.ClientListResponse, error) {
	query := `
        SELECT id, full_name, address, phone, type
        FROM clients
    `

	// Initialize filters
	filters := []string{}
	args := []interface{}{}
	argCounter := 1

	if in.FullName != "" {
		filters = append(filters, fmt.Sprintf("full_name ILIKE $%d", argCounter))
		args = append(args, "%"+in.FullName+"%")
		argCounter++
	}
	if in.Address != "" {
		filters = append(filters, fmt.Sprintf("address ILIKE $%d", argCounter))
		args = append(args, "%"+in.Address+"%")
		argCounter++
	}
	if in.Phone != "" {
		filters = append(filters, fmt.Sprintf("phone = $%d", argCounter))
		args = append(args, in.Phone)
		argCounter++
	}
	if in.Type != "" {
		filters = append(filters, fmt.Sprintf("type = $%d", argCounter))
		args = append(args, in.Type)
		argCounter++
	}

	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}

	query += " ORDER BY created_at"

	// Pagination
	if in.Limit == 0 {
		in.Limit = 10
	}
	if in.Page == 0 {
		in.Page = 1
	}
	offset := (in.Page - 1) * in.Limit

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCounter, argCounter+1)
	args = append(args, in.Limit, offset)

	var dbClients []entity.DBClient
	err := c.db.Select(&dbClients, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve client list: %w", err)
	}

	// Convert DBClient structs to pb.ClientResponse
	var clients []*pb.ClientResponse
	for _, dbClient := range dbClients {
		clients = append(clients, &pb.ClientResponse{
			Id:       dbClient.Id,
			FullName: dbClient.FullName,
			Address:  dbClient.Address,
			Phone:    dbClient.Phone,
			Type:     dbClient.Type,
		})
	}

	return &pb.ClientListResponse{Clients: clients}, nil
}

func (c *UserRepo) UpdateClient(in *pb.ClientUpdateRequest) (*pb.ClientResponse, error) {
	// Start building the query
	query := "UPDATE clients SET "
	updates := []string{}
	args := []interface{}{}
	argIndex := 1

	// Dynamically add fields to update
	if in.FullName != "" {
		updates = append(updates, fmt.Sprintf("full_name = COALESCE(NULLIF($%d, ''), full_name)", argIndex))
		args = append(args, in.FullName)
		argIndex++
	}
	if in.Address != "" {
		updates = append(updates, fmt.Sprintf("address = COALESCE(NULLIF($%d, ''), address)", argIndex))
		args = append(args, in.Address)
		argIndex++
	}
	if in.Phone != "" {
		updates = append(updates, fmt.Sprintf("phone = COALESCE(NULLIF($%d, ''), phone)", argIndex))
		args = append(args, in.Phone)
		argIndex++
	}
	if in.Type != "" {
		updates = append(updates, fmt.Sprintf("type = COALESCE(NULLIF($%d, ''), type)", argIndex))
		args = append(args, in.Type)
		argIndex++
	}

	// Ensure at least one field is being updated
	if len(updates) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	query += strings.Join(updates, ", ") + " WHERE id = $" + strconv.Itoa(argIndex) + " RETURNING id, full_name, address, phone, type"
	args = append(args, in.Id)

	// Execute the query
	var client pb.ClientResponse
	err := c.db.QueryRowx(query, args...).Scan(
		&client.Id,
		&client.FullName,
		&client.Address,
		&client.Phone,
		&client.Type,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update client: %w", err)
	}

	return &client, nil
}

func (c *UserRepo) DeleteClient(in *pb.UserIDRequest) (*pb.MessageResponse, error) {
	query := `
		DELETE FROM clients
		WHERE id = $1
		RETURNING 'Client deleted'
	`
	var message pb.MessageResponse
	err := c.db.QueryRowx(query, in.Id).Scan(&message.Message)
	if err != nil {
		return nil, fmt.Errorf("failed to delete client: %w", err)
	}
	return &message, nil
}
