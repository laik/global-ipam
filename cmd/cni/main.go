package main

import (
	"fmt"

	"net"
	"strings"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	typesVer "github.com/containernetworking/cni/pkg/types/040"
	"github.com/containernetworking/cni/pkg/version"
	"github.com/containernetworking/plugins/pkg/utils/buildversion"
	"github.com/yametech/global-ipam/pkg/allocator"
	"github.com/yametech/global-ipam/pkg/client"
	"github.com/yametech/global-ipam/pkg/dns"
)

const PluginName = "global-ipam"

func main() {
	skel.PluginMain(cmdAdd, cmdChek, cmdDel, version.All, buildversion.BuildString(PluginName))
}

func cmdChek(args *skel.CmdArgs) error { return nil }

func cmdAdd(args *skel.CmdArgs) error {
	ipamConfig, confVersion, err := allocator.LoadIPAMConfig(args.StdinData, args.Args)
	if err != nil {
		return err
	}
	result := &typesVer.Result{}

	if ipamConfig.ResolvConf != "" {
		dns, err := dns.ParseResolvConf(ipamConfig.ResolvConf)
		if err != nil {
			return err
		}
		result.DNS = *dns
	}

	store, err := client.New(ipamConfig.Name, ipamConfig)
	if err != nil {
		return err
	}
	
	defer store.Close()

	// Keep the allocators we used, so we can release all IPs if an error
	// occurs after we start allocating
	var allocators []*allocator.IPAllocator

	// store all requested IPs in a map, so we can easily remove ones we use
	// and error if some remain
	requestedIPs := map[string]net.IP{} //net.IP cannot be a key
	for _, ip := range ipamConfig.IPArgs {
		requestedIPs[ip.String()] = ip
	}

	for idx, rangeSet := range ipamConfig.Ranges {
		allocator := allocator.NewIPAllocator(&rangeSet, store, idx)

		// Check to see if there are any custom IPs requested in this range.
		var requestedIP net.IP
		for k, ip := range requestedIPs {
			if rangeSet.Contains(ip) {
				requestedIP = ip
				delete(requestedIPs, k)
				break
			}
		}

		ipConf, err := allocator.Get(args.ContainerID, requestedIP)
		if err != nil {
			// Deallocate all already allocated IPs
			for _, alloc := range allocators {
				_ = alloc.Release(args.ContainerID)
			}
			return fmt.Errorf("failed to allocate for range %d: %v", idx, err)
		}

		allocators = append(allocators, allocator)

		result.IPs = append(result.IPs, ipConf)
	}

	// If an IP was requested that wasn't fulfilled, fail
	if len(requestedIPs) != 0 {
		for _, alloc := range allocators {
			_ = alloc.Release(args.ContainerID)
		}
		errStr := "failed to allocate all requested IPs:"
		for _, ip := range requestedIPs {
			errStr = errStr + " " + ip.String()
		}
		return fmt.Errorf(errStr)
	}

	result.Routes = ipamConfig.Routes

	return types.PrintResult(result, confVersion)
}

func cmdDel(args *skel.CmdArgs) error {
	ipamConfig, _, err := allocator.LoadIPAMConfig(args.StdinData, args.Args)
	if err != nil {
		return err
	}

	store, err := client.New(ipamConfig.Name, ipamConfig)
	if err != nil {
		return err
	}
	defer store.Close()

	// Loop through all ranges, releasing all IPs, even if an error occurs
	var errors []string
	for idx, rangeSet := range ipamConfig.Ranges {
		ipAllocator := allocator.NewIPAllocator(&rangeSet, store, idx)
		err := ipAllocator.Release(args.ContainerID)
		if err != nil {
			errors = append(errors, err.Error())
		}
	}

	if errors != nil {
		return fmt.Errorf(strings.Join(errors, ";"))
	}

	return nil
}
