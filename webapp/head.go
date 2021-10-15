package webapp

var HeadFrameworkCSSTemplate = []templateHeadLink{
	{
		HRef:        "https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css",
		Rel:         "stylesheet",
		Integrity:   "sha384-1BmE4kWBq78iYhFldvKuhfTAU6auU8tT94WrHftjDbrCEXSU1oBoqyl2QvZ6jIW3",
		CrossOrigin: "anonymous",
	},
	{
		HRef:        "/static/fontawesome-free-5.15.4-web/css/all.min.css",
		Rel:         "stylesheet",
		CrossOrigin: "anonymous",
	},
}

var HeadCSSTemplate = []templateHeadLink{
	{
		HRef: "/static/css/default.css",
		Rel:  "stylesheet",
	},
}

var FooterScriptTemplate = []templateScript{
	{
		Src:         "https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/js/bootstrap.bundle.min.js",
		Integrity:   "sha384-ka7Sk0Gln4gmtz2MlQnikT1wXgYsOg+OMhuP+IlRH9sENBO0LRn5q+8nbTov4+1p",
		CrossOrigin: "anonymous",
	},
}
