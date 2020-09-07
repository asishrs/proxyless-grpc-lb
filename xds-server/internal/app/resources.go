package app

import (
	"context"
	"errors"
	"fmt"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	listenerv2 "github.com/envoyproxy/go-control-plane/envoy/api/v2/listener"
	v2route "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"

	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	lv2 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v2"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	cache "github.com/envoyproxy/go-control-plane/pkg/cache/v2"

	ep "github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"

	logger "github.com/asishrs/proxyless-grpc-lb/common/pkg/logger"
	"go.uber.org/zap"

	wrapperspb "github.com/golang/protobuf/ptypes/wrappers"

	hcm "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
)

type podEndPoint struct {
	IP   string
	Port int32
}

func getK8sEndPoints(serviceNames []string) (map[string][]podEndPoint, error) {
	k8sEndPoints := make(map[string][]podEndPoint)

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	endPoints, err := clientset.CoreV1().Endpoints("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Logger.Error("Received error while trying to get EndPoints", zap.Error(err))
	}
	logger.Logger.Debug("Endpoint in the cluster", zap.Int("count", len(endPoints.Items)))
	for _, serviceName := range serviceNames {
		for _, endPoint := range endPoints.Items {
			name := endPoint.GetObjectMeta().GetName()
			if name == serviceName {
				var ips []string
				var ports []int32
				for _, subset := range endPoint.Subsets {
					for _, address := range subset.Addresses {
						ips = append(ips, address.IP)
					}
					for _, port := range subset.Ports {
						ports = append(ports, port.Port)
					}
				}
				logger.Logger.Debug("Endpoint", zap.String("name", name), zap.Any("IP Address", ips), zap.Any("Ports", ports))
				var podEndPoints []podEndPoint
				for _, port := range ports {
					for _, ip := range ips {
						podEndPoints = append(podEndPoints, podEndPoint{ip, port})
					}
				}
				k8sEndPoints[serviceName] = podEndPoints
			}
		}
	}
	return k8sEndPoints, nil
}

func clusterLoadAssignment(podEndPoints []podEndPoint, clusterName string, region string, zone string) []types.Resource {
	var lbs []*ep.LbEndpoint
	for _, podEndPoint := range podEndPoints {
		logger.Logger.Debug("Creating ENDPOINT", zap.String("host", podEndPoint.IP), zap.Int32("port", podEndPoint.Port))
		hst := &core.Address{Address: &core.Address_SocketAddress{
			SocketAddress: &core.SocketAddress{
				Address:  podEndPoint.IP,
				Protocol: core.SocketAddress_TCP,
				PortSpecifier: &core.SocketAddress_PortValue{
					PortValue: uint32(podEndPoint.Port),
				},
			},
		}}

		lbs = append(lbs, &ep.LbEndpoint{
			HostIdentifier: &ep.LbEndpoint_Endpoint{
				Endpoint: &ep.Endpoint{
					Address: hst,
				}},
			HealthStatus: core.HealthStatus_HEALTHY,
		})
	}

	eds := []types.Resource{
		&v2.ClusterLoadAssignment{
			ClusterName: clusterName,
			Endpoints: []*ep.LocalityLbEndpoints{{
				Locality: &core.Locality{
					Region: region,
					Zone:   zone,
				},
				Priority:            0,
				LoadBalancingWeight: &wrapperspb.UInt32Value{Value: uint32(1000)},
				LbEndpoints:         lbs,
			}},
		},
	}
	return eds
}

func createCluster(clusterName string) []types.Resource {
	logger.Logger.Debug("Creating CLUSTER", zap.String("name", clusterName))
	cls := []types.Resource{
		&v2.Cluster{
			Name:                 clusterName,
			LbPolicy:             v2.Cluster_ROUND_ROBIN,
			ClusterDiscoveryType: &v2.Cluster_Type{Type: v2.Cluster_EDS},
			EdsClusterConfig: &v2.Cluster_EdsClusterConfig{
				EdsConfig: &core.ConfigSource{
					ConfigSourceSpecifier: &core.ConfigSource_Ads{},
				},
			},
		},
	}
	return cls
}

func createVirtualHost(virtualHostName, listenerName, clusterName string) *v2route.VirtualHost {
	logger.Logger.Debug("Creating RDS", zap.String("host name", virtualHostName))
	vh := &v2route.VirtualHost{
		Name:    virtualHostName,
		Domains: []string{listenerName},

		Routes: []*v2route.Route{{
			Match: &v2route.RouteMatch{
				PathSpecifier: &v2route.RouteMatch_Prefix{
					Prefix: "",
				},
			},
			Action: &v2route.Route_Route{
				Route: &v2route.RouteAction{
					ClusterSpecifier: &v2route.RouteAction_Cluster{
						Cluster: clusterName,
					},
				},
			},
		}}}
	return vh

}

func createRoute(routeConfigName, virtualHostName, listenerName, clusterName string) []types.Resource {
	vh := createVirtualHost(virtualHostName, listenerName, clusterName)
	rds := []types.Resource{
		&v2.RouteConfiguration{
			Name:         routeConfigName,
			VirtualHosts: []*v2route.VirtualHost{vh},
		},
	}
	return rds
}

func createListener(listenerName string, clusterName string, routeConfigName string) []types.Resource {
	logger.Logger.Debug("Creating LISTENER", zap.String("name", listenerName))
	hcRds := &hcm.HttpConnectionManager_Rds{
		Rds: &hcm.Rds{
			RouteConfigName: routeConfigName,
			ConfigSource: &core.ConfigSource{
				ConfigSourceSpecifier: &core.ConfigSource_Ads{
					Ads: &core.AggregatedConfigSource{},
				},
			},
		},
	}

	manager := &hcm.HttpConnectionManager{
		CodecType:      hcm.HttpConnectionManager_AUTO,
		RouteSpecifier: hcRds,
	}

	pbst, err := ptypes.MarshalAny(manager)
	if err != nil {
		panic(err)
	}

	lds := []types.Resource{
		&v2.Listener{
			Name: listenerName,
			ApiListener: &lv2.ApiListener{
				ApiListener: pbst,
			},
			Address: &core.Address{
				Address: &core.Address_SocketAddress{
					SocketAddress: &core.SocketAddress{
						Protocol: core.SocketAddress_TCP,
						Address:  "0.0.0.0",
						PortSpecifier: &core.SocketAddress_PortValue{
							PortValue: 10000,
						},
					},
				},
			},
			FilterChains: []*listenerv2.FilterChain{{
				Filters: []*listenerv2.Filter{{
					Name: wellknown.HTTPConnectionManager,
					ConfigType: &listenerv2.Filter_TypedConfig{
						TypedConfig: pbst,
					},
				}},
			}},
		}}
	return lds
}

// GenerateSnapshot creates snapshot for each service
func GenerateSnapshot(services []string) (*cache.Snapshot, error) {
	k8sEndPoints, err := getK8sEndPoints(services)
	if err != nil {
		logger.Logger.Error("Error while trying to get EndPoints from k8s cluster", zap.Error(err))
		return nil, errors.New("Error while trying to get EndPoints from k8s cluster")
	}

	logger.Logger.Debug("K8s", zap.Any("EndPoints", k8sEndPoints))

	var eds []types.Resource
	var cds []types.Resource
	var rds []types.Resource
	var lds []types.Resource
	for service, podEndPoints := range k8sEndPoints {
		logger.Logger.Debug("Creating new XDS Entry", zap.String("service", service))
		eds = append(eds, clusterLoadAssignment(podEndPoints, fmt.Sprintf("%s-cluster", service), "my-region", "my-zone")...)
		cds = append(cds, createCluster(fmt.Sprintf("%s-cluster", service))...)
		rds = append(rds, createRoute(fmt.Sprintf("%s-route", service), fmt.Sprintf("%s-vhost", service), fmt.Sprintf("%s-listener", service), fmt.Sprintf("%s-cluster", service))...)
		lds = append(lds, createListener(fmt.Sprintf("%s-listener", service), fmt.Sprintf("%s-cluster", service), fmt.Sprintf("%s-route", service))...)
	}

	version := uuid.New()
	logger.Logger.Debug("Creating Snapshot", zap.String("version", version.String()), zap.Any("EDS", eds), zap.Any("CDS", cds), zap.Any("RDS", rds), zap.Any("LDS", lds))
	snapshot := cache.NewSnapshot(version.String(), eds, cds, rds, lds, []types.Resource{}, []types.Resource{})

	if err := snapshot.Consistent(); err != nil {
		logger.Logger.Error("Snapshot inconsistency", zap.Any("snapshot", snapshot), zap.Error(err))
	}
	return &snapshot, nil
}
