#/bin/bash
cd $(dirname $0)

# Uncomment this line to automatically watch the server log.
#AUTOTAIL=TRUE

LOG=log

if [[ $1 == "log" ]]
then
	AUTOTAIL=TRUE
fi

if [[ $(pgrep distru | wc -l) != 0 ]]
then
	killall distru
fi

go fmt
if [[ $? != 0 ]]
then
	exit
fi

go build
if [[ $? != 0 ]]
then
	exit
fi

cat /dev/null > $LOG
./distru &> $LOG &

if [[ $AUTOTAIL == TRUE ]]
then
	tail -f $LOG
fi

exit 0
