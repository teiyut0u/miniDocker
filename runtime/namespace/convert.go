// Copyright 2025 The Open Container Initiative.
// Copyright 2025 Chunxu Tu.
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

// This file contains modifications by Chunxu Tu.
// original context can be found in runc/libcontainer/specconv/spec_linux.go

package namespace

import (
	"golang.org/x/sys/unix"
)

var (
	MountPropagationMapping = map[string]int{
		"rprivate":    unix.MS_PRIVATE | unix.MS_REC,
		"private":     unix.MS_PRIVATE,
		"rslave":      unix.MS_SLAVE | unix.MS_REC,
		"slave":       unix.MS_SLAVE,
		"rshared":     unix.MS_SHARED | unix.MS_REC,
		"shared":      unix.MS_SHARED,
		"runbindable": unix.MS_UNBINDABLE | unix.MS_REC,
		"unbindable":  unix.MS_UNBINDABLE,
	}

	MountFlags = map[string]struct {
		clear bool
		flag  int
	}{
		// "acl" cannot be mapped to MS_POSIXACL: https://github.com/opencontainers/runc/issues/3738
		"async":         {true, unix.MS_SYNCHRONOUS},
		"atime":         {true, unix.MS_NOATIME},
		"bind":          {false, unix.MS_BIND},
		"defaults":      {false, 0},
		"dev":           {true, unix.MS_NODEV},
		"diratime":      {true, unix.MS_NODIRATIME},
		"dirsync":       {false, unix.MS_DIRSYNC},
		"exec":          {true, unix.MS_NOEXEC},
		"iversion":      {false, unix.MS_I_VERSION},
		"lazytime":      {false, unix.MS_LAZYTIME},
		"loud":          {true, unix.MS_SILENT},
		"mand":          {false, unix.MS_MANDLOCK},
		"noatime":       {false, unix.MS_NOATIME},
		"nodev":         {false, unix.MS_NODEV},
		"nodiratime":    {false, unix.MS_NODIRATIME},
		"noexec":        {false, unix.MS_NOEXEC},
		"noiversion":    {true, unix.MS_I_VERSION},
		"nolazytime":    {true, unix.MS_LAZYTIME},
		"nomand":        {true, unix.MS_MANDLOCK},
		"norelatime":    {true, unix.MS_RELATIME},
		"nostrictatime": {true, unix.MS_STRICTATIME},
		"nosuid":        {false, unix.MS_NOSUID},
		"nosymfollow":   {false, unix.MS_NOSYMFOLLOW}, // since kernel 5.10
		"rbind":         {false, unix.MS_BIND | unix.MS_REC},
		"relatime":      {false, unix.MS_RELATIME},
		"remount":       {false, unix.MS_REMOUNT},
		"ro":            {false, unix.MS_RDONLY},
		"rw":            {true, unix.MS_RDONLY},
		"silent":        {false, unix.MS_SILENT},
		"strictatime":   {false, unix.MS_STRICTATIME},
		"suid":          {true, unix.MS_NOSUID},
		"sync":          {false, unix.MS_SYNCHRONOUS},
		"symfollow":     {true, unix.MS_NOSYMFOLLOW}, // since kernel 5.10
	}

	RecAttrFlags = map[string]struct {
		clear bool
		flag  uint64
	}{
		"rro":            {false, unix.MOUNT_ATTR_RDONLY},
		"rrw":            {true, unix.MOUNT_ATTR_RDONLY},
		"rnosuid":        {false, unix.MOUNT_ATTR_NOSUID},
		"rsuid":          {true, unix.MOUNT_ATTR_NOSUID},
		"rnodev":         {false, unix.MOUNT_ATTR_NODEV},
		"rdev":           {true, unix.MOUNT_ATTR_NODEV},
		"rnoexec":        {false, unix.MOUNT_ATTR_NOEXEC},
		"rexec":          {true, unix.MOUNT_ATTR_NOEXEC},
		"rnodiratime":    {false, unix.MOUNT_ATTR_NODIRATIME},
		"rdiratime":      {true, unix.MOUNT_ATTR_NODIRATIME},
		"rrelatime":      {false, unix.MOUNT_ATTR_RELATIME},
		"rnorelatime":    {true, unix.MOUNT_ATTR_RELATIME},
		"rnoatime":       {false, unix.MOUNT_ATTR_NOATIME},
		"ratime":         {true, unix.MOUNT_ATTR_NOATIME},
		"rstrictatime":   {false, unix.MOUNT_ATTR_STRICTATIME},
		"rnostrictatime": {true, unix.MOUNT_ATTR_STRICTATIME},
		"rnosymfollow":   {false, unix.MOUNT_ATTR_NOSYMFOLLOW}, // since kernel 5.14
		"rsymfollow":     {true, unix.MOUNT_ATTR_NOSYMFOLLOW},  // since kernel 5.14
	}
)
