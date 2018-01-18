#!/bin/bash -x
#
# http://linux-notes.org/
#
#Vitaliy Natarov
#Set some colors for status OK, FAIL and titles
SETCOLOR_SUCCESS="echo -en \\033[1;32m"
SETCOLOR_FAILURE="echo -en \\033[1;31m"
SETCOLOR_NORMAL="echo -en \\033[0;39m"

SETCOLOR_TITLE="echo -en \\033[1;36m" #Fuscia
SETCOLOR_TITLE_GREEN="echo -en \\033[0;32m" #green
SETCOLOR_TITLE_PURPLE="echo -en \\033[0;35m" #purple 
SETCOLOR_NUMBERS="echo -en \\033[0;34m" #BLUE

KEY="eyJrIjoicjV3bklYVGFPRDZhd3NlS282VTBsdmhZcE1xcW4xNTAiLCJuIjoicmVhZGVyIiwiaWQiOjF9"
HOST="http://10.193.136.37"
DASH_DIR="./dashboards"

if [ ! -d "$DASH_DIR" ]; then
	 mkdir -p dashboards 
else
	 $SETCOLOR_TITLE_PURPLE
	 echo "|-------------------------------------------------------------------------------|";
	 echo "| A $DASH_DIR directory is already exist! |";
	 echo "|-------------------------------------------------------------------------------|";
	 $SETCOLOR_NORMAL
fi
  
#curl -sS -k -H "Authorization: Bearer $KEY" $HOST/api/search\?query\=\& |tr ']{' '\n'| cut -d ':' -f4| cut -d ',' -f1| cut -d '"' -f2 | grep -Ev "(^$|\[)"| cut -c 4-

for dash in $(curl -sS -k -H "Authorization: Bearer $KEY" $HOST/api/search\?query\=\& |tr ']{' '\n'| cut -d ':' -f4| cut -d ',' -f1| cut -d '"' -f2 | grep -Ev "(^$|\[)"| cut -c 4-); do
	curl -sS -k -H "Authorization: Bearer $KEY" $HOST/api/dashboards/db/$dash > dashboards/$dash.json
	$SETCOLOR_TITLE_GREEN
	echo "The [ $dash ] dashboard has been exported"  
	$SETCOLOR_NORMAL
done

$SETCOLOR_TITLE
echo "|-------------------------------------------------------------------------------|";
echo "|----------------------------------FINISHED-------------------------------------|";
echo "|-------------------------------------------------------------------------------|";
$SETCOLOR_NORMAL
