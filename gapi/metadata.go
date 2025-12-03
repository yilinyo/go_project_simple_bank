package gapi

import (
	"context"
	"log"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGateWayUserAgentHeader = "grpcgateway-user-agent"
	grpcGateWayClientIPHeader  = "x-forwarded-for"
	grpcUserAgentHeader        = "user-agent"
	grpcClientIpHeader         = "grpc-client"
)

type Metadata struct {
	UserAgent string
	ClientIp  string
}

func (server *Server) extractMetadata(ctx context.Context) *Metadata {

	var userAgent, clientIp string

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Printf("md: %+v\n", md)
		userAgents := md.Get(grpcGateWayUserAgentHeader)
		clientIps := md.Get(grpcGateWayClientIPHeader)
		userAgents = append(userAgents, md.Get(grpcUserAgentHeader)...)
		if p, ok := peer.FromContext(ctx); ok {
			clientIps = append(clientIps, p.Addr.String())
		}
		if len(userAgents) > 0 {
			userAgent = userAgents[0]
		} else {
			userAgent = "unknown"
		}
		if len(clientIps) > 0 {
			clientIp = clientIps[0]
		} else {
			clientIp = "unknown"
		}
	}
	return &Metadata{
		UserAgent: userAgent,
		ClientIp:  clientIp,
	}
}
