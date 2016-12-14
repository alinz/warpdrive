package warpify

const (
	_ EventKind = 0

	// NoUpdate means there is no update available at the moment
	NoUpdate
	// UpdateAvailable there is an update available for download
	UpdateAvailable
	// UpdateDownloading downloading has started
	UpdateDownloading
	// UpdateDownloaded downlaoing has completed
	// at this moment, a callback from objective c or java should restart the app
	UpdateDownloaded
)

// Setup we need to setup the app
func Setup(bundleVersion, bundlePath, documentPath string) {
	conf.bundleVersion = bundleVersion
	conf.bundlePath = bundlePath
	conf.documentPath = bundlePath

	conf.pubSub = newPubSub()
}

// SourcePath returns the proper path for react-native app to start the process
func SourcePath() string {
	return ""
}

// Subscribe this is a easy to use method to expose to objective-c and jave
// so they can bind their callbacks to known EventKinds
func Subscribe(eventKind EventKind, callback Callback) {
	conf.pubSub.Subscribe(eventKind, callback)
}

// Unsubscribe as it stands, it unsubscribes the any associate
// callback to specific event type. Mainly it's being used for clean up
func Unsubscribe(eventKind EventKind) {
	conf.pubSub.Unsubscribe(eventKind)
}
