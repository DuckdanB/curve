/*
 *  Copyright (c) 2022 NetEase Inc.
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
 * Project: CurveCli
 * Created Date: 2022-05-11
 * Author: chengyi (Cyber-SiKu)
 */

package cmderror

import (
	"fmt"

	"github.com/opencurve/curve/tools-v2/proto/curvefs/proto/copyset"
	"github.com/opencurve/curve/tools-v2/proto/curvefs/proto/mds"
	"github.com/opencurve/curve/tools-v2/proto/curvefs/proto/topology"
)

// It is considered here that the importance of the error is related to the
// code, and the smaller the code, the more important the error is.
// Need to ensure that the smaller the code, the more important the error is
const (
	CODE_BASE_LINE   = 10000
	CODE_SUCCESS     = 0 * CODE_BASE_LINE
	CODE_RPC_RESULT  = 1 * CODE_BASE_LINE
	CODE_HTTP_RESULT = 2 * CODE_BASE_LINE
	CODE_RPC         = 3 * CODE_BASE_LINE
	CODE_HTTP        = 4 * CODE_BASE_LINE
	CODE_INTERNAL    = 9 * CODE_BASE_LINE
	CODE_UNKNOWN     = 10 * CODE_BASE_LINE
)

type CmdError struct {
	Code    int    `json:"code"`    // exit code
	Message string `json:"message"` // exit message
}

var (
	AllError []*CmdError
)

func init() {
	AllError = make([]*CmdError, 0)
}

func (ce *CmdError) ToError() error {
	return fmt.Errorf(ce.Message)
}

func NewSucessCmdError() *CmdError {
	ret := &CmdError{
		Code:    CODE_SUCCESS,
		Message: "success",
	}
	AllError = append(AllError, ret)
	return ret
}

func NewInternalCmdError(code int, message string) *CmdError {
	if code == 0 {
		return NewSucessCmdError()
	}
	ret := &CmdError{
		Code:    CODE_INTERNAL + code,
		Message: message,
	}

	AllError = append(AllError, ret)
	return ret
}

func NewRpcError(code int, message string) *CmdError {
	if code == 0 {
		return NewSucessCmdError()
	}
	ret := &CmdError{
		Code:    CODE_RPC + code,
		Message: message,
	}
	AllError = append(AllError, ret)
	return ret
}

func NewRpcReultCmdError(code int, message string) *CmdError {
	if code == 0 {
		return NewSucessCmdError()
	}
	ret := &CmdError{
		Code:    CODE_RPC_RESULT + code,
		Message: message,
	}
	AllError = append(AllError, ret)
	return ret
}

func NewHttpError(code int, message string) *CmdError {
	if code == 0 {
		return NewSucessCmdError()
	}
	ret := &CmdError{
		Code:    CODE_HTTP + code,
		Message: message,
	}
	AllError = append(AllError, ret)
	return ret
}

func NewHttpResultCmdError(code int, message string) *CmdError {
	if code == 0 {
		return NewSucessCmdError()
	}
	ret := &CmdError{
		Code:    CODE_HTTP_RESULT + code,
		Message: message,
	}
	AllError = append(AllError, ret)
	return ret
}

func (cmd CmdError) TypeCode() int {
	return cmd.Code / CODE_BASE_LINE * CODE_BASE_LINE
}

func (cmd CmdError) TypeName() string {
	var ret string
	switch cmd.TypeCode() {
	case CODE_SUCCESS:
		ret = "success"
	case CODE_INTERNAL:
		ret = "internal"
	case CODE_RPC:
		ret = "rpc"
	case CODE_RPC_RESULT:
		ret = "rpcResult"
	case CODE_HTTP:
		ret = "http"
	case CODE_HTTP_RESULT:
		ret = "httpResult"
	default:
		ret = "unknown"
	}
	return ret
}

func (e *CmdError) Format(args ...interface{}) {
	e.Message = fmt.Sprintf(e.Message, args...)
}

// The importance of the error is considered to be related to the code,
// please use it under the condition that the smaller the code,
// the more important the error is.
func MostImportantCmdError(err []*CmdError) *CmdError {
	if len(err) == 0 {
		return NewSucessCmdError()
	}
	ret := err[0]
	for _, e := range err {
		if e.Code < ret.Code {
			ret = e
		}
	}
	return ret
}

// keep the most important wrong id, all wrong message will be kept
func MergeCmdError(err []*CmdError) CmdError {
	if len(err) == 0 {
		return *NewSucessCmdError()
	}
	var ret CmdError
	ret.Code = CODE_UNKNOWN
	ret.Message = ""
	for _, e := range err {
		if e.Code == CODE_SUCCESS {
			continue
		} else if e.Code < ret.Code {
			ret.Code = e.Code
		}
		ret.Message = e.Message + "\n" + ret.Message
	}
	ret.Message = ret.Message[:len(ret.Message)-1]
	return ret
}

var (
	ErrSuccess = NewSucessCmdError

	// internal error
	ErrHttpCreateGetRequest = func() *CmdError {
		return NewInternalCmdError(1, "create http get request failed, the error is: %s")
	}
	ErrDataNoExpected = func() *CmdError {
		return NewInternalCmdError(2, "data: %s is not as expected, the error is: %s")
	}
	ErrHttpClient = func() *CmdError {
		return NewInternalCmdError(3, "http client gets error: %s")
	}
	ErrRpcDial = func() *CmdError {
		return NewInternalCmdError(4, "dial to rpc server %s failed, the error is: %s")
	}
	ErrUnmarshalJson = func() *CmdError {
		return NewInternalCmdError(5, "unmarshal json error, the json is %s, the error is %s")
	}
	ErrParseMetric = func() *CmdError {
		return NewInternalCmdError(6, "parse metric %s err!")
	}
	ErrGetMetaserverAddr = func() *CmdError {
		return NewInternalCmdError(7, "get metaserver addr failed, the error is: %s")
	}
	ErrGetClusterFsInfo = func() *CmdError {
		return NewInternalCmdError(8, "get cluster fs info failed, the error is: %s")
	}
	ErrGetAddr = func() *CmdError {
		return NewInternalCmdError(9, "invalid %s addr is: %s")
	}
	ErrMarShalProtoJson = func() *CmdError {
		return NewInternalCmdError(10, "marshal proto to json error, the error is: %s")
	}
	ErrUnknownFsType = func() *CmdError {
		return NewInternalCmdError(11, "unknown fs type: %s")
	}
	ErrAligned = func() *CmdError {
		return NewInternalCmdError(12, "%s should aligned with %s")
	}
	ErrUnknownBitmapLocation = func() *CmdError {
		return NewInternalCmdError(13, "unknown bitmap location: %s")
	}
	ErrParseBytes = func() *CmdError {
		return NewInternalCmdError(14, "invalid %s: %s")
	}
	ErrSplitPeer = func() *CmdError {
		return NewInternalCmdError(15, "split peer %s failed!")
	}
	ErrMarshalJson = func() *CmdError {
		return NewInternalCmdError(16, "marshal %s to json error, the error is: %s")
	}
	ErrCopysetKey = func() *CmdError {
		return NewInternalCmdError(17, "copyset key [%d] not found in %s!")
	}
	ErrQueryCopyset = func() *CmdError {
		return NewInternalCmdError(18, "query copyset failed! the error is: %s")
	}
	ErrOfflineCopysetPeer = func() *CmdError {
		return NewInternalCmdError(19, "peer [%s] is offline!")
	}
	ErrStateCopysetPeer = func() *CmdError {
		return NewInternalCmdError(20, "state in peer[%s]: %s")
	}
	ErrListCopyset = func() *CmdError {
		return NewInternalCmdError(21, "list copyset failed! the error is: %s")
	}
	ErrCheckCopyset = func() *CmdError {
		return NewInternalCmdError(22, "check copyset failed! the error is: %s")
	}
	ErrEtcdOffline = func() *CmdError {
		return NewInternalCmdError(23, "etcd[%s] is offline!")
	}
	ErrMdsOffline = func() *CmdError {
		return NewInternalCmdError(24, "mds[%s] is offline!")
	}
	ErrMetaserverOffline = func() *CmdError {
		return NewInternalCmdError(25, "metaserver[%s] is offline!")
	}

	// http error
	ErrHttpUnreadableResult = func() *CmdError {
		return NewHttpResultCmdError(1, "http response is unreadable, the uri is: %s, the error is: %s")
	}
	ErrHttpResultNoExpected = func() *CmdError {
		return NewHttpResultCmdError(2, "http response is not expected, the hosts is: %s, the suburi is: %s, the result is: %s")
	}
	ErrHttpStatus = func(statusCode int) *CmdError {
		return NewHttpError(statusCode, "the url is: %s, http status code is: %d")
	}

	// rpc error
	ErrRpcCall = func() *CmdError {
		return NewRpcReultCmdError(1, "rpc call is fail, the addr is: %s, the func is %s, the error is: %s")
	}
	ErrUmountFs = func(statusCode int) *CmdError {
		var message string
		code := mds.FSStatusCode(statusCode)
		switch code {
		case mds.FSStatusCode_OK:
			message = "success"
		case mds.FSStatusCode_MOUNT_POINT_NOT_EXIST:
			message = "mountpoint not exist"
		case mds.FSStatusCode_NOT_FOUND:
			message = "fs not found"
		case mds.FSStatusCode_FS_BUSY:
			message = "mountpoint is busy"
		default:
			message = fmt.Sprintf("umount from fs failed!, error is %s", code.String())
		}
		return NewRpcReultCmdError(statusCode, message)
	}
	ErrGetFsInfo = func(statusCode int) *CmdError {
		return NewRpcReultCmdError(statusCode, "get fs info failed: status code is %s")
	}
	ErrGetCopysetOfPartition = func(statusCode int) *CmdError {
		code := topology.TopoStatusCode(statusCode)
		message := fmt.Sprintf("get copyset of partition failed: status code is %s", code.String())
		return NewRpcReultCmdError(statusCode, message)
	}
	ErrDeleteFs = func(statusCode int) *CmdError {
		var message string
		code := mds.FSStatusCode(statusCode)
		switch code {
		case mds.FSStatusCode_OK:
			message = "success"
		case mds.FSStatusCode_NOT_FOUND:
			message = "fs not found!"
		default:
			message = fmt.Sprintf("delete fs failed!, error is %s", code.String())
		}
		return NewRpcReultCmdError(statusCode, message)
	}
	ErrCreateFs = func(statusCode int) *CmdError {
		var message string
		code := mds.FSStatusCode(statusCode)
		switch code {
		case mds.FSStatusCode_OK:
			message = "success"
		case mds.FSStatusCode_FS_EXIST:
			message = "fs exist, but s3 info is not inconsistent"
		case mds.FSStatusCode_S3_INFO_ERROR:
			message = "s3 info is not available"
		default:
			message = fmt.Sprintf("delete fs failed!, error is %s", code.String())
		}
		return NewRpcReultCmdError(statusCode, message)
	}
	ErrGetCopysetsInfo = func(statusCode int) *CmdError {
		code := topology.TopoStatusCode(statusCode)
		message := fmt.Sprintf("get copysets info failed: status code is %s", code.String())
		return NewRpcReultCmdError(statusCode, message)
	}
	ErrCopysetOpStatus = func(statusCode copyset.COPYSET_OP_STATUS, addr string) *CmdError {
		var message string
		code := int(statusCode)
		switch statusCode {
		case copyset.COPYSET_OP_STATUS_COPYSET_OP_STATUS_COPYSET_NOTEXIST:
			message = fmt.Sprintf("not exist in %s", addr)
		case copyset.COPYSET_OP_STATUS_COPYSET_OP_STATUS_SUCCESS:
			message = "ok"
		default:
			message = fmt.Sprintf("op status: %s in %s", statusCode.String(), addr)
		}
		return NewRpcReultCmdError(code, message)
	}
)
