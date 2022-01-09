# Gio demo/test

See [gioui.org](https://gioui.org).

./TestWindows/testwindows.go is an application that tests the new window options.
You can maximize, minimize, fullscreen, center, size, position windows either
from the command line or by clicking on buttons.

Example: 
```
>go run testwindows.go -test=n

where n is one of the following:
1 = Maximized window
2 = Minimized window
3 = Centered window
4 = Positioned window
5 = Sized window
6 = Fullscreen window
```

# License

Dual MIT/Unlicense; same as Gio
