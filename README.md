# go-twitter-test

Just a code test I did with Go to create a Twitter like sample app.

# API endpoints

## POST /v1/messages

Used to create tagged messages.

**HTTP Request:**
* `X-User-ID`: added by an API Gateway based upon an `Authorization` header
* `Content-Type`: `application/json`
* Body: `{"text":"A very meaningful message","tag":"philotimo"}`

**HTTP Response:**
* Status codes
  * `201` with Location header on success pointing to the newly created resource
  * `400` if the body is malformed or invalid
  * `401`|`403` for auth errors
  * `500` when there's a backend error (e.g. can't connect to the DB)
* `Content-Type`: doesn't matter
  * the client will just have the status code to understand what is going on for now
  * devs can check the logs for more details about why a backend error has happened
  * later better error messages can be returned as JSONs if needed to give clients more context


## GET /v1/messages

Used to get all messages or a filtered set of messages (getting a single message is not supported for now
e.g. `GET /messages/<id>`).

Supported filters:
* filter by tag
* count by date range

**HTTP Request:**
* `X-User-ID`: added by an API Gateway based upon an `Authorization` header
* `tag` query parameter (to filter by tag)
* `dateStart` and `dateEnd` query parameters
* `count` query parameter (1|0) to instruct the API to return a count instead of a list of messages
  * can be used along with `tag` and `dateStart`, `dateEnd`

**HTTP Response:**
* Status codes
  * `200` OK
  * `400` if one or more the query parameters are invalid
  * `403` if no user ID is specified in the `X-User-ID` header
  * `500` when there's a backend error (e.g. can't connect to the DB)
* `Content-Type`: `application/json`
  * example: `[{"text":"A very meaningful message","tag":"philotimo"}]` or `123` if `count` is `1`
  * alternative: we could potentially have another endpoint (e.g. `GET /messages/count`) just for the count

**Note:** pagination won't be handled here, there's plenty of literature about how to handle pagination properly.
One approach I like is to do tokenized pagination to avoid inconsistencies when browsing back and forth through 
sorted lists.

# Authentication and Authorization

Having an API where you handle both writes and reads makes for a good monolith and whereas I'm not a fan of
monoliths I think they make sense for this exercise (for the sake of simplicity).

Given that I'm stuck with a monolith, I won't be using a discoverable data set like I'd probably do, so for now,
I don't think it makes sense to expose some options (i.e. *count by date range*) via another "management" API but I 
rather put all endpoints in the same API (this one).

Such API would then embed the responsibilities of the mentioned data set (by exposing the underlying data that
lives in the database) along with the responsibilities you would normally have in this kind of APIs
(message creation).

Considering all of the above I propose an API Gateway (which I won't implement here) to deal with resources access.
The API Gateway would sit in front of this API to handle both Authentication and Authorization 
(i.e. *count by date range* for admins only).

# Scalability

As I mentioned in the previous paragraph the only pieces involved in this simplified architecture are the following:

* API Gateway
* Messages API
* HA Cluster (could be SQL or NoSQL)

In front of the API Gateway there should be a clustered pair of (or more) Load Balancers to ensure availability
and fair load distribution (i.e. not Round Robin). More Load Balancers could sit in between the API Gateways and
the Messages APIs as well. The Messages APIs would then connect to a HA Cluster of databases to consume data.

Generally speaking an Active-Active pattern with Sharding could ensure HA and a fair load distribution of both reads and
writes even with multi-clusters deployed across regions. The idea is to make sure both reads and writes are local so
routing during writes is essential (given a key it'll always be routed to the same cluster where one node can accept
the write). The data is then replicated across clusters so that reads are as much as possible local.
If a primary node (the one that accept writes) goes down in a cluster, one of the secondary nodes is immediately elected
as primary.

Although this may differ depending on what DB technology you may choose, what I've been describing is more or less
what you could get with a MongoDB multi-cluster Active-Active configuration with Sharding enabled.

# Considerations

Once the Messages API grows in complexity there could be room for more Microservices and potentially Event Sourcing
and CQRS so that we could separate the reads from the writes to better deal with scalability and concurrent writes.

**Note:** Event Sourcing and CQRS are not a dogma everyone should just blindly follow though, it must make sense for the 
kind of application you're building and here we don't have enough complexity and information to assert whether it would 
be a good fit or not. Given the amount of information we have we can only speculate so what you'll see next is an
hypothetical event driven architecture that could potentially work with the kind of scenarios we might see.

`Event Store` -> `Reducers` -> `Discoverable Data Set` -> `API` -> `Event Dispatcher API` -> `Event Store` -> ∞

1. `Event Store`: as the name implies a store for all the raw events
2. `Reducers`: small Microservices subscribed to the Event Store that would create snapshots (event folding)
   to persist them into a `Discoverable Data Set`
3. `Discoverable Data Set`: it could be a single node database or a HA cluster, it would be discoverable and
   expose the data from one or multiple domains to whoever wants to consume it. Only consumption is allowed so
   the reducers are the only part of the architecture that can update the data set. A data catalogue can be used
   to define boundaries between data domains if needed
4. `API` or other `clients`: at this layer we have APIs or simply clients (e.g. cronjob) that can consume the data 
   from the `Discoverable Data Set` (by querying/aggregating). If needed they can expose such data to a Frontend 
   client and/or create more data by dispatching an event to the `Event Dispatcher API` and then subscribing for 
   an expected change to the `Discoverable Data Set` (Publish/Subscribe)
5. `Event Dispatcher API`: its only responsibility is to allow other actors in the architecture to create events which
   then end up in the `Event Store` thus the loop would start again from point 1 (the events are picked up by `Reducers` 
   which persist a folded state into the `Discoverable Data Set` and so on)
