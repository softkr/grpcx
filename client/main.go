package client

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
	pb "github.com/softkr/grpcx/proto"
	"google.golang.org/grpc"
)

var address string

type FileInfo struct {
	Guid     string
	FileName string
	VideoMD5 string
	SubFile  []string
}

func init() {
	address = fmt.Sprintf("%v%v", os.Getenv("HOST"), os.Getenv("PORT"))
}

// WatchOnOff 워치 온라인 오프라인
func WatchState(watchNo string, prot string, status pb.WatchState_STATUS) (data *pb.WatchStateReply) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	result, err := c.WatchStatus(ctx, &pb.WatchState{Sn: watchNo, Addr: prot, Status: status})
	if err != nil {
		log.Printf("cloud not greet: %v", err)
	}
	return result
}

// 워치 업데이트
func WatchUpdate(watchNo string, step, camerashoot int32) *pb.WatchStateReply {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	result, err := c.WatchUpdate(ctx, &pb.WatchUpdates{Sn: watchNo, Wear: step, TakeMedicine: camerashoot})
	if err != nil {
		log.Printf("cloud not greet: %v", err)
	}
	return result
}

func GetProject(guid string) (data *pb.ProjectReply) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	result, err := c.GetProject(ctx, &pb.ProjectRequest{Guid: guid})
	if err != nil {
		log.Printf("cloud not greet: %v", err)
	}
	// log.Println(result.Project)
	return result
}

func SetFileInfo(guid string, fileName string, videoMd5 string, subFile []string) (data *pb.SetFileInfoReply) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	subFile1 := strings.Join(subFile, "|")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	result, err := c.SetFileInfo(ctx, &pb.SetFileInfoRequest{Guid: guid, FileName: fileName, VideoMd5: videoMd5, SubFile: subFile1})
	if err != nil {
		log.Printf("could not greet: %v", err)
	}
	return result
}

func GetFileInfo(subFile string) (data *pb.GetFindFileInReply) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	result, err := c.FindSubFile(ctx, &pb.GetFindFileInfoRequest{SubFile: subFile})
	if err != nil {
		log.Printf("could not greet: %v", err)
	}
	return result
}

func PutFileInfo(subFile string) (data *pb.PutFileInfoReply) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	result, err := c.PutFileInfo(ctx, &pb.PutFileInfoRequest{SubFile: subFile})
	if err != nil {
		log.Printf("could not greet: %v", err)
	}
	// log.Printf("Greeting: %s", r)
	return result
}

func DeleteFileInfo(videoMd5 string) (data *pb.RemoveFileInfoReply) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	result, err := c.DeleteFileInfo(ctx, &pb.RemoveFileInfoRequest{VideoMd5: videoMd5})
	if err != nil {
		log.Printf("could not greet: %v", err)
	}
	// log.Printf("Greeting: %s", r)
	return result
}

// 카운터 추가
func SubFileCount(videoMd5 string) int32 {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	result, err := c.SubFileCount(ctx, &pb.SubFileCountRequest{VideoMd5: videoMd5})
	if err != nil {
		log.Printf("could not greet: %v", err)
	}
	return result.Message
}
