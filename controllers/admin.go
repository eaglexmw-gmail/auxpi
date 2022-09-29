// Copyright (c) 2019 aimerforreimu. All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
//  GNU GENERAL PUBLIC LICENSE
//                        Version 3, 29 June 2007
//
//  Copyright (C) 2007 Free Software Foundation, Inc. <https://fsf.org/>
//  Everyone is permitted to copy and distribute verbatim copies
// of this license document, but changing it is not allowed.
//
// repo: https://github.com/aimerforreimu/auxpi

package controllers

import (
	auxpi "github.com/auxpi/auxpiAll"
	auxpiLog "github.com/auxpi/log"
	"github.com/auxpi/models"
	"github.com/astaxie/beego"
)

type AdminController struct {
	beego.Controller
}

//渲染前端单页面
func (a *AdminController) Index() {
	var site = auxpi.SiteBase{}
	err := site.UnmarshalJSON([]byte(models.GetOption("site_base", "conf")))
	if err != nil {
		auxpiLog.SetAWarningLog("CONTROLLER", err)
	}

	a.TplName = "admin/index.html"

	a.Data["siteUrl"] = site.SiteUrl
}
