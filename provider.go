package gothicComponents

import "os"

// sigV4Enabled mirrors the GOTHIC_MODE read-once-into-a-package-var idiom
// (see core/router/setup.go). GOTHIC_PROVIDER is read a single time at process
// start; on AWS the app must sign every htmx request with SigV4.
var sigV4Enabled = os.Getenv("GOTHIC_PROVIDER") == "AWS"

// SigV4Enabled reports whether the app is deployed on AWS (read once at process start,
// mirroring the GOTHIC_MODE pattern). When true, RuntimeScripts emits the WASM AWS
// request-signing enable flag + boot-race coordination shim.
func SigV4Enabled() bool { return sigV4Enabled }
