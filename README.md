[![MIT Licence](https://badges.frapsoft.com/os/mit/mit.svg?v=103)](https://opensource.org/licenses/mit-license.php) 

# ZFS-backup-tool
a tool to manage zfs snapshots easily


----------

**Command**
-----------

    zbackup ( target_dataset [rotation_count] | --list [target_dataset [ID]] | --delete [target_dataset [ID]] )


----------

**Create**  `zbackup target_dataset [rotation_count]`

> create a snapshot for target_dataset with rotation_count
> rotation count: The sum of snapshots of this snapshot,it will remove old snapshot automatically (default=20)

![enter image description here](http://i.imgur.com/1uxK5pk.png)
![enter image description here](http://i.imgur.com/fZWn3PQ.png)

**List** `zbackup --list [target_dataset [ID]]`

> list all snapshots created by this tool
> you can use target_dataset and id to modify your list

![enter image description here](http://i.imgur.com/gZFsO6Q.png)

**Delete** `zbackup --delete [target_dataset [ID]]`

>delete snapshots

![enter image description here](http://i.imgur.com/Mxx4CyX.png)

**Help** `zbackup --help`

![enter image description here](http://i.imgur.com/mXRtkDR.png)


----------
**Be careful**
--------------
you should execute "create" and "delete" with sudo.

**Install**
--------------
you can download the latest release [here](https://github.com/FrozenKP/ZFS-backup-tool/releases).
