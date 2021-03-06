package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/jacoblister/noisefloor/app/assets"
	"github.com/jacoblister/noisefloor/app/audiomodule"
	"github.com/jacoblister/noisefloor/pkg/vdom"
	"github.com/jacoblister/noisefloor/pkg/vfs"
)

var mods modules

//GetAudioProcessor returns the audioProcessor to external javascript
func GetAudioProcessor() audiomodule.AudioProcessor {
	return &mods
}

// App is the main entry point for the application
func App(driver Driver, filesystem vfs.FileSystem) {
	vfs.SetDefaultFS(filesystem)
	mods.dspEngine.Load("phase.xml") // Load dsp engine patch

	hardwareDevices := HardwareDevices{}
	driver.Start(hardwareDevices, &mods)

	go func() {
		mods.Init()
		vdom.SetSVGNamespace()
		vdom.SetHeaderElements([]vdom.Element{
			vdom.MakeElement("link",
				"rel", "stylesheet",
				"type", "text/css",
				"href", "assets/files/style.css"),
		})
		vdom.RenderComponentToDom(&mods)
		vdom.ListenAndServe("/assets/files/", assets.Assets)
	}()

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	<-signalChannel

	driver.Stop(hardwareDevices)
}
