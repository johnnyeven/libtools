package testify

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func NewStreamMock(codec grpc.Codec) *StreamMock {
	return &StreamMock{
		codec: codec,
	}
}

type StreamMock struct {
	ctx   context.Context
	codec grpc.Codec
	md    metadata.MD
	bytes []byte
}

func (s *StreamMock) Context() context.Context {
	return s.ctx
}

func (s *StreamMock) SendMsg(m interface{}) error {
	bytes, err := s.codec.Marshal(m)
	if err != nil {
		return err
	}
	s.bytes = bytes
	return nil
}

func (s *StreamMock) RecvMsg(m interface{}) error {
	if s.bytes == nil {
		panic("need set bytes by SendMsg before RecvMsg")
	}
	err := s.codec.Unmarshal(s.bytes, m)
	if err != nil {
		return err
	}
	s.bytes = nil
	return nil
}

func (s *StreamMock) SetHeader(md metadata.MD) error {
	s.md = metadata.Join(s.md, md)

	// just mock to make it can be read both in or out
	s.ctx = metadata.NewIncomingContext(s.ctx, md)
	s.ctx = metadata.NewOutgoingContext(s.ctx, md)
	return nil
}

func (s *StreamMock) SendHeader(metadata.MD) error {
	return nil
}

func (s *StreamMock) SetTrailer(metadata.MD) {
}
