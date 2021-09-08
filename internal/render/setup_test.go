package render

import (
	"encoding/gob"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/gummy789j/bookings/internal/config"
	"github.com/gummy789j/bookings/internal/models"
)

var session *scs.SessionManager

var testApp config.AppConfig

func TestMain(m *testing.M) {
	// what am I going to put in the session
	// use to initialize
	gob.Register(models.Reservation{})

	//  change this when in production
	testApp.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = testApp.InProduction

	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())
}

// Making a Test ResponseWriter by myself

type myWriter struct{}

func (this *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (this *myWriter) Write(data []byte) (int, error) {

	length := len(data)
	return length, nil
}

func (this *myWriter) WriteHeader(statuscode int) {

}
