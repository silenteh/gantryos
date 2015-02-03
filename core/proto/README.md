PROTOCOL

This modules contains:
- protocol definition
- messaging format


Proto simple description: all messages are idempotent !
1) Master starts
2) Slave starts
3) Slave register (syn -> ack)
4) Slave starts to send offers
5) Master accepts tasks from API and send them to the slave
6) Slave starts the task and keeps the master updated


-- Error cases
1) Slave starts sooner than the master
2) Master dies but slave alive
3) Slave die, but master alive
4) Slave alive but cannot communicate with master and viceversa
5) Not enough resources on the slave, abort completely ?
6) Persistence storage fails, then what ? We need a replicated persistence storage



-- Security
1) communication between master and slaves should happen in TLS 1.2 with PFS




