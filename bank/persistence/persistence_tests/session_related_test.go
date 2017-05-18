package persistence_test

import (
	"github.com/iliyanmotovski/bankv1/bank/domain"
	"github.com/iliyanmotovski/bankv1/bank/persistence"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"testing"
	"time"
)

func TestItStartsUpdatesAndDeletesSessionFromDB(t *testing.T) {
	session, _ := mgo.DialWithTimeout("localhost:27016", time.Second*10)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	sessionStore := persistence.NewMongoSessionStore(*session, "bankTest")

	startInstant, _ := time.Parse("Jan 2, 2006 at 3:04pm (MST)", "Feb 3, 2013 at 7:54pm (EET)")
	updatedInstant, _ := time.Parse("Jan 2, 2006 at 3:04pm (MST)", "Feb 13, 2013 at 7:54pm (EET)")

	testUser := domain.User{UserID: "54321"}

	userSession, err := sessionStore.StartSession(testUser, startInstant)
	if err != nil {
		t.Error(err.Error())
	}

	expected := domain.Session{UserID: "54321", SessionID: userSession.SessionID, Expires: updatedInstant}

	sessionStore.UpdateSession(userSession.SessionID, expected.Expires)

	var result domain.Session
	session.DB("bankTest").C("sessions").Find(bson.M{"userid": testUser.UserID}).One(&result)

	if result != expected {
		t.Errorf("expected session: %v got: %v", expected, result)
	}

	result.SessionID = ""

	sessionStore.DeleteSession(userSession.SessionID)

	sessionStore.Session.DB("bankTest").C("sessions").Find(bson.M{"userid": testUser.UserID}).One(&result)

	if result.SessionID != "" {
		t.Error("expected to have 0 accounts in the DB, but there was more")
	}

	sessionStore.Session.DB("bankTest").C("sessions").RemoveAll(nil)
}

func TestItChecksSessionIsValidAtAGivenTimeInstant(t *testing.T) {
	session, _ := mgo.DialWithTimeout("localhost:27016", time.Second*10)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	sessionStore := persistence.NewMongoSessionStore(*session, "bankTest")

	expires, _ := time.Parse("Jan 2, 2006 at 3:04pm (MST)", "Feb 3, 2013 at 7:54pm (EET)")
	validInstant, _ := time.Parse("Jan 2, 2006 at 3:04pm (MST)", "Feb 3, 2013 at 7:53pm (EET)")
	invalidInstant, _ := time.Parse("Jan 2, 2006 at 3:04pm (MST)", "Feb 3, 2013 at 7:55pm (EET)")

	sessionStore.Session.DB("bankTest").C("sessions").Insert(&domain.Session{Expires: expires, SessionID: "123456", UserID: "654321"})

	var sessionFromDB domain.Session
	session.DB("bankTest").C("sessions").Find(bson.M{"userid": "654321"}).One(&sessionFromDB)

	_, ok := sessionStore.FindSessionAvailableAt("123456", validInstant)
	if ok != true {
		t.Error("expected session to be valid, but it wasn't")
	}

	_, ok = sessionStore.FindSessionAvailableAt("123456", invalidInstant)
	if ok != false {
		t.Error("expected session to be invalid, but it wasn't")
	}

	sessionStore.Session.DB("bankTest").C("sessions").RemoveAll(nil)
}
