<h1 align="center">auto53 📂  </h1>

<h5 align="center">The missing link between AWS AutoScaling Groups and Route53</h5>

<br/>

[![Build Status](https://travis-ci.com/cirocosta/auto53.svg?token=ixZ9XiEPW4YH62ixq7Av&branch=master)](https://travis-ci.com/cirocosta/auto53)


### Overview

`auto53` solves the issue of keeping a route53 zone up to date with the changes that an autoscaling group might face. 

For instance, consider the following state in EC2:

```yaml
# internal state retrieved from the 
# inspection of EC2
autoscaling_groups:
  - name: asg1
    instances:
      - id: 'i-0123'
        private_ip: '10.0.0.2'
      - id: 'i-0321'
        private_ip: '10.0.0.3' 
```

and also the following formatting configuration:

```yaml
# for each zone we can tie several
# autoscaling groups that present
# an automatic record creation rule.
--- 
- AutoScalingGroup: 'asg1'
  Zone': 
    ID: 'zone123'
    Name: 'ciro-test'
  Record: 'asg1-machines'

- AutoScalingGroup: 'asg1'
  Zone': 
    ID: 'zone123'
    Name: 'ciro-test'
  Record: '{{ .Id }}-machine'
```

with that we'd end up with the following records:

```yaml
asg1-machines.ciro-test:
  - 10.0.0.2
  - 10.0.0.3

i-0123-machine.ciro-test:
  - 10.0.0.2

i-0321-machine.ciro-test:
  - 10.0.0.3
```

Now we could consider the situation where `asg1` scales out to `3` instances while at the same time having the death of `i-0321`:


```diff
 autoscaling_groups:
   asg1:
     machines:
       - id: 'i-0123'
         private_ip: '10.0.0.2'
-      - id: 'i-0321'
-        private_ip: '10.0.0.3' 
+      - id: 'i-0444'
+        private_ip: '10.0.0.4' 
+      - id: 'i-0555'
+        private_ip: '10.0.0.5' 
```

That would mean that `auto53` would notice the change in the desired state and then update `route53` accordingly:



```yaml
asg1-machines.ciro-test:
  - 10.0.0.2
  - 10.0.0.4
  - 10.0.0.5

i-0123-machine.ciro-test:
  - 10.0.0.2

i-0444-machine.ciro-test:
  - 10.0.0.4

i-0555-machine.ciro-test:
  - 10.0.0.5
```

### Usage

`auto53` aims at being a single binary that is capable of running in 2 modes:

- server mode: sits in an instance querying route53 state from time to time as well as being (optionally) notified by SNS in case of any changes in an autoscaling group;
- single execution mode: runs once whenever the binary is fired - suitable for executions in the context of a one-off lambda function.

In either case, the necessary user permissions are needed:

- EC2 Describe Instances
- Route53 - ListResourceRecordSets, ChangeResourceRecordSets

The AWS credentials are accessed via the default behavior of AWS CLI (either environment variables or config file under `~/.aws`).

```
Usage: auto53 [opts ...]

Options:
  --config CONFIG [default: ./auto53.yaml]
  --debug                activates debug-level logging
  --dry                  run without performing modifications
  --interval INTERVAL [default: 2m0s]
  --listen
  --once                 run one time and exit
  --port PORT [default: 8080]
  --help, -h             display this help and exit
```

