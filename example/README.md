The simplest bood example
=========================

From this directory, try the following commands:

#### Install bood
```
$ go get -u github.com/roman-mazur/bood/cmd/bood
```

#### Build the program
```
$ bood
INFO 2021/03/22 19:41:56 Ninja build file is generated at out/build.ninja
INFO 2021/03/22 19:41:56 Starting the build now
[4/4] Archive binary myBood to zip task2
```

#### Run the program
```
$ out/bin/hello
Hello, World!
```

#### Run build again (and see nothing is done)
```
$ bood
INFO 2021/03/22 19:43:27 Ninja build file is generated at out/build.ninja
INFO 2021/03/22 19:43:27 Starting the build now
[1/1] Archive binary myBood to zip task2
Archive is current
```
