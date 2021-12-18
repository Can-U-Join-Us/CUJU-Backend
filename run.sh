#!/bin/bash

cmd=$1

case "$cmd" in
	-serve)
	go run main.go
	;;

	*)
	echo "'$cmd' is unknown command"
	;;
esac
