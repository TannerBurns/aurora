package routes

import "../controllers"

func AuthenticationRoutes(c *controllers.Controller) (r Routes) {
	r = Routes{
		Route{
			"GetToken",
			"Post",
			c.Session.LiteConfig.Config["api"]["route"] + "/login",
			c.GetToken,
		},
	}
	return
}
