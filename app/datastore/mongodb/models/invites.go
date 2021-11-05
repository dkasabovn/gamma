package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Invite struct {
	ID primitive.ObjectID `bson:"_id" json:"uuid"`
	OrgID string 		  `bson:"org_id" json:"orgId"`
	Permissions int        `bson:"permissions" json:"permissions"`
	CreatedAt time.Time   `bson:"createdAt"`
}

func InviteByID(ctx context.Context, coll *mongo.Collection, id primitive.ObjectID) (*Invite,error) {
	invite := new(Invite)
	filter := bson.M{"_id" : id}

	err := coll.FindOne(ctx, filter, findOneOpts).Decode(&invite)
	if err != nil {
		return nil, err
	}
	return invite, nil
}

func InvitesByOrg(ctx context.Context, coll *mongo.Collection, orgId string) ([]Invite, error){
	var invites []Invite
	filter := bson.M{"org_id" : orgId}
	
	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &invites); err != nil {
		return nil, err
	}

	return invites, nil
}

func NewInvite(ctx context.Context, coll *mongo.Collection, invite Invite) (*mongo.InsertOneResult, error){
	invite.CreatedAt = time.Now()
	result, err := coll.InsertOne(ctx, invite)
	if err != nil {
		return nil, err
	}
	return result, nil
}