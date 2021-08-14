package routing

import (
	"down_tip/service"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func Routing(s *ghttp.Server) {
	group := s.Group("/api")
	group.GET("/key_log", func(r *ghttp.Request) {
		r.Response.WriteJson(g.Map{
			"msg":  "获取成功",
			"code": 0,
			"data": service.GetKeyLog(),
		})
	})
}
