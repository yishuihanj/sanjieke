cur_dir=$(shell pwd)
bin_dir=$(cur_dir)/bin/
BUILD_GCFLAG="-N -l"

BUILD_CMD=go build
BUILD_FLAGS=-gcflags $(BUILD_GCFLAG)  -trimpath -a -o $(bin_dir)

sdl:
	$(BUILD_CMD) $(BUILD_FLAGS)/sdl.exe ./