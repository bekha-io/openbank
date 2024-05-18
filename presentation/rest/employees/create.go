package employees

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ctrl *EmployeeController) CreateEmployee(c *gin.Context) {
	type req struct {
		Email      string `json:"email" binding:"required"`
		Password   string `json:"password" binding:"required"`
		FirstName  string `json:"first_name" binding:"required"`
		LastName   string `json:"last_name" binding:"required"`
		MiddleName string `json:"middle_name"`
	}

	var r req
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := ctrl.EmployeeService.CreateEmployee(c, r.Email, r.Password, r.FirstName, r.LastName, r.MiddleName)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, "employee created")
}
