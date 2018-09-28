# SPFLoop

## What is this?

**SPF** records are a rudimentary form of spam protection for email. When a mail server receives a message, it is supposed to look up TXT DNS records to determine if the message is coming from an ip address authorized by the domain owner.

These records can include other records, which cause additional dns lookups. Those can include other records and cause further lookups.

SPF resolvers are supposed to stop at **10** dns lookups. If they do not, they can be used as ddos soures on unsuspecting dns servers. This is an attempt to see how respectful common main providers are of that limit.

## What does it do?

This is an extremely stupid dns server. It will respond to TXT queries for *any* domain with responses like so:

- `example.com` -> `v=spf1 include:_spf1.example.com -all`
- `_spf1.example.com` -> `v=spf1 include:_spf2.example.com -all`
- ...
- `_spf300.example.com` -> `v=spf1 include:_spf301.example.com -all`

It will log queries and see how deep they go. 

You can send a message with an arbitrary subdomain, and watch what queries get sent from the spf resolver.