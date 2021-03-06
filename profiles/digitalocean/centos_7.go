// Copyright © 2017 The Kubicorn Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package digitalocean

import (
	"fmt"

	"github.com/kubicorn/kubicorn/apis/cluster"
	"github.com/kubicorn/kubicorn/pkg/kubeadm"
)

// NewCentosCluster creates a basic CentOS DigitalOcean cluster.
func NewCentosCluster(name string) *cluster.Cluster {
	return &cluster.Cluster{
		Name:     name,
		Cloud:    cluster.CloudDigitalOcean,
		Location: "sfo2",
		SSH: &cluster.SSH{
			PublicKeyPath: "~/.ssh/id_rsa.pub",
			User:          "root",
		},
		KubernetesAPI: &cluster.KubernetesAPI{
			Port: "443",
		},
		Values: &cluster.Values{
			ItemMap: map[string]string{
				"INJECTEDTOKEN": kubeadm.GetRandomToken(),
			},
		},
		Components: &cluster.Components{
			ComponentVPN: false,
		},
		ServerPools: []*cluster.ServerPool{
			{
				Type:     cluster.ServerPoolTypeMaster,
				Name:     fmt.Sprintf("%s-master", name),
				MaxCount: 1,
				Image:    "centos-7-x64",
				Size:     "s-2vcpu-2gb",
				BootstrapScripts: []string{
					"bootstrap/vpn/openvpnMaster-centos.sh",
					"bootstrap/digitalocean_k8s_centos_7_master.sh",
				},
				Firewalls: []*cluster.Firewall{
					{
						Name: fmt.Sprintf("%s-master", name),
						IngressRules: []*cluster.IngressRule{
							{
								IngressToPort:   "22",
								IngressSource:   "0.0.0.0/0",
								IngressProtocol: "tcp",
							},
							{
								IngressToPort:   "443",
								IngressSource:   "0.0.0.0/0",
								IngressProtocol: "tcp",
							},
							{
								IngressToPort:   "1194",
								IngressSource:   "0.0.0.0/0",
								IngressProtocol: "udp",
							},
							{
								IngressToPort:   "all",
								IngressSource:   fmt.Sprintf("%s-node", name),
								IngressProtocol: "tcp",
							},
						},
						EgressRules: []*cluster.EgressRule{
							{
								EgressToPort:      "all", // By default all egress from VM
								EgressDestination: "0.0.0.0/0",
								EgressProtocol:    "tcp",
							},
							{
								EgressToPort:      "all", // By default all egress from VM
								EgressDestination: "0.0.0.0/0",
								EgressProtocol:    "udp",
							},
						},
					},
				},
			},
			{
				Type:     cluster.ServerPoolTypeNode,
				Name:     fmt.Sprintf("%s-node", name),
				MaxCount: 2,
				Image:    "centos-7-x64",
				Size:     "s-1vcpu-2gb",
				BootstrapScripts: []string{
					"bootstrap/vpn/openvpnNode-centos.sh",
					"bootstrap/digitalocean_k8s_centos_7_node.sh",
				},
				Firewalls: []*cluster.Firewall{
					{
						Name: fmt.Sprintf("%s-node", name),
						IngressRules: []*cluster.IngressRule{
							{
								IngressToPort:   "22",
								IngressSource:   "0.0.0.0/0",
								IngressProtocol: "tcp",
							},
							{
								IngressToPort:   "1194",
								IngressSource:   "0.0.0.0/0",
								IngressProtocol: "udp",
							},
							{
								IngressToPort:   "all",
								IngressSource:   fmt.Sprintf("%s-master", name),
								IngressProtocol: "tcp",
							},
						},
						EgressRules: []*cluster.EgressRule{
							{
								EgressToPort:      "all", // By default all egress from VM
								EgressDestination: "0.0.0.0/0",
								EgressProtocol:    "tcp",
							},
							{
								EgressToPort:      "all", // By default all egress from VM
								EgressDestination: "0.0.0.0/0",
								EgressProtocol:    "udp",
							},
						},
					},
				},
			},
		},
	}
}
