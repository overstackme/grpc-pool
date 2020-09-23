package ns

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)

func Dial(callerService string) (conn *grpc.ClientConn, err error) {
	s := fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err = grpc.DialContext(
		ctx,
		fmt.Sprintf("ns:///%s", callerService),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(s),
	)
	return
}
