# people's cloud
using the ideas of botnets for good

# aim & intent
to create a technological solution so that people with computing power can share it with others.  

in the current iteration, this is focused on botmakers so that botmakers *with* servers can run swarms where botmakers *without* them can run their creations, while (hopefully!) reducing the overhead needed for the server people to share their resources.

# terminology
* it's coming!
* nodes
* swarms
* programs

# usage
more examples coming! but the help shows a lot of usecases, just make sure you have fulfilled all the requirements and that you're running a pc daemon in another window
`go run pc --help`
`go run pc daemon config.toml`

# alpha version, work in progress etc

probably has lots of bugs, but i'm hunting them with a mace and hammer (yeah both) so don't worry

documentation forthcoming!

# what you can currently do
* run self-contained nodejs, python2 and python3 programs on other people's computers

## python 2/3
* make sure your main file that does all the work is called bot.py
* have a file called requirements.txt in the same folder as bot.py, which lists all of the modules you've downloaded to make your program or bot run

## nodejs
* make sure you have a package.json file that lists the start script (example of this forth-coming - don't worry!)

# coming up
* binaries & github releases so you don't have to bother with using go run to run it all!
* password protected channels
* investigations into running everything on raspberry pis (i.e. cross-compiling binaries for the pi)
* investigations into exposing ports so potentially webservices could run easily
* investigations into supporting other languages, and potentially most languages!
* lots and lots of bugfixes :bug:

# requirements
* ipfs
* docker
* go

# contributions
would love some! i'll think about and write some instructions on what kind of format they should follow to make life easier for intergrating changes and new features.

# issues/bugs/crashes
let me know! either write a nice github issue about it, or ping me on [twitter](https://twitter.com/cblgh)

# credit
organization branding by [osavox](https://twitter.com/osavox)
