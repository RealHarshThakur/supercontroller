# Supercontroller
 Supercontroller is an example controller pattern that can connect to multiple Kubernetes clusters, focus on events of entire API group (example.com) and react to them.

It is-
* Multi-cluster aware
* Multi-resource: will watch all resources in an API group
* Dynamic : doesn't need static schema of resources
* Multi-modal: multiple handlers can be hooked (not sure if this is a good idea)


## Use case
Useful for common operations that need to be performed across multiple clusters for an API-group. For example-
* Create a global view in a single datastore
* Process CR for quota, accounting, etc per tenant 

