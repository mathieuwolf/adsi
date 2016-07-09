package adsi

import (
	"sync"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/scjalliance/comshim"
	"gopkg.in/adsi.v0/api"
)

// ADSI Objects of LDAP:  https://msdn.microsoft.com/library/aa772208
// ADSI Objects of WinNT: https://msdn.microsoft.com/library/aa772211

// Object provides access to Active Directory objects.
type Object struct {
	object
}

// NewObject returns an object that manages the given COM interface.
func NewObject(iface *api.IADs) *Object {
	comshim.Add(1)
	return &Object{object{iface: iface}}
}

type object struct {
	m     sync.RWMutex
	iface *api.IADs
}

func (o *object) closed() bool {
	return (o.iface == nil)
}

// Close will release resources consumed by the object. It should be
// called when the object is no longer needed.
func (o *object) Close() {
	o.m.Lock()
	defer o.m.Unlock()
	if o.closed() {
		return
	}
	defer comshim.Done()
	run(func() error {
		o.iface.Release()
		return nil
	})
	// FIXME: What happens if the run returns an error?
	o.iface = nil
}

// Name retrieves the name of the object.
func (o *object) Name() (name string, err error) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.closed() {
		return "", ErrClosed
	}
	err = run(func() error {
		name, err = o.iface.Name()
		if err != nil {
			return err
		}
		return nil
	})
	return
}

// Class retrieves the class of the object.
func (o *object) Class() (class string, err error) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.closed() {
		return "", ErrClosed
	}
	err = run(func() error {
		class, err = o.iface.Class()
		if err != nil {
			return err
		}
		return nil
	})
	return
}

// GUID retrieves the globally unique identifier of the object.
func (o *object) GUID() (guid *ole.GUID, err error) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.closed() {
		return nil, ErrClosed
	}
	err = run(func() error {
		var sguid string
		sguid, err = o.iface.GUID()
		if err != nil {
			return err
		}

		guid = ole.NewGUID(sguid)
		if guid == nil {
			return ErrInvalidGUID
		}

		return nil
	})
	return
}

// Path retrieves the fully qualified path of the object.
func (o *object) Path() (path string, err error) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.closed() {
		return "", ErrClosed
	}
	err = run(func() error {
		path, err = o.iface.AdsPath()
		if err != nil {
			return err
		}
		return nil
	})
	return
}

// Parent retrieves the fully qualified path of the object's parent.
func (o *object) Parent() (path string, err error) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.closed() {
		return "", ErrClosed
	}
	err = run(func() error {
		path, err = o.iface.Parent()
		if err != nil {
			return err
		}
		return nil
	})
	return
}

// Schema retrieves the fully qualified path of the object's schema class
// object.
func (o *object) Schema() (path string, err error) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.closed() {
		return "", ErrClosed
	}
	err = run(func() error {
		path, err = o.iface.Schema()
		if err != nil {
			return err
		}
		return nil
	})
	return
}

// ToContainer attempts to acquire a container interface for the object.
func (o *object) ToContainer() (c *Container, err error) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.closed() {
		return nil, ErrClosed
	}
	err = run(func() error {
		idispatch, err := o.iface.QueryInterface(api.IID_IADsContainer)
		if err != nil {
			return err
		}
		iface := (*api.IADsContainer)(unsafe.Pointer(idispatch))
		c = NewContainer(iface)
		return nil
	})
	return
}

// ToComputer attempts to acquire a computer interface for the object.
func (o *object) ToComputer() (c *Computer, err error) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.closed() {
		return nil, ErrClosed
	}
	err = run(func() error {
		idispatch, err := o.iface.QueryInterface(api.IID_IADsComputer)
		if err != nil {
			return err
		}
		iface := (*api.IADsComputer)(unsafe.Pointer(idispatch))
		c = NewComputer(iface)
		return nil
	})
	return
}

// ToGroup attempts to acquire a group interface for the object.
func (o *object) ToGroup() (g *Group, err error) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.closed() {
		return nil, ErrClosed
	}
	err = run(func() error {
		idispatch, err := o.iface.QueryInterface(api.IID_IADsGroup)
		if err != nil {
			return err
		}
		iface := (*api.IADsGroup)(unsafe.Pointer(idispatch))
		g = NewGroup(iface)
		return nil
	})
	return
}
