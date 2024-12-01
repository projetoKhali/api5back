package server

import (
	"fmt"
	"net/http"

	"api5back/ent"
	"api5back/src/service"

	"github.com/gin-gonic/gin"
)

func HiringProcessDashboard(
	engine *gin.Engine,
	dbClient *ent.Client,
	dwClient *ent.Client,
) {
	v1 := engine.Group("/api/v1")
	{
		hiringProcess := v1.Group("/hiring-process")
		{
			hiringProcess.POST("/dashboard", Dashboard(dwClient))
			hiringProcess.POST("/table", VacancyTable(dwClient))
		}

		suggestions := v1.Group("/suggestions")
		{
			suggestions.GET("/recruiter", UserList(dwClient))
			suggestions.POST("/process", HiringProcessList((dwClient)))
			suggestions.POST("/vacancy", VacancyList(dwClient))
			suggestions.GET("/department", ListDepartments(dbClient))
		}
		authentication := v1.Group("/authentication")
		{
			authentication.GET("/users", ListUsers(dbClient))
			authentication.POST("/login", LoginUser(dbClient))
			authentication.POST("/create", CreateUser(dbClient))
		}
		groupAccess := v1.Group("/group-access")
		{
			groupAccess.GET("", ListGroupAcess(dbClient))
			groupAccess.POST("", CreateGroupAcess(dbClient))
		}
	}
}

// Dashboard godoc
// @Summary dashboard
// @Schemes
// @Description show dashboard
// @Tags hiring-process
// @Accept json
// @Param body body service.FactHiringProcessFilter true "Metrics filter"
// @Produce json
// @Success 200 {string} Dashboard
// @Router /hiring-process/dashboard [post]
func Dashboard(
	dwClient *ent.Client,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		// TODO: change to pointer
		var dashboardMetricsFilter service.FactHiringProcessFilter
		if err := c.ShouldBindJSON(&dashboardMetricsFilter); err != nil {
			c.JSON(http.StatusBadRequest, DisplayError(err))
			return
		}

		metricsData, err := service.GetMetrics(
			c, dwClient, dashboardMetricsFilter,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, DisplayError(err))
			return
		}

		c.JSON(http.StatusOK, metricsData)
	}
}

func TableData(
	dbClient *ent.Client,
	dwClient *ent.Client,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		var userIDs []int

		// Parse the body for user IDs
		if err := c.ShouldBindJSON(&userIDs); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
	}
}

// UserList godoc
// @Summary List users
// @Schemes
// @Description Return a list of users with id and name
// @Tags suggestions
// @Accept json
// @Produce json
// @Success 200 {array} model.Suggestion
// @Router /suggestions/recruiter/ [get]
func UserList(dwClient *ent.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		users, err := service.GetUsers(c, dwClient)
		if err != nil {
			c.JSON(http.StatusInternalServerError, DisplayError(err))
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

// HiringProcessList godoc
// @Summary List hiring processes
// @Schemes
// @Description Return a list of hiring processes with id and title
// @Tags suggestions
// @Accept json
// @Param body body []int true "User IDs"
// @Produce json
// @Success 200 {array} model.Suggestion
// @Router /suggestions/process [post]
func HiringProcessList(
	dbClient *ent.Client,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		var userIDs *[]int

		// Parse the body for user IDs
		if err := c.ShouldBindJSON(&userIDs); err != nil {
			c.JSON(http.StatusBadRequest, DisplayError(err))
			return
		}

		processes, err := service.ListHiringProcesses(
			c, dbClient,
			userIDs,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, DisplayError(err))
			return
		}

		c.JSON(http.StatusOK, processes)
	}
}

// HiringProcessList godoc
// @Summary List hiring processes
// @Schemes
// @Description Return a list of hiring processes with id and title
// @Tags suggestions
// @Accept json
// @Param body body []int false "User IDs"
// @Produce json
// @Success 200 {array} model.Suggestion
// @Router /suggestions/vacancy [post]
func VacancyList(
	dwClient *ent.Client,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		var processesIds *[]int
		if err := c.ShouldBindJSON(&processesIds); err != nil {
			c.JSON(http.StatusBadRequest, DisplayError(fmt.Errorf("error: Invalid request body")))
			return
		}

		vacancies, err := service.GetVacancySuggestions(
			c, dwClient,
			processesIds,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, DisplayError(err))
			return
		}

		c.JSON(http.StatusOK, vacancies)
	}
}

// HiringProcessList godoc
// @Summary List hiring processes
// @Schemes
// @Description Return a list of vacancies with summarized information
// @Tags hiring-process
// @Accept json
// @Param body body service.FactHiringProcessFilter true "Metrics filter"
// @Produce json
// @Success 200 {array} model.Suggestion
// @Router /hiring-process/table [post]
func VacancyTable(
	dwClient *ent.Client,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		var filter service.FactHiringProcessFilter
		if err := c.ShouldBindJSON(&filter); err != nil {
			c.JSON(http.StatusBadRequest, DisplayError(fmt.Errorf("error Invalid request body")))
			return
		}

		vacancies, err := service.GetVacancyTable(
			c, dwClient,
			filter,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, DisplayError(err))
			return
		}

		c.JSON(http.StatusOK, vacancies)
	}
}

// ListDepartments godoc
// @Summary List departments
// @Schemes
// @Description Return a list of departments with id and title
// @Tags departments
// @Produce json
// @Success 200 {array} model.Suggestion
// @Router /suggestions/departments [get]
func ListDepartments(
	client *ent.Client,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		departments, err := service.ListDepartments(c, client)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to list departments",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, departments)
	}
}

// CreateGroupAcess godoc
// @Summary Create a new group access
// @Schemes
// @Description Create a new group access with name and related departments
// @Tags group_acess
// @Accept json
// @Produce json
// @Param body body service.CreateGroupAcessRequest true "Group Access Info"
// @Success 201 {object} ent.GroupAcess
// @Router /group-access [post]
func CreateGroupAcess(client *ent.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		var request service.CreateGroupAcessRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		group, err := service.CreateGroupAcess(c, client, request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, group)
	}
}

// ListGroupAcess godoc
// @Summary List group access with departments
// @Schemes
// @Description Return a list of group access with id, name, and departments
// @Tags group_acess
// @Produce json
// @Success 200 {array} service.GroupAcessReturn
// @Router /group-access [get]
func ListGroupAcess(client *ent.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		groups, err := service.GetGroupAcessWithDepartments(c, client)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, groups)
	}
}

// ListUsers godoc
// @Summary List all users
// @Description Retrieve a list of all users with ID, name, and email
// @Tags authentication
// @Accept json
// @Produce json
// @Success 200 {array} ent.Authentication
// @Failure 500 {object} gin.H{"error": "description of the error"}
// @Router /authentication/users [get]
func ListUsers(client *ent.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		users, err := client.Authentication.Query().All(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

// LoginUser godoc
// @Summary User login
// @Description Authenticate a user with email and password
// @Tags authentication
// @Accept json
// @Produce json
// @Param body body map[string]string true "User credentials: {email, password}"
// @Success 200 {object} gin.H{"message": "login successful", "token": "jwt_token"}
// @Failure 400 {object} gin.H{"error": "Invalid email or password"}
// @Router /authentication/login [post]
func LoginUser(client *ent.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Parse the request body
		var creds map[string]string
		if err := c.ShouldBindJSON(&creds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Validate required fields
		email, emailOk := creds["email"]
		password, passOk := creds["password"]
		if !emailOk || !passOk {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing email or password"})
			return
		}

		// Call the Login function from the service
		loginResponse, err := service.Login(c.Request.Context(), client, service.LoginRequest{
			Email:    email,
			Password: password,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Generate a JWT or session token (placeholder)
		token := "jwt_token_placeholder"

		// Return the login response
		c.JSON(http.StatusOK, gin.H{
			"message": "login successful",
			"user":    loginResponse,
			"token":   token,
		})
	}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with name, email, password, and group ID
// @Tags authentication
// @Accept json
// @Produce json
// @Param body body service.CreateUserRequest true "User info: {name, email, password, groupID}"
// @Success 201 {object} ent.Authentication
// @Failure 400 {object} gin.H{"error": "Invalid request body"}
// @Failure 500 {object} gin.H{"error": "Error message"}
// @Router /authentication/create [post]
func CreateUser(client *ent.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		var request service.CreateUserRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		user, err := service.CreateUser(c, client, request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, user)
	}
}
