[![MIT Licence](https://badges.frapsoft.com/os/mit/mit.svg?v=103)](https://opensource.org/licenses/mit-license.php) 

# **ZFS-backup-tool**
a tool to manage zfs snapshots easily


----------
# **Command**
    zbackup ( target_dataset [rotation_count] | --list [target_dataset [ID]] | --delete [target_dataset [ID]] | --daemon [--config path/to/your/conf])

#Create
`zbackup target_dataset [rotation_count]`

create a snapshot for target_dataset with rotation_count
rotation count: The sum of snapshots of this snapshot,it will remove old snapshot automatically (default=20)

![enter image description here](http://i.imgur.com/1uxK5pk.png)
![enter image description here](http://i.imgur.com/fZWn3PQ.png)

#List
`zbackup --list [target_dataset [ID]]`

list all snapshots created by this tool
you can use target_dataset and id to modify your list

![enter image description here](http://i.imgur.com/gZFsO6Q.png)

#Delete
`zbackup --delete [target_dataset [ID]]`

delete snapshots

![enter image description here](http://i.imgur.com/Mxx4CyX.png)

#Help
`zbackup --help`

![enter image description here](http://i.imgur.com/4DEnHWR.png)

#Daemon
`zbackup --daemon [--config path/to/your/conf]`

Backup automatically background.
You can use custom conf path with `--config path/to/your/conf`,or zbackup will use /usr/local/etc/zbackup.conf in default.

Here is a conf file sample.
```bash
[zroot/home]		# Dataset to be snapshot
enabled=no			# (default:yes)
policy=4x15m		# Snapshot policy
					# e.g. AxB
					# snapshots are created with interval B, and there can be A snapshots (created by zbackup) at most
					# supported units of interval are:
					# m: minutes, h: hours, d: days, w: weeks
[zroot/data]
policy=5x1d

[zroot/etc]
policy=12x1d

[zroot/home]
policy=4x1w
```


----------
#**Be careful**

 - You should execute "create" and "delete" with sudo.
 - You should put zbackup under $PATH, or it may corrupt when executing some features.

#**Install**
you can download the latest release [here](https://github.com/FrozenKP/ZFS-backup-tool/releases).
