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

package middleware

import (
	"net/http"
	"strconv"

	"github.com/auxpi/auxpiAll"
	"github.com/auxpi/auxpiAll/e"
	"github.com/auxpi/bootstrap"
	"github.com/auxpi/log"
	"github.com/auxpi/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/validation"
)

var site = auxpi.SiteBase{}

func init() {
	err := site.UnmarshalJSON([]byte(models.GetOption("site_base", "conf")))
	if err != nil {
		auxpiLog.SetAWarningLog("CONTROLLER", err)
	}
}

//未登录用户重定向
var CookieAuthCheck = func(ctx *context.Context) {
	sid := ctx.GetCookie("AuXPI_id")
	id, _ := strconv.Atoi(sid)
	at, _ := ctx.GetSecureCookie(bootstrap.SiteConfig.AuxpiSalt, "AuXPI_at")
	un, _ := ctx.GetSecureCookie(bootstrap.SiteConfig.AuxpiSalt, "AuXPI_uname")
	em, _ := ctx.GetSecureCookie(bootstrap.SiteConfig.AuxpiSalt, "AuXPI_email")
	var userCookie = auxpi.AuxpiCookie{
		UName:      un,
		Email:      em,
		ID:         id,
		Version:    ctx.GetCookie("AuXPI_v"),
		AuxpiToken: at,
	}
	valid := validation.Validation{}
	b, err := valid.Valid(&userCookie)

	if err != nil {
		auxpiLog.SetAWarningLog("COOKIE", err)
		beego.Alert("没有 Cookie")
	}
	if !b {
		beego.Alert(userCookie)
		ctx.Redirect(http.StatusFound, site.SiteUrl + "/login")
		beego.Alert("没有 Cookie")
		return
	}

	var sName = `_email_` + userCookie.Email +
		`_user_` + userCookie.UName +
		`_id_` + sid +
		`_version_` + userCookie.Version

	token := ctx.Input.Session(sName)

	if token != userCookie.AuxpiToken {
		ctx.Redirect(http.StatusFound, site.SiteUrl + "/login")
		return
	}

}

//已登录用户重定向
var CookieAuthedCheck = func(ctx *context.Context) {
	sid := ctx.GetCookie("AuXPI_id")
	id, _ := strconv.Atoi(sid)
	at, _ := ctx.GetSecureCookie(bootstrap.SiteConfig.AuxpiSalt, "AuXPI_at")
	un, _ := ctx.GetSecureCookie(bootstrap.SiteConfig.AuxpiSalt, "AuXPI_uname")
	em, _ := ctx.GetSecureCookie(bootstrap.SiteConfig.AuxpiSalt, "AuXPI_email")

	var userCookie = auxpi.AuxpiCookie{
		UName:      un,
		Email:      em,
		ID:         id,
		Version:    ctx.GetCookie("AuXPI_v"),
		AuxpiToken: at,
	}
	valid := validation.Validation{}
	b, err := valid.Valid(&userCookie)

	if err != nil {
		auxpiLog.SetAWarningLog("COOKIE", err)

	}
	if b {

		var sName = `_email_` + userCookie.Email +
			`_user_` + userCookie.UName +
			`_id_` + sid +
			`_version_` + userCookie.Version

		token := ctx.Input.Session(sName)

		if token != userCookie.AuxpiToken {
			ctx.Redirect(http.StatusFound, site.SiteUrl + "/login")
			return
		}

		//定位到用户首页
		ctx.Redirect(http.StatusFound, site.SiteUrl + "/users/index")

		return
	}

}

//验证 cookie 合法性
var CookieSignCheck = func(ctx *context.Context) {
	sid := ctx.GetCookie("AuXPI_id")
	id, _ := strconv.Atoi(sid)
	at, _ := ctx.GetSecureCookie(bootstrap.SiteConfig.AuxpiSalt, "AuXPI_at")
	un, _ := ctx.GetSecureCookie(bootstrap.SiteConfig.AuxpiSalt, "AuXPI_uname")
	em, _ := ctx.GetSecureCookie(bootstrap.SiteConfig.AuxpiSalt, "AuXPI_email")

	//如果全部是空，才能判定为是游客，否者直接销毁所有的 cookie 才能上传
	if sid == "" &&
		at == "" &&
		un == "" &&
		em == "" && ctx.GetCookie("AuXPI_v") == "" {
		return
	}

	//开始验证是否为合法用户
	var userCookie = auxpi.AuxpiCookie{
		UName:      un,
		Email:      em,
		ID:         id,
		Version:    ctx.GetCookie("AuXPI_v"),
		AuxpiToken: at,
	}
	valid := validation.Validation{}
	b, err := valid.Valid(&userCookie)

	if err != nil {
		auxpiLog.SetAWarningLog("UPLOAD_COOKIE", err)
		destoryCookie(ctx)

	}
	if !b {
		//不合法 cookie 直接销毁,然后重定向主页
		destoryCookie(ctx)

		if ajaxErrorResp(ctx) {
			return
		}

		//ctx.Redirect(http.StatusFound, site.SiteUrl)
		return
	}

	//开始验证 cookie 合法性
	var sName = `_email_` + userCookie.Email +
		`_user_` + userCookie.UName +
		`_id_` + sid +
		`_version_` + userCookie.Version

	token := ctx.Input.Session(sName)

	if token != userCookie.AuxpiToken {
		destoryCookie(ctx)
		if ajaxErrorResp(ctx) {
			return
		}
		//ctx.Redirect(http.StatusFound, site.SiteUrl )
		return
	}

}

var CookieUploadControl = func(ctx *context.Context) {
	if bootstrap.SiteConfig.AllowTourists == true {
		return
	}

	sid := ctx.GetCookie("AuXPI_id")
	id, _ := strconv.Atoi(sid)
	at, _ := ctx.GetSecureCookie(bootstrap.SiteConfig.AuxpiSalt, "AuXPI_at")
	un, _ := ctx.GetSecureCookie(bootstrap.SiteConfig.AuxpiSalt, "AuXPI_uname")
	em, _ := ctx.GetSecureCookie(bootstrap.SiteConfig.AuxpiSalt, "AuXPI_email")

	var userCookie = auxpi.AuxpiCookie{
		UName:      un,
		Email:      em,
		ID:         id,
		Version:    ctx.GetCookie("AuXPI_v"),
		AuxpiToken: at,
	}
	valid := validation.Validation{}
	b, err := valid.Valid(&userCookie)

	if err != nil {
		auxpiLog.SetAWarningLog("COOKIE", err)
	}
	if !b {
		errorInfo := auxpi.RespJson{
			Code: e.ERROR_USER_UN_LOGIN,
			Msg:  e.GetMsg(e.ERROR_USER_UN_LOGIN),
		}
		info, _ := errorInfo.MarshalJSON()
		ctx.Output.Header("Content-Type", "application/json; charset=UTF-8")
		ctx.ResponseWriter.Write(info)
		return
	}

	var sName = `_email_` + userCookie.Email +
		`_user_` + userCookie.UName +
		`_id_` + sid +
		`_version_` + userCookie.Version

	token := ctx.Input.Session(sName)

	if token != userCookie.AuxpiToken {
		errorInfo := auxpi.RespJson{
			Code: e.ERROR_USER_UN_LOGIN,
			Msg:  e.GetMsg(e.ERROR_USER_UN_LOGIN),
		}
		info, _ := errorInfo.MarshalJSON()
		ctx.Output.Header("Content-Type", "application/json; charset=UTF-8")
		ctx.ResponseWriter.Write(info)
		return
	}

}

func destoryCookie(ctx *context.Context) {
	ctx.SetCookie("AuXPI_uname", "", -1)
	ctx.SetCookie("AuXPI_email", "", -1)
	ctx.SetCookie("AuXPI_id", "", -1)
	ctx.SetCookie("AuXPI_v", "", -1)
	ctx.SetCookie("AuXPI_at", "", -1)

	if ctx.GetCookie("AuXPI_Admin-Token") != "" {
		ctx.SetCookie("AuXPI_r", "", -1)
		ctx.SetCookie("AuXPI_Admin-Token", "", -1)
	}
}

func ajaxErrorResp(ctx *context.Context) bool {
	if ctx.Input.IsAjax() {
		errorInfo := auxpi.RespJson{
			Code: e.ERROR_USER_COOKIE,
			Msg:  e.GetMsg(e.ERROR_USER_COOKIE),
		}
		info, _ := errorInfo.MarshalJSON()
		ctx.Output.Header("Content-Type", "application/json; charset=UTF-8")
		ctx.ResponseWriter.Write(info)
		return true
	}
	return false

}
