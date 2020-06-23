// Copyright 2019 HAProxy Technologies
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
//

package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	client_native "github.com/haproxytech/client-native/v2"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/maps"
)

//MapsCreateRuntimeMapHandlerImpl implementation of the MapsCreateRuntimeMapHandler interface using client-native client
type MapsCreateRuntimeMapHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h *MapsCreateRuntimeMapHandlerImpl) Handle(params maps.CreateRuntimeMapParams, principal interface{}) middleware.Responder {
	file, header, err := params.HTTPRequest.FormFile("fileUpload")
	if err != nil {
		return maps.NewCreateRuntimeMapBadRequest()
	}
	defer file.Close()

	me, err := h.Client.Runtime.CreateMap(file, *header)
	if err != nil {
		status := misc.GetHTTPStatusFromErr(err)
		return maps.NewCreateRuntimeMapDefault(status).WithPayload(misc.SetError(status, err.Error()))
	}
	// Just a quick test, if this resolves the problem, we'll need to update the payload to include a reload id
	h.ReloadAgent.Reload()
	return maps.NewCreateRuntimeMapCreated().WithPayload(me)
}

//GetMapsHandlerImpl implementation of the GetAllRuntimeMapFilesHandler interface using client-native client
type GetMapsHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//Handle executing the request and returning a response
func (h *GetMapsHandlerImpl) Handle(params maps.GetAllRuntimeMapFilesParams, principal interface{}) middleware.Responder {
	mapFiles, err := h.Client.Runtime.ShowMaps()
	if err != nil {
		status := misc.GetHTTPStatusFromErr(err)
		return maps.NewShowRuntimeMapDefault(status).WithPayload(misc.SetError(status, err.Error()))
	}
	return maps.NewGetAllRuntimeMapFilesOK().WithPayload(mapFiles)
}

//GetMapHandlerImpl implementation of the MapsGetOneRuntimeMapHandler interface using client-native client
type GetMapHandlerImpl struct {
	Client *client_native.HAProxyClient
}

func (h *GetMapHandlerImpl) Handle(params maps.GetOneRuntimeMapParams, principal interface{}) middleware.Responder {
	m, err := h.Client.Runtime.GetMap(params.Name)
	if err != nil {
		status := misc.GetHTTPStatusFromErr(err)
		return maps.NewGetOneRuntimeMapDefault(status).WithPayload(misc.SetError(status, err.Error()))
	}
	if m == nil {
		return maps.NewGetOneRuntimeMapNotFound()
	}
	return maps.NewGetOneRuntimeMapOK().WithPayload(m)
}

//ClearMapHandlerImpl implementation of the ClearRuntimeMapHandler interface using client-native client
type ClearMapHandlerImpl struct {
	Client *client_native.HAProxyClient
}

func (h *ClearMapHandlerImpl) Handle(params maps.ClearRuntimeMapParams, principal interface{}) middleware.Responder {
	forceDelete := false
	if params.ForceDelete != nil {
		forceDelete = *params.ForceDelete
	}
	err := h.Client.Runtime.ClearMap(params.Name, forceDelete)
	if err != nil {
		status := misc.GetHTTPStatusFromErr(err)
		return maps.NewClearRuntimeMapDefault(status).WithPayload(misc.SetError(status, err.Error()))
	}
	return maps.NewClearRuntimeMapNoContent()
}

//ShowMapHandlerImpl implementation of the ShowMapHandlerImpl interface using client-native client
type ShowMapHandlerImpl struct {
	Client *client_native.HAProxyClient
}

func (h *ShowMapHandlerImpl) Handle(params maps.ShowRuntimeMapParams, principal interface{}) middleware.Responder {
	m, err := h.Client.Runtime.ShowMapEntries(params.Map)
	if err != nil {
		status := misc.GetHTTPStatusFromErr(err)
		return maps.NewShowRuntimeMapDefault(status).WithPayload(misc.SetError(status, err.Error()))
	}
	if m == nil {
		return maps.NewShowRuntimeMapNotFound()
	}
	return maps.NewShowRuntimeMapOK().WithPayload(m)
}

//AddMapEntryHandlerImpl implementation of the AddMapEntryHandler interface using client-native client
type AddMapEntryHandlerImpl struct {
	Client *client_native.HAProxyClient
}

func (h *AddMapEntryHandlerImpl) Handle(params maps.AddMapEntryParams, principal interface{}) middleware.Responder {
	err := h.Client.Runtime.AddMapEntry(params.Map, params.Data.Key, params.Data.Value)
	if err != nil {
		status := misc.GetHTTPStatusFromErr(err)
		return maps.NewAddMapEntryDefault(status).WithPayload(misc.SetError(status, err.Error()))
	}
	return maps.NewAddMapEntryCreated().WithPayload(params.Data)
}

//GetRuntimeMapEntryHandlerImpl implementation of the GetRuntimeMapEntryHandler interface using client-native client
type GetRuntimeMapEntryHandlerImpl struct {
	Client *client_native.HAProxyClient
}

func (h *GetRuntimeMapEntryHandlerImpl) Handle(params maps.GetRuntimeMapEntryParams, principal interface{}) middleware.Responder {
	m, err := h.Client.Runtime.GetMapEntry(params.Map, params.ID)
	if err != nil {
		status := misc.GetHTTPStatusFromErr(err)
		return maps.NewGetRuntimeMapEntryDefault(status).WithPayload(misc.SetError(status, err.Error()))
	}
	if m == nil {
		return maps.NewGetRuntimeMapEntryNotFound()
	}
	return maps.NewGetRuntimeMapEntryOK().WithPayload(m)
}

//ReplaceRuntimeMapEntryHandlerImpl implementation of the ReplaceRuntimeMapEntryHandler interface using client-native client
type ReplaceRuntimeMapEntryHandlerImpl struct {
	Client *client_native.HAProxyClient
}

func (h *ReplaceRuntimeMapEntryHandlerImpl) Handle(params maps.ReplaceRuntimeMapEntryParams, principal interface{}) middleware.Responder {
	err := h.Client.Runtime.SetMapEntry(params.Map, params.ID, *params.Data.Value)
	if err != nil {
		status := misc.GetHTTPStatusFromErr(err)
		return maps.NewGetRuntimeMapEntryDefault(status).WithPayload(misc.SetError(status, err.Error()))
	}

	e, err := h.Client.Runtime.GetMapEntry(params.Map, params.ID)
	if err != nil {
		return maps.NewReplaceRuntimeMapEntryNotFound()
	}
	return maps.NewGetRuntimeMapEntryOK().WithPayload(e)
}

//DeleteRuntimeMapEntryHandlerImpl implementation of the DeleteRuntimeMapEntryHandler interface using client-native client
type DeleteRuntimeMapEntryHandlerImpl struct {
	Client *client_native.HAProxyClient
}

func (h *DeleteRuntimeMapEntryHandlerImpl) Handle(params maps.DeleteRuntimeMapEntryParams, principal interface{}) middleware.Responder {
	err := h.Client.Runtime.DeleteMapEntry(params.Map, params.ID)
	if err != nil {
		status := misc.GetHTTPStatusFromErr(err)
		return maps.NewDeleteRuntimeMapEntryDefault(status).WithPayload(misc.SetError(status, err.Error()))
	}
	return maps.NewDeleteRuntimeMapEntryNoContent()
}
