# tuhirender

The `tuhirender` utility processes
[Tuhi](https://github.com/tuhiproject/tuhi)'s JSON output files and
produces PNG or GIF output files from them. Some example images can be
found at https://blog.gnoack.org/post/tuhirender/.

# Installation

To install, run:

```
go get github.com/gnoack/tuhirender/cmd/tuhirender
```

The self-contained `tuhirender` binary is saved at `~/go/bin/tuhirender`.

# Example: Convert to fit

Create a new image of width 100 pixels; scale the input image to fit.

```
tuhirender -width 100 -fit -o img100.png < input.json
```

# Example: Convert to an animated GIF

```
tuhirender -width 200 -fit --format gif -o out.gif < in/1593622352.json
```

Note: Large drawings can take some time to convert and the output
files can become large.

# Example: Convert in batch using a makefile

```
# Create a fresh working directory
mkdir tuhipics; cd tuhipics
# Link "in" to the Tuhi directory with JSON files.
ln -s ~/.local/share/tuhi/12:34:56:78:90:ab in
# Create an empty "out" directory for the PNG images.
mkdir out
```

Create the following `Makefile`:

```
OUTS := $(patsubst in/%.json,out/%.png,$(wildcard in/*.json))

all: $(OUTS)

out/%.png: in/%.json
	~/go/bin/tuhirender -fit -o $@ < $<
```

Now you should have a directory with the subdirectories `in` and `out`
and a `Makefile`. Run `make` to convert all JSON images at once. When
new JSON files are added to `in/`, the `make` command will only run
incrementally for the new ones.
