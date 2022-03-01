package server

import (
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
