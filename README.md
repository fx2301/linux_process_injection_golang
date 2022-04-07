# Overview

Linux project injection demo. Creates an memory-only file descriptor to write executable to, and then executes it using a fake process name of "injectedprocess". 

# Example

```
$ go build memfd_launch.go 
$ ./memfd_launch --inject
Injecting into process...
File descriptor is 3
Bytes written to memfd 2030005 at path /proc/self/fd/3
$ ps $(pidof injectedprocess)
    PID TTY      STAT   TIME COMMAND
 173277 ?        Ssl    0:00 injectedprocess
```

# References

* This work is derived from https://www.guitmz.com/linux-elf-runtime-crypter/ / https://github.com/guitmz/ezuri and should work on multiple architectures.
* https://medium.com/confluera-engineering/reflective-code-loading-in-linux-a-new-defense-evasion-technique-in-mitre-att-ck-v10-da7da34ed301
* https://github.com/golang/go/issues/227
