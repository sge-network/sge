# **Overview**

The Bet module is responsible for receiving and processing requests to place and settle bets. In the case of placement, it validates the request and places the bet, and in the case of settlement, it checks the result of the sport-event, determines the bet result, and settles the bet using `OrderBook` module.
