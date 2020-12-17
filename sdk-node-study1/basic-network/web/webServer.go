
package web



import (

	"net/http"

	"fmt"

	"basic-network/web/controller"

)



func  WebStart(app *controller.Application)  {



	fs := http.FileServer(http.Dir("web/static"))

	http.Handle("/static/", http.StripPrefix("/static/", fs))



	http.HandleFunc("/company/add", app.Addcompanyinfo)



	http.HandleFunc("/company/search", app.QueryInfo)

	http.HandleFunc("/PurchaseRecord/search", app.GetPurchaseRecord)
	http.HandleFunc("/TicketRecord/search", app.GetTicket)
	http.HandleFunc("/PurchaseRecord/add", app.AddPurchaseRecord)
	http.HandleFunc("/TicketRecord/add", app.AddTicket)
	http.HandleFunc("/History/search", app.GetHistoryForKey)


	fmt.Println("启动Web服务, 监听端口号: 8888")



	err := http.ListenAndServe(":8888", nil)

	if err != nil {

		fmt.Println("启动Web服务错误")

	}



}
