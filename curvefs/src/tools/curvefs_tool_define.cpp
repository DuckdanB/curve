/*
 *  Copyright (c) 2021 NetEase Inc.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

/*
 * Project: curve
 * Created Date: 2021-09-14
 * Author: chengyi01
 */
#include "curvefs/src/tools/curvefs_tool_define.h"

DEFINE_string(mdsAddr, "127.0.0.1:6700", "mds addr");
DEFINE_bool(example, false, "print the example of usage");
DEFINE_string(confPath, "/etc/curvefs/tools.conf", "config file path of tools");
DEFINE_string(fsname, "curvefs", "fs name");
DEFINE_string(mountpoint, "127.0.0.1:/mnt/curvefs-umount-test",
              "curvefs mount in local path");

// topology
DEFINE_string(mds_addr, "127.0.0.1:6700",
              "mds ip and port, separated by \",\"");  // NOLINT
DEFINE_string(cluster_map, "topo_example.json", "cluster topology map.");
DEFINE_uint32(rpcTimeOutMs, 5000u, "rpc time out");
DEFINE_string(op, "", "operation: build_topology");