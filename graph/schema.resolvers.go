package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"blumer-ms-refers/graph/model"
	"blumer-ms-refers/middleware"
	"context"
)

// Refer is the resolver for the refer field.
func (r *mutationResolver) Refer(ctx context.Context, username string) (bool, error) {
	ctxUser, errExtension := middleware.GetCurrentUserFromCTX(ctx)
	if errExtension != nil {
		return false, PresentTypedError(ctx, *errExtension)
	}

	errExtension = r.Repository.Refer(ctxUser.UserID, username)
	if errExtension != nil {
		return false, PresentTypedError(ctx, *errExtension)
	}

	return true, nil
}

// ReferrerInfo is the resolver for the referrerInfo field.
func (r *queryResolver) ReferrerInfo(ctx context.Context) (*model.ReferrerInfo, error) {
	ctxUser, errExtension := middleware.GetCurrentUserFromCTX(ctx)
	if errExtension != nil {
		return nil, PresentTypedError(ctx, *errExtension)
	}

	profile, err := r.Repository.Find(ctxUser.UserID)
	if err != nil {
		return nil, PresentTypedError(ctx, ErrorExtensionParams{AppError: DataSourceError, Reason: err.Error()})
	}

	return &model.ReferrerInfo{
		UserID: profile.UserID,
		Reward: profile.Reward,
	}, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
