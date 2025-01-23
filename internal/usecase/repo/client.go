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
		INSERT INTO clients (full_name, address, phone, type, client_type, company_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, full_name, address, phone, type, company_id
	`
	var client pb.ClientResponse
	err := c.db.QueryRowx(query, in.FullName, in.Address, in.Phone, in.Type, in.ClientType, in.CompanyId).Scan(
		&client.Id,
		&client.FullName,
		&client.Address,
		&client.Phone,
		&client.Type,
		&client.CompanyId,
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
		WHERE id = $1 AND company_id = $2
	`
	var client pb.ClientResponse
	err := c.db.QueryRowx(query, in.Id, in.CompanyId).Scan(
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
	// Базовый запрос
	baseQuery := `
        FROM clients
        WHERE company_id = $1 AND client_type = $2
    `
	args := []interface{}{in.CompanyId, in.ClientType}
	argCounter := 3

	// Добавляем фильтры
	if in.FullName != "" {
		baseQuery += fmt.Sprintf(" AND full_name ILIKE $%d", argCounter)
		args = append(args, "%"+in.FullName+"%")
		argCounter++
	}
	if in.Address != "" {
		baseQuery += fmt.Sprintf(" AND address ILIKE $%d", argCounter)
		args = append(args, "%"+in.Address+"%")
		argCounter++
	}
	if in.Phone != "" {
		baseQuery += fmt.Sprintf(" AND phone = $%d", argCounter)
		args = append(args, in.Phone)
		argCounter++
	}
	if in.Type != "" {
		baseQuery += fmt.Sprintf(" AND type = $%d", argCounter)
		args = append(args, in.Type)
		argCounter++
	}

	// Запрос для total_count (без LIMIT и OFFSET)
	countQuery := "SELECT COUNT(*) " + baseQuery

	// Запрос для данных с пагинацией
	dataQuery := "SELECT id, full_name, address, phone, type, company_id " + baseQuery
	if in.Limit > 0 && in.Page > 0 {
		offset := (in.Page - 1) * in.Limit
		dataQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCounter, argCounter+1)
		args = append(args, in.Limit, offset)
	}

	// Выполнение запросов
	var totalCount int
	var dbClients []entity.DBClient

	tx, err := c.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Проверяем длину args перед использованием
	if len(args) < argCounter-1 {
		return nil, fmt.Errorf("insufficient arguments for query")
	}

	// Выполняем запрос для total_count
	if err := tx.Get(&totalCount, countQuery, args[:len(args)-2]...); err != nil { // Use appropriate slice
		return nil, fmt.Errorf("failed to retrieve total count: %w", err)
	}

	// Выполняем запрос для данных
	if err := tx.Select(&dbClients, dataQuery, args...); err != nil {
		return nil, fmt.Errorf("failed to retrieve client list: %w", err)
	}

	// Преобразование результата
	clients := make([]*pb.ClientResponse, len(dbClients))
	for i, dbClient := range dbClients {
		clients[i] = &pb.ClientResponse{
			Id:        dbClient.Id,
			FullName:  dbClient.FullName,
			Address:   dbClient.Address,
			Phone:     dbClient.Phone,
			Type:      dbClient.Type,
			CompanyId: dbClient.CompanyId,
		}
	}

	// Возвращаем результат
	return &pb.ClientListResponse{
		Clients:    clients,
		TotalCount: int64(totalCount),
	}, nil
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

	query += strings.Join(updates, ", ") + " WHERE id = $" + strconv.Itoa(argIndex) + " AND company_id = $" + strconv.Itoa(argIndex+1) +
		" RETURNING id, full_name, address, phone, type"
	args = append(args, in.Id, in.CompanyId)

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
		WHERE id = $1 AND company_id = $2
		RETURNING id, full_name
	`
	var id string
	var fullName string
	err := c.db.QueryRowx(query, in.Id, in.CompanyId).Scan(&id, &fullName)
	if err != nil {
		return nil, fmt.Errorf("failed to delete client: %w", err)
	}
	return &pb.MessageResponse{Message: fmt.Sprintf("Client %s (ID: %s) deleted", fullName, id)}, nil
}
