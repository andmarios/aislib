# aislib #

A Go library to read AIS sentences and decode AIS messages.

The most simple example of decoding and printing an AIS message:

```
$ cd examples/example
$ echo '!AIVDM,1,1,,A,14eGrSPP00ncMJTO5C6aBwvP2D0?,0*7A' | go run example.go
=== Class A Position Report (1) ===
 Repeat       : 0
 MMSI         : 316013198 [Ship, Canada]
 Nav.Status   : Under way using engine
 Turn (ROT)   : no turn information
 Speed (SOG)  : 0.0 knots
 Accuracy     : High accuracy (<10m)
 Coordinates  : 130°18.9742'W  54°19.2666N
 Course (COG) : 237.9°
 Heading (HDG): not available
 Manuever ind.: not available
 RAIM         : in use
```

# History

While I was working for _marine.travel_ —a social site for sea lovers— I wanted to polish my Go
skills. I thought a nice project would be an AIS parser, which later we could use to implement
new functionality for our product.

I wrote _aislib_ outside my work duties and time, but intended to use it there. Unfortunately
_marine.travel_ didn't make it. Thankfully we are on new, equally exciting adventures now. Since
_aislib_ didn't have any prospects anymore, I asked for permission to release my code under a
permissive license. It was granted and here is the repository. :)

# Quality and Use Disclaimer

Please bare in mind that this was my first big-ish project in golang, some poor choices are to
be expected. I am pretty certain that the main functionality is robust though, so if it doesn't
fit your needs, you could use the decoding parts and implement your ideal management parts.

# Features and State

Inside the `examples/` directory you will find two examples to get you going. Start with the
simple one, `example.go`. To run:

     $ cat nmea-sample.txt | go run example.go

**aislib** can decode type 1, 2, 3 (Class A Position Report), 4 (Base Station Report),
5 (Static Voyage Data), 18 (Class B Position Report) messages. It may also understand type 8
(binary Broadcast) messages, report their respective type and extract the binary payload.

These are the most common types you will find. If you are interested in extending aislib, it is
worth implementing type 21 and 24 decoding.

A limitation is multi-sentence messages. Messages that span across AIS sentences will only be
decoded if (a) they come in order and (b) do not interleave with other multi-sentence messages.

In my experience from capturing AIS data from ais1.shipraiser.net for many hours, the messages
you receive almost always meet these criteria. I think I never saw a sentence parsing failing
due to this restriction.

# How it Works

As stated, some poor choices may have been made.

Each AIS message type has different fields. Thus I implement a router, that receives sentences,
checks their type, their checksum and if the message they carry spans across many sentences
(and it reconstructs the message's payload) and pass back the type and the payload. This work
is done through channels, an incoming channel to the router, to pass AIS sentences and an
outgoing channel to receive message payloads (and types). There is also a second outgoing
channel where you receive the sentences that the router failed to recognize for any reason
(e.g bad checksum or out of order multi-span message). It is useful for debugging.

So in sort you send AIS sentences into the router and get tuples with AIS message type and
payload.

You should switch on the message type to the proper decoding function.

In retrospect, now I know that I could return just the message and let you use type assertion.
This wouldn't change much though, as again you would need a select statement to deal with the
result.

Check `example.go` to understand how the router and decoding function works.

# License

Check `LICENSE` file. In sort it is GPL version 3 or greater.
