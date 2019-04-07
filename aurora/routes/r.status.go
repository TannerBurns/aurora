package routes

import "../controllers"

func StatusRoutes(c *controllers.Controller) (r Routes) {
	r = Routes{
		Route{
			"Status",
			"GET",
			c.Session.LiteConfig.Config["api"]["route"] + "/status",
			c.GetStatus,
		},
	}
	return
}
