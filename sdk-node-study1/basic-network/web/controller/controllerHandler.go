package controller

import (
	"net/http"

	"basic-network/service"
	"fmt"
)

type Application struct {
	Fabric *service.ServiceSetup
}

func (app *Application) Addcompanyinfo(w http.ResponseWriter, r *http.Request) {

	// 上传提交数据

	company_name := r.FormValue("company_name")
	location := r.FormValue("location")
	restaurant := r.FormValue("restaurant")
	legalperson := r.FormValue("legalperson")
	linkname := r.FormValue("linkname")
	tel := r.FormValue("tel")
	add0 := []string{company_name, location, restaurant, legalperson, linkname, tel}
	// 调用业务层, 反序列化

	res, err := app.Fabric.AddCompany(add0...)
	if err != nil {
		fmt.Fprintln(w,err.Error())
	} else {

		// 响应客户端

		fmt.Fprintln(w, res)
	}
}

func (app *Application) AddTicket(w http.ResponseWriter, r *http.Request) {

	//上传提交数据

	company_name := r.FormValue("company_name")
	ticket := r.FormValue("ticket_record")
	add0 := []string{company_name,ticket}
	// 调用业务层, 反序列化

	res, err := app.Fabric.AddTicket(add0...)
	if err != nil {
		fmt.Fprintln(w,err.Error())
	} else {

		// 响应客户端

		fmt.Fprintln(w, res)
	}
}
func (app *Application) AddPurchaseRecord(w http.ResponseWriter, r *http.Request) {

	// 获取提交数据

	company_name := r.FormValue("company_name")
	purchase := r.FormValue("purchase_record")
	add0 := []string{company_name,purchase}
	// 调用业务层, 反序列化

	res, err := app.Fabric.AddPurchaseRecord(add0...)
	if err != nil {
		fmt.Fprintln(w,err.Error())
	} else {

		// 响应客户端

		fmt.Fprintln(w, res)
	}
}
func (app *Application) GetPurchaseRecord(w http.ResponseWriter, r *http.Request) {

	// 获取提交数据

	company_name := r.FormValue("company_name")

	res, err := app.Fabric.GetPurchaseRecord(company_name)
	if err != nil {
		fmt.Fprintln(w,err.Error())
	} else {

		// 响应客户端

		fmt.Fprintln(w, res)
	}
}
func (app *Application) GetTicket(w http.ResponseWriter, r *http.Request) {

	// 获取提交数据

	company_name := r.FormValue("company_name")

	res, err := app.Fabric.GetTicket(company_name)
	if err != nil {
		fmt.Fprintln(w,err.Error())
	} else {

		// 响应客户端

		fmt.Fprintln(w, res)
	}
}
func (app *Application) GetHistoryForKey(w http.ResponseWriter, r *http.Request) {

	// 获取提交数据

	company_name := r.FormValue("company_name")

	res, err := app.Fabric.GetHistoryForKey(company_name)
	if err != nil {
		fmt.Fprintln(w,err.Error())
	} else {

		// 响应客户端

		fmt.Fprintln(w, res)
	}
}
func (app *Application) QueryInfo(w http.ResponseWriter, r *http.Request) {

	// 获取提交数据

	name := r.FormValue("company_name")

	msg, err := app.Fabric.QueryByName(name)

	if err != nil {
		fmt.Fprintln(w,err.Error())
	} else {
		fmt.Fprintln(w, msg)
	}
}
