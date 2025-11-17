// Formally verified E-Voting using Dafny
// Copyright (C) 2025 Authors Gruppe EinS
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package api

import (
	"context"
	"e-voting-service/data/dto"
	"e-voting-service/usecases"
	"log"

	pb "e-voting-service/proto/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func (s *WahlServices_Server) RegisterWahlleiter(_ context.Context, in *pb.RegisterRequest) (*pb.CheckAnmeldungResponse, error) {
	if in.Username == "" || in.Password == "" {
		return nil, status.Errorf(codes.FailedPrecondition, "Username or Password is empty")
	}

	wahlleiter := dto.Wahlleiter{
		Username: in.Username,
		Password: in.Password,
		Email:    in.Email,
	}

	authToken, err := usecases.RegisterWahlleiter(wahlleiter)
	if err != nil {
		log.Printf("error in RegisterWahlleiter_Usecase: %v", err)
		return nil, err
	}

	return &pb.CheckAnmeldungResponse{AuthBearerToken: authToken}, nil
}

func (s *WahlServices_Server) CheckAnmeldung(_ context.Context, in *pb.CheckAnmeldungRequest) (*pb.CheckAnmeldungResponse, error) {
	log.Printf("Received CheckAnmeldung-Request with Username: %s", in.Username)

	if in.Username == "" || in.Password == "" {
		return nil, status.Errorf(codes.FailedPrecondition, "Username or Password is empty")
	}

	authToken, err := usecases.CheckAnmeldung(in.Username, in.Password)
	if err != nil {
		log.Printf("error in CheckAnmeldung_Usecase: %v", err)
		return nil, err
	}

	return &pb.CheckAnmeldungResponse{AuthBearerToken: authToken}, nil
}

func (s *WahlServices_Server) Abmelden(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	md, _ := metadata.FromIncomingContext(ctx)

	isAuth, err := IsAuthenticatedFromHeaderMetadata(md)

	if err != nil {
		log.Printf("Error in Abmelden: %v", err)
		return &emptypb.Empty{}, err
	}

	if isAuth {
		// Token rauslöschen

		token, _ := checkHeadersForAuthToken(md) // error wurde schon gecheckt
		usecases.Abmelden(token)
		return &emptypb.Empty{}, nil
	}

	return &emptypb.Empty{}, status.Errorf(codes.Unauthenticated, "Authtoken nicht angenommen")
}

func (s *WahlServices_Server) CheckToken(ctx context.Context, _ *emptypb.Empty) (*pb.CheckTokenResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)

	isAuth, err := IsAuthenticatedFromHeaderMetadata(md)

	if err != nil {
		log.Printf("Error in CheckToken: %v", err)
		return nil, err
	}

	err = nil
	if !isAuth {
		err = status.Errorf(codes.Unauthenticated, "Authtoken nicht angenommen")
		return nil, err
	}

	return &pb.CheckTokenResponse{IsValid: isAuth}, err
}
