MODULES:
- core
  + master
  + slaves
  + resources offer
  + resource declined
  + protocol definition (proto buf)
  + comunication via message passing
  + reconcilation
  + master election
  + master request forwarding in case the requests do not come to the current master
  + high availability
  + roles definition and resource allocation
  
- REST
  + API
  + UI

- SECURITY
  + resource isolation (cgroups)
  + namespaces, non-root, grsecurity, signatures, TLS, authenticated encryption
  + monitoring (resources, networking, user abuse)

- TASKS
  + Docker
  + LXC
  + XEN (High availability will be problematic because of the VM size)
  + Long tasks, short tasks (pure shell commands)

- REPORTING
  + logging (vms, containers logs)
  + reporting (resources utilization: cadvisor ?) 
  + monitoring (resources, networking, user abuse)
  + healthchecks
  + alerting in case of tasks or masters or slaves or networking problems or scarse resources
  + event bus: all events are exposed to subscribers

- ORCHESTRATION
  + networking (openflow - distributed switch - customizable network scenario)
  + service discovery: zookeeper, etcd, consul
  + proxying
  + URL rewriting
  + container and VMs central repositories (S3, local disk)

- Supported cloud providers (roadmap):
  + AWS
  + Rackspace


