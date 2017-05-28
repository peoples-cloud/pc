# people's cloud
using the ideas of botnets for good  

_the documentation on this page is always under construction :warning:_ 
### aim & intent
to create a technological solution so that people with computing power can share it with others.  

in its current iteration, this is focused on botmakers so that botmakers *with* servers can run swarms where botmakers *without* servers, or the monetary means to acquire & run them, can run their creations, while (hopefully!) reducing the overhead needed for the server people to share their resources.

### what is it
a tool for creating decentralized swarms of computers, where nodes in the swarms can deploy programs to be run on other nodes 

### who is it for
obviously everyone since it is open source, but  
it's intended for small communities with an already established trust, as technical trust is a complex subject  
it's for [art bot / folk bot](https://youtu.be/87yiUjGnXdI) people that want to run their bots with the help of their friends  
it's for private communities that want to experiment with weird stuff    
basically, if you don't trust the people you will join the swarm with - **don't run this**

### how do i use it
first make sure you have fulfilled all the requirements (ipfs/docker/go)

then, in a terminal window, do:  
`go run pc.go daemon config.toml`

using another terminal window you can now issue commands to the daemon:
```sh
go run pc.go --help
Usage:
  pc [command]

Available Commands:
  create      Create a new swarm
  daemon      Starts the pc daemon
  deploy      Deploy a program to be run by a node in the swarm
  join        Join a new swarm
  leave       Leave a swarm
  list        Lists all connected swarms and deployed programs
  stop        Stop a deployed program
```
_binaries are coming. they can also be built using `go build pc.go`_

### terminology
* nodes
  * simply a computer running the `pc daemon` command
* swarms
  * a group of connected nodes, sharing the load of the deployed programs within the swarm
* programs
  * self-contained javascript/python code, e.g. a python twitter bot

### what you can currently do
* run self-contained nodejs, python2 and python3 programs on other people's computers

#### python 2/3
* make sure your main file that does all the work is called bot.py
* have a file called requirements.txt in the same folder as bot.py, which lists all of the modules you've downloaded to make your program or bot run

#### nodejs
* make sure you have a package.json file that lists the start script (example of this forth-coming - don't worry!)

### requirements
* ipfs
* docker
* go

## coming up
* binaries & github releases so you don't have to bother with using go run to run it 
* standalone clients, which would included
  * binaries with ipfs bundled
  * binaries without the docker requirement (still mulling this about)
* support for password protected swarms
* configuration examples
* modularizing the code even further, for reuse in other projects
* investigations into... 
  * a dockerized setup so that you can just download an image and launch that instead
  * running everything on raspberry pis (i.e. cross-compiling binaries for the pi)
    * creating a build process to automate cross-compiling, alongside other binaries
  * exposing ports so potentially webservices could run & communicate with the outside world
  * supporting languages other than python & js
  * granularity in configuration of your node, allowing you to restrict runtime, RAM usage and the like
  * technical audits that will allow for traceability of deployed programs (and thus some form of security for hosts)
  * potentially using pubsub as a communication mechanism
* exposing under-the-hood details via options for communities that want to tinker
* bugfixes :bug:


## contributions
would love some! i'll think about and write some instructions on what kind of format they should follow to make life easier for integrating changes and new features.

## issues/bugs/crashes
#### alpha version, work in progress etc
there are probably lots of bugs lurking around, so let me know if you find any! 

either write a nice github issue about it, or ping me on [twitter](https://twitter.com/cblgh)

### credit
organization branding by [osavox](https://twitter.com/osavox)
