package cors

type CorsOptions struct {
	AllowOrigins       []string
	AllowHeaders       []string
	AllowMethods       []string
	AllowExposeHeaders []string
	RequestHeaders     []string
	TimingAllowOrigin  []string
	MaxAge             int
	AllowCredentials   bool
	AllowNoOrigin      bool
}

func NewCorsOptions() *CorsOptions {
	return &CorsOptions{
		MaxAge:           -1,
		AllowCredentials: false,
		AllowNoOrigin:    false,
	}
}

var wildcard = []string{"*"}

func Wildcard() []string {
	return wildcard
}
