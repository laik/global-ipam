package main

import (
	"fmt"
	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	"github.com/containernetworking/cni/pkg/types/current"
	"github.com/containernetworking/cni/pkg/version"
	bv "github.com/containernetworking/plugins/pkg/utils/buildversion"
	"github.com/yametech/global-ipam/pkg/allocator"
	"github.com/yametech/global-ipam/pkg/dns"
	"github.com/yametech/global-ipam/pkg/etcd"
	"net"
	"strings"
)

const PluginName = "global-ipam"

func main() {
	skel.PluginMain(cmdAdd, cmdChek, cmdDel, version.All, bv.BuildString(PluginName))
}

func cmdChek(args *skel.CmdArgs) error {
	return nil
}

func cmdAdd(args *skel.CmdArgs) error {
	_IPAMConfig, confVersion, err := allocator.LoadIPAMConfig(args.StdinData, args.Args)
	if err != nil {
		return err
	}

	result := &current.Result{}

	if _IPAMConfig.ResolvConf != "" {
		dns, err := dns.ParseResolvConf(_IPAMConfig.ResolvConf)
		if err != nil {
			return err
		}
		result.DNS = *dns
	}

	store, err := etcd.New(_IPAMConfig.Name, _IPAMConfig)
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

	for _, ip := range _IPAMConfig.IPArgs {
		requestedIPs[ip.String()] = ip
	}

	for idx, rangeSet := range _IPAMConfig.Ranges {
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

	result.Routes = _IPAMConfig.Routes

	return types.PrintResult(result, confVersion)
}

func cmdDel(args *skel.CmdArgs) error {
	_IPAMConfig, _, err := allocator.LoadIPAMConfig(args.StdinData, args.Args)
	if err != nil {
		return err
	}

	store, err := etcd.New(_IPAMConfig.Name, _IPAMConfig)
	if err != nil {
		return err
	}
	defer store.Close()

	// Loop through all ranges, releasing all IPs, even if an error occurs
	var errors []string
	for idx, rangeSet := range _IPAMConfig.Ranges {
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
