package interfaces

import "context"

type MessageService interface {
	CreateMessage(c context.Context, req *CreateMessageReq) (*CreateMessageRes, error)
	GetMessagesByRoomID(c context.Context, roomID string, limit int) ([]*CreateMessageRes, error)
}
