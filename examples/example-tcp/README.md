`example-tcp.go` reads live NMEA data and plots the ships on a google map, served at <http://localhost:8080>.
You can click on the ships to get info about them.
It uses <ais1.shipraiser.net> (hardcoded), which unfortunately isn't active anymore.

If you have another source for AIS data, set it inside `example-tcp.go` and run it:

    go run example-tcp.go
