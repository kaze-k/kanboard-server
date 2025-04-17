package admin

import (
	"server/internal/app/admin/handlers"
	"server/internal/constant"
	"server/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func InitApi(router *gin.RouterGroup) {
	admin := router.Group("/admin")

	public := admin.Group("/")

	user := admin.Group("/:id")
	user.Use(middlewares.Auth(constant.ADMIN_TOKEN))

	publicHandler := handlers.NewPublicHandler()
	{
		user.GET("/dashboard", publicHandler.Dashboard)
		user.GET("/getRecentTasks", publicHandler.GetRecentTasks)

		public.GET("/captcha", publicHandler.GetCaptcha)
	}

	wsHandler := handlers.NewWSHandler()
	{
		user.GET("/ws", wsHandler.WebSocket)
	}

	messageHandler := handlers.NewMessageHandler()
	{
		user.GET("/getAllMsgs", messageHandler.GetAllMsgs)
		user.GET("/getUnReadMsgs", messageHandler.GetUnReadMsgs)
		user.GET("/getReadedMsgs", messageHandler.GetReadedMsgs)
		user.POST("/markReadMsg", messageHandler.MarkReadMsg)
		user.DELETE("/deleteMsg", messageHandler.DeleteMsg)
	}

	userHandler := handlers.NewUserHandler()
	{

		admin.POST("/login", userHandler.Login)
		// admin.POST("/register", userHandler.Register)

		user.POST("/createUser", userHandler.CreateUser)
		user.GET("/getUsers", userHandler.GetUsers)
		user.POST("/getUser/", userHandler.GetUser)
		user.PUT("/updateUser", userHandler.UpdateUser)
		user.POST("/searchUser", userHandler.SearchUser)
		user.DELETE("/deleteUser", userHandler.DeleteUser)
		user.PUT("/changePassword", userHandler.ChangePassword)
		user.GET("/getMembers", userHandler.GetMembers)
		user.POST("/uploadAvatar", userHandler.UploadAvatar)

	}

	projectHandler := handlers.NewProjectHandler()
	{

		user.POST("/createProject", projectHandler.CreateProject)
		user.GET("/getProjectList", projectHandler.GetProjectList)
		user.POST("/getProjectById", projectHandler.GetProject)
		user.PUT("/updateProject", projectHandler.UpdateProject)
		user.DELETE("/deleteProject", projectHandler.DeleteProject)
		user.POST("/addProjectMember", projectHandler.AddProjectMember)
		user.POST("/removeProjectMember", projectHandler.RemoveProjectMember)
		user.POST("/setProjectAssignee", projectHandler.SetProjectAssignee)

	}

	taskHandler := handlers.NewTaskHandler()
	{
		user.GET("/getProjectTask", taskHandler.GetProjectTask)
	}
}
