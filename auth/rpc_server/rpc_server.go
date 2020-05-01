package rpc_server

import (
	"context"
	"errors"
	"log"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/nikvas0/dc-homework/auth/objects"
	pb "github.com/nikvas0/dc-homework/lib/proto/auth"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
}

func (s *AuthServer) Validate(ctx context.Context, in *pb.ValidateRequest) (*pb.ValidateReply, error) {
	token := objects.Token{}
	tokenInfo, err := jwt.ParseWithClaims(in.GetToken(), &token, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ACCESS_TOKEN_KEY")), nil
	})
	if err != nil || !tokenInfo.Valid {
		log.Println("Validate request error: Got bad token.")
		return &pb.ValidateReply{}, errors.New("bad token")
	}

	log.Printf("Validate request: success (id=%d).", token.UserID)

	return &pb.ValidateReply{User: token.UserID, Email: token.Email, Role: token.Role}, nil
}
