package clientcommon

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
)

type GrpcBuilder struct {
	target  string
	options []grpc.DialOption
	timeout time.Duration
	logger  slog.Logger
}

func NewGrpcBuilder(target string) *GrpcBuilder {
	return &GrpcBuilder{
		target:  target,
		timeout: 5 * time.Second,
	}
}

func (builder *GrpcBuilder) SetTimeout(timeout time.Duration) {
	builder.logger.Info("Setting timeout: ", "timeout", timeout)
	builder.timeout = timeout
}

func (builder *GrpcBuilder) SetGrpcDialOptions(dialOptions ...grpc.DialOption) {
	builder.logger.Info("Setting grpc dial option")
	builder.options = append(builder.options, dialOptions...)
}

func (builder *GrpcBuilder) Build() (*grpc.ClientConn, error) {
	builder.logger.Info("Building client connection")
	ctx, cancel := context.WithTimeout(context.Background(), builder.timeout)
	defer cancel()
	conn, err := grpc.DialContext(ctx, builder.target, builder.options...)

	if err != nil {
		builder.logger.Error("Error in creating client connection")
		return nil, err
	}
	builder.logger.Info("Successfully created client connetion")
	return conn, nil
}
