package main

//set on build with -ldflags

// The commit hash that was compiled. This will be filled in by the compiler.
var CommitHash string

// The date of the build
var BuildDate string

// The main version number that is being run at the moment.
const Version = "0.1.0"

// A pre-release marker for the version. If this is "" (empty string)
// then it means that it is a final release. Otherwise, this is a pre-release
// such as "dev" (in development), "beta", "rc1", etc.
const VersionPrerelease = "dev"

// The name of this application
const AppName = "example"
