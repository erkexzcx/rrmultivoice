# rrmultivoice

This service grants ability for RockRobo vacuum robot to have different voice-lines each time it says something.

- [Get started](#get-started)
  * [Verifying compatibility](#verifying-compatibility)
    + [Root & SSH](#root---ssh)
    + [Process name](#process-name)
    + [Sounds location](#sounds-location)
    + [Available storage](#available-storage)
  * [Get audio files](#get-audio-files)
  * [Upload audio files](#upload-audio-files)
    + [Note about hard links](#note-about-hard-links)
    + [Delete existing files](#delete-existing-files)
    + [Upload multiple folders of sounds](#upload-multiple-folders-of-sounds)
    + [Populate with hard links](#populate-with-hard-links)
  * [Test if service works](#test-if-service-works)
  * [Start on boot](#start-on-boot)

# Get started

## Verifying compatibility

### Root & SSH

First thing first - your RockRobo must be rooted. You must be able to access it via ssh.

### Process name

There is a process on your robot vacuum that plays those sounds. Play the sound (e.g. "locate robot" action) and issue this command **while audio is playing**:

```
[root@rockrobo data]# lsof | grep -i wav
533     /opt/rockrobo/cleaner/bin/RoboController        /opt/rockrobo/resources/sounds/en/findme.wav
```

Great, you found it. It's called `RoboController`. If this is different for you - please raise an issue.

### Sounds location

From the previous step, you've found the location of the sounds (`/opt/rockrobo/resources/sounds/en/`) as well as process name.

### Available storage

Using below command you can check how much space you have in your robot vacuum:

```
[root@rockrobo data]# df -h /opt/rockrobo/resources/sounds/en/
Filesystem                Size      Used Available Use% Mounted on
rootfs                  493.9M    200.4M    268.0M  43% /
```

From this output you can see that the sounds are using `/` mountpoint, which has `493.9M` total storage and `268.0M` available. Great, because your uploaded files will need to be within the same partition as the sounds directory.

With workaround it's possible to use different partition, but more on that later.

## Get audio files

You can use my provided sound packs `additional_sounds.tar.gz` from this repo. Alternatively, your next best choice is to generate them yourself. See my other project https://github.com/erkexzcx/rrvoicegen.

## Upload audio files

### Note about hard links

Explanation of what is hard link:

_A hard link is a reference or pointer to a file's physical location on a disk, similar to a shortcut, but within the file system itself. Unlike a shortcut, a hard link is indistinguishable from the original file and has the same inode (the data structure in a Unix-style file system that describes a filesystem object such as a file or a directory). This means that even if you delete the original file, the hard link will still contain the file's data._

Hard link **limitation**:

_A hard link is a direct reference to the physical data on the disk and it must be on the same partition or disk as the original file. This is because hard links are essentially pointers to the same inode (data block) on the disk where the original file data is stored. If you attempt to create a hard link across different partitions or disks, it would fail because each partition or disk maintains its own set of inodes._

Which means you have 2 options here:

**Option 1**: Use as is (the rest of example is based on this method), considering empty space of `/` partition (where `/opt/rockrobo/resources/sounds/en/` is located).

**Option 2**: If you need more storage, maybe you can find another partition, such as mountpoint `/mnt/data`. In this case, you should delete everything from `en` folder, create (for example) `/mnt/data/rrmultivoice/en` dir and mount it to `/opt/rockrobo/resources/sounds/en` using `mount --bind` command as well as appropriate `/etc/fstab` entry (to mount on boot). Then you would work only with `/mnt/data/rrmultivoice/en`. **Note** that I have not tried this, so you are on your own.

### Delete existing files

SSH to the vacuum robot and delete (or backup somewhere) original audio files:
```bash
rm -rf /opt/rockrobo/resources/sounds/en/*
```

### Upload multiple folders of sounds

Now create a new directory in there:
```bash
mkdir /opt/rockrobo/resources/sounds/en/additional_sounds
```

Now upload your voice-lines packs (folders) to that folder. Here is the command that I used to upload all folders from `additional_sounds` on my PC to vacuum robot at `/opt/rockrobo/resources/sounds/en/additional_sounds/`:
```bash
scp -r -O additional_sounds/* 123.123.123.123:/opt/rockrobo/resources/sounds/en/additional_sounds/
```

Here is how I can verify if things are in place:
```bash
[root@rockrobo ]# cd /opt/rockrobo/resources/sounds/en
[root@rockrobo en]# ls -l | grep -v .wav
drwxr-xr-x    7 root     root          1024 Oct 15 17:47 additional_sounds
[root@rockrobo en]# ls -l additional_sounds | grep -v .wav
drwxr-xr-x    2 root     root          4096 Oct 15 17:43 custom1
drwxr-xr-x    2 root     root          5120 Oct 15 17:44 custom2
drwxr-xr-x    2 root     root          4096 Oct 15 17:44 custom3
drwxr-xr-x    2 root     root          5120 Oct 15 17:44 custom4
drwxr-xr-x    2 root     root          4096 Oct 15 17:44 custom5
```

### Populate with hard links

So the last missing thing is that the `/opt/rockrobo/resources/sounds/en/` directory contains no sounds. We need to populate it with initial sounds by using hard links. Note that soft links do not work (audio files are not being played).

Let's say one of the sound packs (folders) you have is this: `/opt/rockrobo/resources/sounds/en/additional_sounds/custom1`

In this case, run below command:
```bash
find /opt/rockrobo/resources/sounds/en/additional_sounds/custom1 -type f -exec ln {} /opt/rockrobo/resources/sounds/en/ \;
```

Now your `/opt/rockrobo/resources/sounds/en/` is populated with files that robot can play.

## Test if service works

SSH to your robot vacuum and inspect CPU architecture. In my case it's armv7:
```bash
[root@rockrobo ]# cat /proc/cpuinfo | grep Processor
Processor       : ARMv7 Processor rev 5 (v7l)
```

Since now you know the architecture, you can download an appropriate binary for your robot CPU from here: <insert_latest_release_link>

Recommended place to put binary file is `/mnt/data/rrmultivoice`.

Before you try, understand command line arguments and their default values (if not specified):
```bash
[root@rockrobo ]# /mnt/data/rrmultivoice -help
Usage of /mnt/data/rrmultivoice:
  -interval duration
        Interval for scanning fds. (default 300ms)
  -packsdir string
        Directory of additional sound directories (default "/opt/rockrobo/resources/sounds/en/additional_sounds")
  -soundsdir string
        Original directory from which robot plays voice-lines. (default "/opt/rockrobo/resources/sounds/en/")
```

Note about `-interval` - this service scans `RoboController` process for opened files. If one of those opened files is `/opt/rockrobo/resources/sounds/en/` - it updates hard link. So basically it's polling interval. The whole check takes about 1-2ms, so setting it to something very low (e.g. `20ms`) should be fine. However, I suggest keeping it default.

Now you can test the service to see if it works like this:
```bash
[root@rockrobo ]# /mnt/data/rrmultivoice
2023/10/15 20:15:13 Found RoboController PID: 533
2023/10/15 20:15:16 Replacing file: start.wav
2023/10/15 20:15:24 Replacing file: stop_clean.wav
2023/10/15 20:15:25 Replacing file: charging.wav
```

Try to spam some commands, such as start cleaning or go back to the dock multiple times. Each time it should say something different (depending on what voice-lines packs (dirs) you supplied it). It means everything is works and the remaining part - set it to start on boot.

## Start on boot

I am not certainly sure about this procedure, so I will cover **only what worked for me**. For better guidance and for different models - you should check out another repo that contains instructions on how to start binary on boot: https://github.com/porech/roborock-oucher

I have RoboRock S5 and for me, using `/etc/init` folder to put shell scripts works. I've created a file `/etc/init/S12rrmultivoice` file with the following contents:
```sh
#!/bin/sh

load() {
    curtime=`cat /proc/uptime | awk -F ' ' '{print $1}'`
    echo "[$curtime] start rrmultivoice"
    start-stop-daemon -S -b -q -m -p /var/run/rrmultivoice.pid -x /mnt/data/rrmultivoice
}

unload() {
    echo "Stopping rrmultivoice" >/dev/kmsg
    start-stop-daemon -K -q -p /var/run/rrmultivoice.pid
}

case "$1" in
    start)
        load
        ;;
    stop)
        unload
        ;;
    restart)
        unload
        load
        ;;
    *)
        echo "$0 <start/stop/restart>"
        ;;
esac
```

And lastly do a `chmod +x /etc/init/S12rrmultivoice`. Now reboot and test if voice-lines are changing.

Note that I am not quite sure how should I pass command-line arguments for this, but I never needed to.
