// Copyright 2015 CNI authors
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

package store

import "net"

const GLOBAL_IPAM = "global-ipam"

const UNIX_SOCK_PATH = "/var/run/global-ipam.sock"

type Store interface {
	Lock() error
	Unlock() error
	Close() error
	GetByID(id, ip string) (net.IP, error)
	Reserve(id string, ip net.IP, rangeID string) (bool, error)
	LastReservedIP(rangeID string) (net.IP, error)
	Release(ip net.IP) error
	ReleaseByID(id string) error
}

type LastReservedIPResponse struct {
	IP    string `json:"ip"`
	Error error  `json:"error"`
}
type ReserveResponse struct {
	Reserved bool  `json:"reserved"`
	Error    error `json:"error"`
}
type ReleaseResponse struct {
	IsRelease bool  `json:"isRelease"`
	Error     error `json:"error"`
}
