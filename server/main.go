/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
package main

import (
	"context"
	_ "github.com/joho/godotenv/autoload"
	"github.com/softkr/grpcx/db"
	pb "github.com/softkr/grpcx/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"strings"
)

type server struct {
	pb.UnimplementedGreeterServer
}

// WatchStatus 워치 온라인 오프라인
func (s *server) WatchStatus(ctx context.Context, in *pb.WatchState) (*pb.WatchStateReply, error) {
	db.WatchSocketStatus(in.GetSn(), in.GetAddr(), in.GetStatus().String())
	return &pb.WatchStateReply{Status: 200, Message: "성공했습니다."}, nil
}

// 워치 상태 업데이트
func (s *server) WatchUpdate(ctx context.Context, in *pb.WatchUpdates) (*pb.WatchStateReply, error) {
	db.WatchUpdate(in.GetSn(), in.GetWear(), in.GetTakeMedicine())
	return &pb.WatchStateReply{Status: 200, Message: "성공했습니다."}, nil
}

// GetProject 프로젝트 구하기
func (s *server) GetProject(ctx context.Context, in *pb.ProjectRequest) (*pb.ProjectReply, error) {
	result := db.FindWatch(in.Guid)
	return &pb.ProjectReply{Status: 200, Project: result}, nil
}

// SetFileInfo 파일 업로드전 파일 상세정보 저장
func (s *server) SetFileInfo(ctx context.Context, in *pb.SetFileInfoRequest) (*pb.SetFileInfoReply, error) {
	// log.Printf("Received: %v", in.GetSubFile())
	subFile := strings.Split(in.SubFile, "|")
	data := db.FileInfo{in.GetGuid(), in.FileName, in.VideoMd5, subFile, 0}
	db.Insert(&data)
	return &pb.SetFileInfoReply{Message: in.GetSubFile()}, nil
}

// PutFileInfo 서브파일 업로드 카운트 추가
func (s *server) PutFileInfo(ctx context.Context, in *pb.PutFileInfoRequest) (*pb.PutFileInfoReply, error) {
	result := db.Update(in.GetSubFile())
	// fmt.Println(in.GetSubfile())
	return &pb.PutFileInfoReply{SubFileCount: result}, nil
}

// FindSubFile 서브파일 존재여부 확인
func (s *server) FindSubFile(ctx context.Context, in *pb.GetFindFileInfoRequest) (*pb.GetFindFileInReply, error) {
	// fmt.Println(in.GetSubfile())
	result := db.Find(in.GetSubFile())
	subFile := strings.Join(result.SubFile, " ")

	return &pb.GetFindFileInReply{Guid: result.Guid, FileName: result.FileName, VideoMd5: result.VideoMD5, SubFile: subFile}, nil
}

// DeleteFileInfo 업로드 파일 끝나면 파일정도 삭제
func (s *server) DeleteFileInfo(ctx context.Context, in *pb.RemoveFileInfoRequest) (*pb.RemoveFileInfoReply, error) {
	db.DeleteOne(in.GetVideoMd5())
	return &pb.RemoveFileInfoReply{
		Message: "ok",
	}, nil
}

// SubFileCount 서브파일 숫자 카운팅하는 함수임
func (s *server) SubFileCount(ctx context.Context, in *pb.SubFileCountRequest) (*pb.SubFileCountReply, error) {
	count := db.Count(in.VideoMd5)
	return &pb.SubFileCountReply{
		Message: count,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", os.Getenv("PORT"))
	if err != nil {
		log.Printf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	// log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Printf("failed to serve: %v", err)
	}
}
