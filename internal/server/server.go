package server

import (
	"context"
	"errors"
	"io"
	"log"

	"github.com/ac2393921/grpc-chat-go/internal/chat/pb"
)

type server struct {
	rooms []room
}

type room struct {
	id       string
	contents []message
}

type message struct {
	user    string
	content string
}

func (s *server) AddRoom(ctx context.Context, request *pb.RoomRequest) (*pb.RoomInfo, error) {
	s.rooms = append(s.rooms, room{id: request.Id, contents: []message{}})

	index, err := searchRooms(s.rooms, request.Id)
	if err != nil {
		return nil, err
	}
	room := s.rooms[index]
	return &pb.RoomInfo{
		Id:           room.id,
		MessageCount: int32(len(room.contents)),
	}, nil
}

func (s *server) GetRoomInfo(ctx context.Context, request *pb.RoomRequest) (*pb.RoomInfo, error) {
	index, err := searchRooms(s.rooms, request.Id)
	if err != nil {
		return nil, err
	}
	room := s.rooms[index]
	return &pb.RoomInfo{
		Id:           room.id,
		MessageCount: int32(len(room.contents)),
	}, nil
}

func searchRooms(r []room, id string) (int, error) {
	for i, v := range r {
		if v.id == id {
			return i, nil
		}
	}
	return -1, errors.New("Not Found")
}

func (s *server) SendMessage(stream pb.Broadcast_SendMessageServer) error {
	for {
		msg, err := stream.Recv()
		log.Printf("Receive message>> [%s] %s", msg.Name, msg.Content)

		if err == io.EOF {
			return stream.SendAndClose(&pb.SendResult{
				Result: true,
			})
		}
		if err != nil {
			return err
		}

		if msg.Content == "/exit" {
			return stream.SendAndClose(&pb.SendResult{
				Result: true,
			})
		}

		index, err := searchRooms(s.rooms, msg.Id)
		if err != nil {
			return nil
		}

		s.rooms[index].contents = append(
			s.rooms[index].contents,
			message{
				user:    msg.Name,
				content: m.Content,
			},
		)
	}
}

func (s *server) GetMessages(p *pb.MessagesRequest, stream pb.Broadcast_GetMessagesServer) error {
	index, err := searchRooms(s.rooms, p.Id)
	if err != nil {
		return err
	}

	targetRoom := s.rooms[index]

	previousCount := len(targetRoom.contents)
	currentCount := 0

	for {
		targetRoom = s.rooms[index]
		currentCount = len(targetRoom.contents)

		if previousCount < currentCount {
			msg, _ := latestMessage(targetRoom.contents)
			if err := stream.Send(&pb.Message{Id: targetRoom.id, Name: msg.author, Content: msg.content}); err != nil {
				return err
			}
		}
		previousCount = currentCount
	}
}
