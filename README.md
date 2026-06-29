# aethoria_marketplace
A mystical bazaar where seekers buy and bid on legendary artifacts of Aethoria.


## Architecture Decisions (ADR)

### Why PostgreSQL?

The main features of this project, like placing bids and buying items, need database transactions. PostgreSQL supports ACID transactions and row-level locking (`SELECT FOR UPDATE`), which makes these operations safe. If I used Redis or another database, I would have to implement the locking logic myself.

### Why sqlc?

I chose `sqlc` because it lets me write normal SQL queries while still getting type-safe Go code. It is lightweight, has no ORM overhead, and the generated code is fast because it does not use reflection.

### Why DDD layering?

The project has several business rules, such as bid increment rules, legendary item restrictions, and wallet reserve/release logic. I put these rules inside the domain layer so they are always checked, no matter which part of the application calls them.

### What I left out

* No authentication. (since we don't have a frontend) The `guild_id` is sent in the request body.
* No pagination for list endpoints.
* The scheduler is implemented with a simple ticker instead of a job queue.

### What I would add if I had more time

* More comments on the code to be well-documented
* More tests (both unit and integration tests)
* An outbox pattern for wallet transactions to make sure every transaction is recorded.
