// Package client provides a datasvc client based on a predefined Consul
// service name and relevant tags. Users must only provide the address of a
// Consul server.
package client

import (
	"io"
	"time"

	consulapi "github.com/hashicorp/consul/api"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	"github.com/harmonsoftwaresolutions/hss-inmem-db/datasvc"
)

// New ...
func New(consulAddr string, logger log.Logger) (datasvc.Service, error) {
	apiclient, err := consulapi.NewClient(&consulapi.Config{
		Address: consulAddr,
	})
	if err != nil {
		return nil, err
	}

	// As the implementer of datasvc, we declare and enforce these
	// parameters for all of the datasvc consumers.
	var (
		consulService = "datasvc"
		consulTags    = []string{"prod"}
		passingOnly   = true
		retryMax      = 3
		retryTimeout  = 500 * time.Millisecond
	)

	var (
		sdclient  = consul.NewClient(apiclient)
		instancer = consul.NewInstancer(sdclient, logger, consulService, consulTags, passingOnly)
		endpoints datasvc.Endpoints
	)
	{
		factory := factoryFor(datasvc.MakePostNoteEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.PostNoteEndpoint = retry
	}
	{
		factory := factoryFor(datasvc.MakeGetNoteEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.GetNoteEndpoint = retry
	}
	{
		factory := factoryFor(datasvc.MakePutNoteEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.PutNoteEndpoint = retry
	}
	{
		factory := factoryFor(datasvc.MakePatchNoteEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.PatchNoteEndpoint = retry
	}
	{
		factory := factoryFor(datasvc.MakeDeleteNoteEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.DeleteNoteEndpoint = retry
	}

	return endpoints, nil
}

func factoryFor(makeEndpoint func(datasvc.Service) endpoint.Endpoint) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		service, err := datasvc.MakeClientEndpoints(instance)
		if err != nil {
			return nil, nil, err
		}
		return makeEndpoint(service), nil, nil
	}
}
