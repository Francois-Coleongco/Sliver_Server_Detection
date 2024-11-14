# Sliver_Server_Detection

Detecting a SliverC2 Server Running on an Ubuntu Machine

In order to do this, I must be able to get information about processes connecting to the internet

once i know this, i can find  the process and see what syscalls it is making

i can see if it uses a terminal based on this info, and if it uses event polling then it's seriously suspicious.



do this in golang and the os lib (probably the os lib)

use the `ss` command programmatically in golang to display the active connections

get each connection, and find the PIDs through `lsof`

using the PIDs, do `strace` on them to see the syscalls being called

parse those syscalls to detect rev shell related activity.

terminal activity, epolling (event polling), really really suspicious...

Another thing you can check is file access. i think in strace u can see the uh whats it. i think you can see some stuff when u ls'd i forgot when exactly though.

------------------------------------------------------------------------------------------------

Upon testing on my own virtual machines, a Sliver Server running on Kali Linux and a victim Ubuntu Desktop, I found that the syscalls made by the executable were the following: futex, epoll_ctl, fcntl, and of course, read and write.

Note: the command used to generate the sliver payload was `generate --mtls <ip_of_kali_machine> --os linux`

futex for fast non-blocking synchronization of threads
epoll_ctl to wait for server commands
fcntl to manage file descriptor operations

read and write i have seen receive and send tls-looking data starting with:

read(3, "\27\3\3\\" 

AND

write(3, "\27\3\3\\"

because they go to the same fd and the data appears to be tls, this is likely the sliver server's commands encrypted in transit.

Using the `ps -aux` command, I found that the process uses a pts (pseudo terminal slave) which means on top of encrypting it's read and writes, it also uses the terminal. This is... extremely suspicious. So the idea now is when I detect a process using a pts and encrypting it's traffic with tls, i will kill it.


further checks:

`inotify` to listen to file reads and writes, as well as creations and modifies.

`file hashes` use sha256sum on that sucker and use programmatically get a result from virustotal using their api

---------------------------------------------------------------------------------------------------

since the SliverC2 framework does not have linux stager compatibility, the only thing i think i need to worry about is statistically analyzing the single executable binary after killing the process.



---------------------------------------------------------------------------------------------------


