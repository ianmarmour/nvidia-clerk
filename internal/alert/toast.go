package alert

import (
	"fmt"
	"os/exec"

	"github.com/ianmarmour/nvidia-clerk/third_party/toast"
)

var execCommand = exec.Command

func linuxToast(name string) error {
	err := execCommand("notify-send", "NVIDIA Clerk", fmt.Sprintf("%s Is ready for checkout", name), "-u", "critical").Start()
	if err != nil {
		return err
	}

	return nil
}

func darwinToast(name string) error {
	notification := fmt.Sprintf("display notification \"NVIDIA Clerk\" with title \"NVIDIA Clerk Inventory Alert\" subtitle \"%s\"", fmt.Sprintf("%s Is ready for checkout", name))
	err := execCommand("osascript", "-e", notification).Start()
	if err != nil {
		return err
	}

	return nil
}

func windowsToast(name string) error {
	notification := toast.Notification{
		AppID:    "NVIDIA Clerk",
		Title:    "NVIDIA Clerk Inventory Alert",
		Message:  fmt.Sprintf("%s Is ready for checkout", name),
		Duration: "long",
	}

	err := notification.Push()
	if err != nil {
		return err
	}

	return nil
}

// SendToast Sends a Toast alert for desktop notifications.
func SendToast(os string, name string) error {
	var err error

	switch os {
	case "linux":
		err = linuxToast(name)
	case "windows":
		err = windowsToast(name)
	case "darwin":
		err = darwinToast(name)
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return err
}
