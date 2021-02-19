# Datastore Emulator UI

I wrote in Go what can display data of Datastore Emulator.
We didn't introduce a big framework to make io easier for you to change.

I have no intention of using it for anything other than the emulator, but it may be possible if I devise something around authentication.

# Install

go install github.com/secondarykey/dem/_cmd/dem-view

or dem release.

# Run

```
$ dem-view
```

Boot on port 8088(default)
The list of endpoints is saved in $HOME/.dem.gob(default)

Read Help for details on the arguments.


