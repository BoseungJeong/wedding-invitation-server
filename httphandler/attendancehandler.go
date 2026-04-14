package httphandler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/boseungjeong/wedding-invitation-server/env"
	"github.com/boseungjeong/wedding-invitation-server/sqldb"
	"github.com/boseungjeong/wedding-invitation-server/types"
)

type AttendanceHandler struct {
	http.Handler
}

func (h *AttendanceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		decoder := json.NewDecoder(r.Body)
		var attendance types.AttendanceCreate
		err := decoder.Decode(&attendance)
		if err != nil {
			log.Printf("[ATTENDANCE] decode failed: %v (remote=%s)", err, r.RemoteAddr)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("BadRequest"))
			return
		}

		err = sqldb.CreateAttendance(attendance.Side, attendance.Name, attendance.Meal, attendance.Count)

		if err != nil {
			log.Printf("[ATTENDANCE] insert failed: %v (side=%q name=%q meal=%q count=%d)",
				err, attendance.Side, attendance.Name, attendance.Meal, attendance.Count)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("InternalServerError"))
			return
		}

		log.Printf("[ATTENDANCE] new submission at=%s side=%q name=%q meal=%q count=%d remote=%s",
			time.Now().Format(time.RFC3339),
			attendance.Side,
			attendance.Name,
			attendance.Meal,
			attendance.Count,
			r.RemoteAddr,
		)

		w.Header().Set("Content-Type", "application/json")
	} else if r.Method == http.MethodGet {
		password := r.URL.Query().Get("password")
		if password == "" {
			password = r.Header.Get("X-Admin-Password")
		}

		if env.AdminPassword == "" || password != env.AdminPassword {
			log.Printf("[ATTENDANCE] unauthorized read attempt (remote=%s)", r.RemoteAddr)
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Forbidden"))
			return
		}

		response, err := sqldb.GetAttendance()
		if err != nil {
			log.Printf("[ATTENDANCE] query failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("InternalServerError"))
			return
		}

		pbytes, err := json.Marshal(response)
		if err != nil {
			log.Printf("[ATTENDANCE] marshal failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("InternalServerError"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(pbytes)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed"))
	}
}
