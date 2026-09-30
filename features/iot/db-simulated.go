package iot

import (
	"sync"

	"github.com/aldesgroup/goald"
)

// ----------------------------------------------------------------------------
// CRUD for user management - simulated for now
// TODO remove this ultimately
// ----------------------------------------------------------------------------

type User struct {
	goald.BusinessObject

	UserFullName string
}

func (thisUser *User) GetUsername() string {
	return thisUser.UserFullName
}

// GetMemberships implements [goald.IUser].
func (thisUser *User) GetMemberships() []goald.IUserGroup {
	panic("unimplemented")
}

var users = []goald.IUser{
	newUser(125, "John Doe"),
	newUser(233, "Jane Doe"),
	newUser(984, "Peter Smith"),
}

func newUser(id int, fullname string) (usr *User) {
	usr = &User{}
	usr.ID = goald.BObjID(id)
	usr.UserFullName = fullname

	return
}

// TODO remove once we have an ORM
func GetUser(id int) goald.IUser {
	for _, usr := range users {
		if int(usr.GetID()) == id {
			return usr
		}
	}

	return nil
}

// ----------------------------------------------------------------------------
// Link Device <-> User
// TODO properly
// ----------------------------------------------------------------------------

func (thisDevice *Device) IsUserAssociated(user goald.IUser) bool {
	for _, associated := range thisDevice.AssociatedUsers {
		if associated.GetID() == user.GetID() {
			return true
		}
	}
	return false
}

// ----------------------------------------------------------------------------
// Device Bootstrap Payloads, and device bootstraps
// TODO properly + job for purging them
// ----------------------------------------------------------------------------

// TODO remove - we won't need them when we'll use the DB
var deviceBoostrapPayloads = map[string]*DeviceBootstrapPayload{}
var deviceBoostrapPayloadsMx sync.Mutex

func (thisPayload *DeviceBootstrapPayload) InsertInDb() {
	deviceBoostrapPayloadsMx.Lock()
	defer deviceBoostrapPayloadsMx.Unlock()
	deviceBoostrapPayloads[thisPayload.Serial+thisPayload.Nonce] = thisPayload
}

func (thisPayload *DeviceBootstrapPayload) ExistsInDb() bool {
	deviceBoostrapPayloadsMx.Lock()
	defer deviceBoostrapPayloadsMx.Unlock()
	return deviceBoostrapPayloads[thisPayload.Serial+thisPayload.Nonce] != nil
}

// ----------------------------------------------------------------------------
// Device Bootstraps
// TODO properly + job for purging them
// ----------------------------------------------------------------------------

var deviceBoostraps = map[string]*DeviceBootstrap{}
var deviceBoostrapsMx sync.Mutex

func GetDeviceBoostrap(serial string) *DeviceBootstrap {
	deviceBoostrapsMx.Lock()
	defer deviceBoostrapsMx.Unlock()

	// getting it, initialising it beforehand if needed
	// TODO fetching it from the DB
	deviceBootstrap := deviceBoostraps[serial]

	if deviceBootstrap == nil {
		deviceBootstrap = &DeviceBootstrap{Status: BootstrapStatusAUTHENTICATEDxONLY}
		deviceBoostraps[serial] = deviceBootstrap
	}

	return deviceBootstrap
}

// Later:

// type Obj struct {
//     ID     string
//     Email  string
//     // other fields...
// }

// // Pseudo: db is *sql.DB, tx optional if you create related rows
// func GetOrCreateByEmail(ctx context.Context, db *sql.DB, email string) (*Obj, error) {
//     // 1) Try to insert
//     _, err := db.ExecContext(ctx,
//         `INSERT INTO users (email) VALUES (@p1)`, email) // Use placeholders for your driver/DB
//     if err != nil {
//         // 2) If unique violation → someone else inserted first; fetch and return
//         if isUniqueViolation(err) {
//             obj := &Obj{}
//             if err := db.QueryRowContext(ctx,
//                 `SELECT id, email FROM users WHERE email = @p1`, email).
//                 Scan(&obj.ID, &obj.Email); err != nil {
//                 return nil, err
//             }
//             return obj, nil
//         }
//         return nil, err
//     }

//     // 3) Insert succeeded → fetch and return (or use RETURNING if supported)
//     obj := &Obj{}
//     if err := db.QueryRowContext(ctx,
//         `SELECT id, email FROM users WHERE email = @p1`, email).
//         Scan(&obj.ID, &obj.Email); err != nil {
//         return nil, err
//     }
//     return obj, nil
// }

// // Map DB errors to "unique violation"
// func isUniqueViolation(err error) bool {
//     // Postgres: SQLSTATE 23505
//     // MySQL: ER_DUP_ENTRY (1062)
//     // SQL Server: 2601 (duplicate key), 2627 (unique index violation)
//     // Implement using driver-specific error inspection.
//     return matchUniqueConstraintError(err)
// }
// ``
