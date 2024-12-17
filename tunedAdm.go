package tunedAdm

import (
    "fmt"
    "errors"
    "time"
    "github.com/godbus/dbus/v5"
)

const (
    systemd1_dest = "org.freedesktop.systemd1"
    systemd1_path = "/org/freedesktop/systemd1/unit/tuned_2eservice"
    systemd1_interface = "org.freedesktop.systemd1.Unit"
    systemd1_properties = "org.freedesktop.DBus.Properties.Get"
    tuned_dest = "com.redhat.tuned"
    tuned_path = "/Tuned"
    tuned_interface = "com.redhat.tuned.control"
)

func connect() (*dbus.Conn, error) {
    c, err := dbus.ConnectSystemBus()
    if err != nil {
        return nil, errors.New("failed to connect to session bus: " + fmt.Sprint(err))
    }
    return c, nil
}

func Disable() (error) {
    if err := serviceRestore(); err != nil {
        return err
    }
    c, err := connect()
    if err != nil {
        return err
    }
    defer c.Close()
    obj := c.Object(tuned_dest, tuned_path)
    if err := obj.Call(tuned_interface + ".disable", 0).Store(); err != nil {
        return err
    }
    return nil
}

func Profiles() ([]string, error) {
    var resp []string
    if err := serviceRestore(); err != nil {
        return nil, err
    }
    c, err := connect()
    if err != nil {
        return nil, err
    }
    defer c.Close()
    obj := c.Object(tuned_dest, tuned_path)
    if err := obj.Call(tuned_interface + ".profiles", 0).Store(&resp); err != nil {
        return nil, err
    }
    return resp, nil
}

func IsRunning() (bool, error) {
    var resp bool
    if err := serviceRestore(); err != nil {
        return false, err
    }
    c, err := connect()
    if err != nil {
        return false, err
    }
    defer c.Close()
    obj := c.Object(tuned_dest, tuned_path)
    if err := obj.Call(tuned_interface + ".is_running", 0).Store(&resp); err != nil {
        return false, err
    }
    return resp, nil
}

func ActiveProfile() (string, error) {
    var resp string
    if err := serviceRestore(); err != nil {
        return "", err
    }
    c, err := connect()
    if err != nil {
        return "", err
    }
    defer c.Close()
    obj := c.Object(tuned_dest, tuned_path)
    if err := obj.Call(tuned_interface + ".active_profile", 0).Store(&resp); err != nil {
        return "", err
    }
    return resp, nil
}

func VerifyProfile() (bool, error) {
    var resp bool
    if err := serviceRestore(); err != nil {
        return false, err
    }
    c, err := connect()
    if err != nil {
        return false, err
    }
    defer c.Close()
    obj := c.Object(tuned_dest, tuned_path)
    if err := obj.Call(tuned_interface + ".verify_profile", 0).Store(&resp); err != nil {
        return false, err
    }
    return resp, nil
}

func SwitchProfile(profile string) (error) {
    var resp []interface{}
    if err := serviceRestore(); err != nil {
        return err
    }
    c, err := connect()
    if err != nil {
        return err
    }
    defer c.Close()
    obj := c.Object(tuned_dest, tuned_path)
    if err := obj.Call(tuned_interface + ".switch_profile", 0, profile).Store(&resp); err != nil {
        return err
    }

    if resp[0] == false {
        return errors.New(resp[1].(string))
    } else {
        return nil
    }
}

func AutoProfile() (error) {
    var resp []interface{}
    if err := serviceRestore(); err != nil {
        return err
    }
    c, err := connect()
    if err != nil {
        return err
    }
    defer c.Close()
    obj := c.Object(tuned_dest, tuned_path)
    if err := obj.Call(tuned_interface + ".auto_profile", 0).Store(&resp); err != nil {
        return err
    }

    if resp[0] == false {
        return errors.New(resp[1].(string))
    } else {
        return nil
    }
}

func RecommendProfile() (string, error) {
    var resp string
    if err := serviceRestore(); err != nil {
        return "", err
    }
    c, err := connect()
    if err != nil {
        return "", err
    }
    defer c.Close()
    obj := c.Object(tuned_dest, tuned_path)
    if err := obj.Call(tuned_interface + ".recommend_profile", 0).Store(&resp); err != nil {
        return "", err
    }
    return resp, nil
}

func ServiceStatus() (string, error) {
    var resp string
    c, err := connect()
    if err != nil {
        return "", err
    }
    defer c.Close()
    obj := c.Object(systemd1_dest, systemd1_path)
    if err := obj.Call(systemd1_properties, 0, "org.freedesktop.systemd1.Unit", "ActiveState").Store(&resp); err != nil {
        return "", err
    }
    return resp, nil
}

func ServiceStart() (error) {
    var resp string
    c, err := connect()
    if err != nil {
        return err
    }
    defer c.Close()
    obj := c.Object(systemd1_dest, systemd1_path)
    if err := obj.Call(systemd1_interface + ".Start", 0, "replace").Store(&resp); err != nil {
        return err
    }
    return nil
}

func ServiceStop() (error) {
    var resp string
    c, err := connect()
    if err != nil {
        return err
    }
    defer c.Close()
    obj := c.Object(systemd1_dest, systemd1_path)
    if err := obj.Call(systemd1_interface + ".Stop", 0, "replace").Store(&resp); err != nil {
        return err
    }
    return nil
}

func serviceRestore() (error) {
    if respStatus, _ := ServiceStatus(); respStatus == "inactive" {
        if err := ServiceStart(); err != nil {
            return err
        } else {
            for i := 0; i < 10; i++ {
                if status, _ := ServiceStatus(); status == "active" {
                    return nil
                }
                time.Sleep(1 * time.Second)
            }
            return errors.New("timeout starting tuned")
        }
    }
    return nil
}
