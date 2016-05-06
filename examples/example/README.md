`example.go` is a sample code on how to use the AIS library.

For a quick demonstration, try to run:

    cat nmea-sample.txt | go run example.go

Or with live data:

    netcat ais1.shipraiser.net 6492 | go run example.go

Or with a sentence:

    echo '!AIVDM,1,1,,A,14eGrSPP00ncMJTO5C6aBwvP2D0?,0*7A' | go run example.go
