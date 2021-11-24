package queue

// JobAddDomain adds a new domain to dns redis
const JobAddDomain = "AddDomain"

// JobRemoveDomain removes a domain from dns redis
const JobRemoveDomain = "RemoveDomain"

// JobUpdateSubDomain updates redis dns with records from db for a given subdomain
const JobUpdateSubDomain = "UpdateSubDomain"
