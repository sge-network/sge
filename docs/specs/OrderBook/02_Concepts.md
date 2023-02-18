# **Concepts**

The Order Book is tasked with maintaining the book for sport events. At creation of a sport event, a book for the created sport event is inititated which is maintained throughout till sport event resolution and th the event of sport event resolution the book is settled.

The Book for a given sport event, is made up of Book Participants. The first book participant for any book of a sport event will be Strategic Reserve. At the creation of a sport event, in order to faciliate betting on the created sport event. The Strategic Reserve initiates a book for the sport event and becomes the first book participant by traferring an amount for the Book.

Once the Strategic Reserve has initiated a book for a sport event, users can either bet against the house or become a part of the house by depositing some amount through the House module. When a user deposits some amount through the house module, the house module will call the order book module to update the order book and add the user as a book participant for the requested sport event.

The deposit amount of book participants is used to facilitate betting on the sport event.
