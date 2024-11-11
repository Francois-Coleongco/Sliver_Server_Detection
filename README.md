#Detecting a SliverC2 Server Running on an Ubuntu Machine

In order to do this, I must be able to get information about processes connecting to the internet

once i know this, i can find  the process and see what syscalls it is making

i can see if it uses a terminal based on this info, and if it uses event polling then it's seriously suspicious.



do this in golang and the os lib (probably the os lib)

use the `ss` command programmatically in golang to display the active connections

get each connection, and find the PIDs through `lsof`

using the PIDs, do `strace` on them to see the syscalls being called

parse those syscalls to detect rev shell related activity.

terminal activity, epolling (event polling), really really suspicious...

# Sliver_Server_Detection
