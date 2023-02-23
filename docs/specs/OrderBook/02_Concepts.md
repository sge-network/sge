# **Concepts**

The Order Book is tasked with maintaining the order book, participations, exposures and order book settlement,
each order book will be created as a one-to-one dependency of sport-event. in action, sport-event module calls
book initiation method of the order book module to create the corresponsing order book, paarticipation and exposures.
The created book for the sport-event is inititated will be maintained until the sport-event marked as settled.

The Book for a given sport-event, is made up of Book Participants. The first book participant for any book of
a sport-event will be Strategic Reserve. At the creation of a sport-event, in order to faciliate betting on
the created sport-event. The order book initiates a book for the sport-event and sets the first participation for the
strategic reserve module.

Once the order book has initiated a book for a sport-event, users can either bet against the house or
become a part of the house by depositiion of chosen amount through the House module. When a user deposits chosen
amount through the house module, the house module will call the order book module to update the order book
and set the participation for the user on the requested sport-event.

The deposit amount of book participants is used to facilitate betting on the sport-event.

The payout that need to be paid by the system is named Exposure, there are two types of bet odds exposures:

- The odds exposure are the payouts should expected to be paid.
- The participation exposure are the payout should go to the participant.
