## TODO

- Most frequent messages are types 1, 3, 4, 5, 18. We should decode all of them. Currently we decode 1, 3, 4 and 5.
- Find out what ports' reports are. I think they announce when a ship arrives/leaves from port.
- Found out how satellite can be used to capture all AIS.
- Decode some Type 8 messages: 1-11, 1-31 (meteorological), Area Notice, Extended Static and Voyage Related Data, Route Information, etc
- Decode Type 9 messages: search and rescue aircrafts position report
- Decode Type 12 messages: addressed safety related message
- Decode Type 14 messages: broadcast safety related message
- Decode Type 19 messages: Class B extended position report (something between type 18 and type 5)
- Decode Type 21 messages: aid to navigation report
- Decode Type 27 messages: long range Class A position reports
- Decode Type 6 DAC-FIDs in order to have some stats. Though Type 6 mesages are very rare it seems.
- DONE: Decode MMSI (includes country or other info as well)
- DONE: Implement a generic router where we feed it messages and it returns message type and payload or error.
- DONE: implement a parser for payloads spanning across 2 or more AIS messages. (tricky: variable length, out of order maybe, maybe we should expire if we don't receive all parts after some time)


## Notes:

- Type 1/2/3/4/9/11/18/22 -        168bits -       1 sentence
- Type 5                  -        424bits -       2 sentences
- Type 6/8/12/14          - up to 1008bits - up to 5 sentences
- Type 7/13               -     72-168bits -       1 sentence
- Type 10                 -         72bits -       1 sentence
- Type 15                 -     88-160bits -       1 sentence
- Type 16                 -  96 or 144bits -       1 sentence
- Type 17                 -  80 to 816bits - up to 4 sentences
- Type 19                 -        312bits -       2 sentences
- Type 20                 -     72-160bits -       1 sentence
- Type 21                 -    272-360bits -       2 sentences
- Type 23                 -        160bits -       1 sentence
- Type 24                 -    328-336bits -       2 sentences
- Type 25                 -  up to 168bits -       1 sentence
- Type 26                 -    60-1064bits - up to 5 sentences
- Type 27                 -  96 or 168bits -       1 sentence

## Stats

### Out of 76175 type 1, 3, 4 messages analyzed:

71162 / 76175 Ship
 4344 / 76175 Coastal Station
  316 / 76175 Group of ships
  213 / 76175 Aids to navigation
  109 / 76175 SAR —Search and Rescue Aircraft
   20 / 76175 Diver's radio
   01 / 76175 invalid MMSI

### Out of 701088 messages analyzed:

Type  1: 218359
Type  3:  55917
Type  4:  23239
Type  5:  65672
Type  8: 272742
Type  9:    226
Type 12:      3
Type 13:      5
Type 18:  25117
Type 21:  17380
Type 24:  22428
