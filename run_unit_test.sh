#!/bin/bash

# Execute all unit tests.

TRUE=1
FALSE=0
EXCEPT=$FALSE


#----- Run the unit tests
go test -v || EXCEPT=$TRUE

cd perfTestUtils
go test -v || EXCEPT=$TRUE
cd -

cd testStrategies
go test -v || EXCEPT=$TRUE
cd -

#----- Report overall success
MSG="All Unit Tests PASSED"
if [ $EXCEPT -eq $TRUE ]; then
	MSG="!! Unit Test FAILED !!   Search '--- FAIL' in output above for details."
fi

echo
echo
echo $MSG
echo

exit 0