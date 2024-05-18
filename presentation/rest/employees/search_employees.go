package employees

import "github.com/gin-gonic/gin"

func (ctrl *EmployeeController) SearchEmployees(c *gin.Context) {
	query := c.Query("query")
	if len(query) < 4 {
		c.JSON(400, gin.H{"error": "query should be at least 4 characters long"})
		return
	}

	employees, err := ctrl.EmployeeService.SearchEmployees(c, query)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, employees)
}
