#/bin/bash
cd $(dirname $0)

# Uncomment this line to automatically watch the server log.
#AUTOTAIL=TRUE

LOG=log

if [[ $(pgrep distru | wc -l) != 0 ]]
then
	killall distru
fi

go build
cat /dev/null > $LOG
./distru &> $LOG &

if [[ $AUTOTAIL == TRUE ]]
then
	tail -f $LOG
fi

exit 0
