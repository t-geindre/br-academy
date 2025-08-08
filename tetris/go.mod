module tetris

go 1.24.1

require (
	component v0.0.0-00010101000000-000000000000
	control v0.0.0-00010101000000-000000000000
	debug v0.0.0-00010101000000-000000000000
	github.com/hajimehoshi/ebiten/v2 v2.8.8
	layout v0.0.0-00010101000000-000000000000
	pool v0.0.0-00010101000000-000000000000
	ui v0.0.0-00010101000000-000000000000
)

require (
	github.com/ebitengine/gomobile v0.0.0-20240911145611-4856209ac325 // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/oto/v3 v3.3.3 // indirect
	github.com/ebitengine/purego v0.8.0 // indirect
	github.com/go-text/typesetting v0.2.0 // indirect
	github.com/hajimehoshi/go-mp3 v0.3.4 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	golang.org/x/image v0.20.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.18.0 // indirect
)

replace (
	component => ../libs/component
	control => ../libs/control
	debug => ../libs/debug
	layout => ../libs/layout
	pool => ../libs/pool
	ui => ../libs/ui
)
