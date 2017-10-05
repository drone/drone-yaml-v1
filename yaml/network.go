package yaml

// Networks represent a container network list.
type Networks struct {
	Networks []*Network
}

// Network represent a container network.
type Network struct {
	Name    string
	Aliases []string
}

// networks:
// 	- frontend
// 	- backend

// networks:
// 	some-network:
// 		aliases:
// 		 - alias1
// 		 - alias3
// 	other-network:
// 		aliases:
// 		 - alias2

// networks:
// 	new:
// 		aliases:
// 			- database
// 	legacy:
// 		aliases:
// 			- mysql

// networks:
// 	app_net:
// 		ipv4_address: 172.16.238.10
// 		ipv6_address: 2001:3984:3989::10

// UnmarshalYAML implements the Unmarshaller interface.
func (n *Networks) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var sliceType []string
	if err := unmarshal(&sliceType); err == nil {
		for _, name := range sliceType {
			n.Networks = append(n.Networks, &Network{
				Name: name,
			})
		}
		return nil
	}

	var mapType map[string]*Network
	if err := unmarshal(&sliceType); err != nil {
		return err
	}
	for name, network := range mapType {
		network.Name = name
		n.Networks = append(n.Networks, network)
	}
	return nil
}
