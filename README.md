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
formatting_rules:
  - autoscaling_group: 'asg1'
    zone': 'ciro-test'
    record: 'asg1-machines'

  - autoscaling_group: 'asg1'
    zone': 'ciro-test'
    record: '{{ .Id }}-machine'
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

Now we could consider the situation where `asg1` scales out to `3` instances while at the same time having the death of `i-0321:


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

