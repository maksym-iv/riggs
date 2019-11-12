# Answers
* Prove how it aligns to 12factor app best practices
  1. Codebase  
      True. In git.
  2. Dependencies  
      True. `go.mod` for deps.
  3. Config  
      True. Configurable via env. vars. and cli.
  4. Backing services  
      True. GeoDB datastore is attached resource.
  5. Build, release, run  
      Strictly separate build and run stages
  6. Processes  
      True. Executed in isolated env.
  7. Port binding   
      True. Docker + K8s port binding
  8. Concurrency  
      True. Scaling can be done in horizontal manner with K8s HPA.
  9. Disposability  
      True. Fast startup/graceful shutdown
  10. Dev/prod parity  
      True. Mean to run separately for dev and prod.
  11. Logs  
      True. Logs can be gathered by _filebeat_ or _Grafana Loki_ and forwarded as events.
  12. Admin processes  
      True. Admin/management tasks are running externally (ci/k8s/any other container orchestrator)

* Prove how it fits and uses the best cloud native understanding
  * CI/CD - possible/present
  * OS Independent
  * Autoscaling enabled by adding HPA + Cluster Scaler (for K8s)
  * ...

* How would you expand on this service to allow for the use of an eventstore?  
  Since _BE_ is querying third party GeoIP API we could add Eventstore in order not to face limits of ThirdParty API.

  Flow:
  1. client requests Geoip data
  2. BE puts `geoRequested` event
  3. BE Worker consumes `geoRequested` event
  4. BE Worker requests data from API
  5. BE worker stores data to datastore (DB)
  6. client requests BE to retrieve data from datastore

* How would this service be accessed and used from an external client from the cluster?  
    Following flow:
    1. Client
    2. HTTP2 enabled LB (ALB in AWS) with HTTPS
    3. Internal K8s _Ingress Controller_ (Nginx)
    4. K8s Service -> Pod

    Issues:
    1. Authentication   
       * Can be solved via AWS ALB + Cognito or AWS Lambda
       * By adding additional Auth microservice/layer 
