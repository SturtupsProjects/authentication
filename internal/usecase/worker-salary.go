package usecase

import (
	pb "authentification/internal/generated/user"
	"context"
)

// Methods for salary
func (a *AuthServiceServer) CreateSalary(ctx context.Context, in *pb.SalaryRequest) (*pb.SalaryResponse, error) {

	res, err := a.workers.CreateSalary(in)
	if err != nil {
		a.log.Error("Error in create salary to worker", "error", err)
		return nil, err
	}

	a.log.Info("created salary to worker", "result", res)

	return res, err

}

func (a *AuthServiceServer) UpdateSalary(ctx context.Context, in *pb.SalaryUpdate) (*pb.SalaryResponse, error) {

	res, err := a.workers.UpdateSalary(in)
	if err != nil {
		a.log.Error("Error in update salary to worker", "error", err)
		return nil, err
	}

	a.log.Info("updated salary to worker", "result", res)

	return res, err
}

func (a *AuthServiceServer) GetSalaryByID(ctx context.Context, in *pb.ID) (*pb.SalaryResponse, error) {

	res, err := a.workers.GetSalaryByID(in)
	if err != nil {
		a.log.Error("Error in get salary by id to worker", "error", err)
		return nil, err
	}

	a.log.Info("get salary by id to worker", "result", res)

	return res, err
}

func (a *AuthServiceServer) ListSalaries(ctx context.Context, in *pb.GetSalaryRequest) (*pb.GetSalaryList, error) {

	res, err := a.workers.ListSalaries(in)
	if err != nil {
		a.log.Error("Error in list salaries to worker", "error", err)
		return nil, err
	}

	a.log.Info("list salaries to worker", "result", res)

	return res, err
}

// Methods for Bonuses
func (a *AuthServiceServer) CreateAdjustment(ctx context.Context, in *pb.AdjustmentRequest) (*pb.AdjustmentResponse, error) {

	res, err := a.workers.CreateAdjustment(in)
	if err != nil {
		a.log.Error("Error in create adjustment to worker", "error", err)
		return nil, err
	}

	a.log.Info("created adjustment to worker", "result", res)

	return res, err
}

func (a *AuthServiceServer) UpdateAdjustment(ctx context.Context, in *pb.AdjustmentUpdate) (*pb.AdjustmentResponse, error) {

	res, err := a.workers.UpdateAdjustment(in)
	if err != nil {
		a.log.Error("Error in update adjustment to worker", "error", err)
		return nil, err
	}

	a.log.Info("updated adjustment to worker", "result", res)

	return res, err
}

func (a *AuthServiceServer) CloseAdjustment(ctx context.Context, in *pb.ID) (*pb.AdjustmentResponse, error) {

	res, err := a.workers.CloseAdjustment(in)
	if err != nil {
		a.log.Error("Error in close adjustment to worker", "error", err)
		return nil, err
	}

	a.log.Info("close adjustment to worker", "result", res)

	return res, err
}

func (a *AuthServiceServer) GetAdjustmentByID(ctx context.Context, in *pb.ID) (*pb.AdjustmentResponse, error) {

	res, err := a.workers.GetAdjustmentByID(in)
	if err != nil {
		a.log.Error("Error in get adjustment by id to worker", "error", err)
		return nil, err
	}

	a.log.Info("get adjustment by id to worker", "result", res)

	return res, err
}

func (a *AuthServiceServer) ListAdjustments(ctx context.Context, in *pb.GetAdjustmentRequest) (*pb.AdjustmentList, error) {

	res, err := a.workers.ListAdjustments(in)
	if err != nil {
		a.log.Error("Error in list adjustments to worker", "error", err)
		return nil, err
	}

	a.log.Info("list adjustments to worker", "result", res)

	return res, err
}

func (a *AuthServiceServer) GetWorkerAllInfo(ctx context.Context, in *pb.ID) (*pb.WorkerAllInfo, error) {

	res, err := a.workers.GetWorkerAllInfo(in)
	if err != nil {
		a.log.Error("Error in get worker info to worker", "error", err)
		return nil, err
	}

	a.log.Info("get worker info to worker", "result", res)

	return res, err
}
