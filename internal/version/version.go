package version

// Version identifies the application build for audit trails and diagnostics.
// Override at build time via: -ldflags "-X emissioncalculator/internal/version.Version=1.2.3".
var Version = "dev"
