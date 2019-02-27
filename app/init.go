package app

import (
	"flag"
	"github.com/MonkeyBuisness/screencast/app/service"
	"github.com/revel/revel"
	"log"
	"os"
)

var (
	// AppVersion revel app version (ldflags)
	AppVersion string

	// BuildTime revel app build-time (ldflags)
	BuildTime string

	// flags
	port = flag.Int("p", 9000, "usage port")
)

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		CORSFilter,                    // CORS filter
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.BeforeAfterFilter,       // Call the before and after filter functions
		revel.ActionInvoker,           // Invoke the action.
	}

	revel.OnAppStart(StartupScript)
	revel.OnAppStop(StopScript)
}

var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	fc[0](c, fc[1:])
}

var CORSFilter = func(c *revel.Controller, fc []revel.Filter) {
	c.Response.Out.Header().Set("Access-Control-Allow-Origin", "*")
	c.Response.Out.Header().Set("Access-Control-Allow-Methods", "WS, POST, GET, OPTIONS")
	c.Response.Out.Header().Set("Access-Control-Allow-Headers",
		"Accept, Accept-Language, Content-Type, Content-Length, Accept-Encoding, Authorization")

	if c.Request.Method == "OPTIONS" {
		return
	}

	// next filter
	fc[0](c, fc[1:])
}

func StartupScript() {
	// parse flags
	flag.Parse()
	log.Println(*port)
	log.Println(os.Getenv("hello"))

	// register template functions
	revel.TemplateFuncs["socket_addr"] = func() string {
		return "ws://dc7b5834.ngrok.io/mirror"
	}

	// start screencast service
	service.StartScreencast(service.ScreenCastConfig{
		BitRate: 14,
		Quality: 37,
	})
}

func StopScript() {
	// stop screencast service
	service.StopScreencast()
}
