package routes

import "../controllers"

func UserRoutes(c *controllers.Controller) (r Routes) {
	r = Routes{
		Route{
			"CreateUser",
			"Post",
			c.Session.LiteConfig.Config["api"]["route"] + "/user",
			c.CreateUser,
		},
		Route{
			"ReadUser",
			"Get",
			c.Session.LiteConfig.Config["api"]["route"] + "/user",
			controllers.AuthenticationMiddleware(c.ReadUser),
		},
		Route{
			"UpdateUser",
			"Put",
			c.Session.LiteConfig.Config["api"]["route"] + "/user/{id}",
			controllers.AuthenticationMiddleware(c.UpdateUser),
		},
		/*
			Route{
				"DeleteUser",
				"Delete",
				c.Session.LiteConfig.Config["api"]["route"] + "/user/{id}",
				controllers.AuthenticationMiddleware(c.DeleteUser),
			},
		*/
		Route{
			"CreateUserTask",
			"Post",
			c.Session.LiteConfig.Config["api"]["route"] + "/task",
			controllers.AuthenticationMiddleware(c.CreateTask),
		},
		Route{
			"GetTask",
			"Get",
			c.Session.LiteConfig.Config["api"]["route"] + "/task/{id}",
			controllers.AuthenticationMiddleware(c.GetTask),
		},
		Route{
			"CreateComment",
			"Post",
			c.Session.LiteConfig.Config["api"]["route"] + "/task/{id}/comment",
			controllers.AuthenticationMiddleware(c.CreateComment),
		},
		Route{
			"UpdateComment",
			"Post",
			c.Session.LiteConfig.Config["api"]["route"] + "/task/{tid}/" +
				"comment/{cid}",
			controllers.AuthenticationMiddleware(c.UpdateComment),
		},
		Route{
			"CreateTag",
			"Post",
			c.Session.LiteConfig.Config["api"]["route"] + "/task/{tid}/" +
				"comment/{cid}/tag",
			controllers.AuthenticationMiddleware(c.CreateTag),
		},
		Route{
			"DeleteTag",
			"Delete",
			c.Session.LiteConfig.Config["api"]["route"] + "/task/{tid}/" +
				"comment/{cid}/tag/{taid}",
			controllers.AuthenticationMiddleware(c.DeleteTag),
		},
	}
	return
}
