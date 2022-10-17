package routes

import "github.com/panhdjf/server_management_system/controllers"

type Excel_route_controller struct {
	Excelcontroller controllers.ExcelController
}

func New_route_excel_controller(excelcontroller controllers.ExcelController) Excel_route_controller {
	return Excel_route_controller{excelcontroller}
}

func (c *Excel_route_controller) Excel_Route() {}
