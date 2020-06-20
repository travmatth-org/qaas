package handlerfuncs

import (
	"io/ioutil"
	"net/http"

	"github.com/Travmatth/faas/internal/config"
	"github.com/rs/zerolog/hlog"
)

type Home struct {
	file   []byte
	config *config.Config
}

func NewHome(c *config.Config) (*Home, error) {
	filename := c.GetIndexHtml()
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Home{f, c}, nil
}

func (h *Home) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if bytes, err := w.Write(h.file); err != nil {
		hlog.FromRequest(r).Error().
			Err(err).
			Msg("Error completing request")
	} else {
		hlog.FromRequest(r).Info().
			Int("sent", bytes).
			Msg("Request completed")
	}
}
