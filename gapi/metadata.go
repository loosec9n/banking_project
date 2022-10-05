package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

const (
	grpcGatewayUserAgentKey = "grpcgateway-user-agent"
	grpcUserAgentKey        = "user-agent"
	xForwardForHeader       = "x-forwarded-for"
)

func (serer *Server) extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		//log.Printf("md: %+v\n", md)
		if userAgent := md.Get(grpcGatewayUserAgentKey); len(userAgent) > 0 {
			mtdt.UserAgent = userAgent[0]
		}

		if userAgent := md.Get(grpcUserAgentKey); len(userAgent) > 0 {
			mtdt.UserAgent = userAgent[0]
		}

		if clinetIP := md.Get(xForwardForHeader); len(clinetIP) > 0 {
			mtdt.ClientIP = clinetIP[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIP = p.Addr.String()
	}

	return mtdt
}
