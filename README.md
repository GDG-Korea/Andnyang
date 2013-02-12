Andnyang
========

An IRC Bot that written in [GOLANG](http://golang.org/) for [GDG Korea](https://developers.google.com/groups/directory/south-korea) (Google Developers Group Korea) channels 
 * Server: [irc://irc.ozinger.org](irc://irc.ozinger.org)
 * [GDG Korea Android](https://plus.google.com/communities/100903743067544956282)'s channel: #gdgand
 * [GDG Korea Women](https://plus.google.com/communities/116463742909053357630)'s channel: #gdgwomen

Setup
-------
Install:

    sudo apt-get install postgresql
    go get github.com/dalinaum/Andnyang
    sh <your GOPATH>/src/dalinaum/Andnyang/genDB.sh

You can see your `GOPATH` by `export | grep GOPATH`.

If it does not exist, you should add your `GOPATH` to your `.bashrc` as follows. 

    export GOPATH=~/mygo
    export PATH=$GOPATH/bin:$PATH

After setting your `.bashrc`, enter `source ~/.bashrc` command into your terminal to use modified setting.

Run:

    Andnyang &

You can use `hohup` command to hide standard output.

    nohup Andnyang &

Authors
-------
 * Leonardo YongUk kIm dalinaum@gmail.com
 * Homin Lee homin.lee@suapapa.net

