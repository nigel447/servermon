## Server Monitor via Simple SSH Applicaton

### Standard libary with EC25519 authentication

You can use this repo to bootstrap a ssh pipe with a desktop client  
This example displays a remote servers realtime  system metrics  

Code to mashall [EC25519](https://cr.yp.to/ecdh.html) keys for the ssh pipe  

Simple front end client  
<img src="https://github.com/nigel447/servermon/blob/master/serverMon.gif" width="400"> 

Remote Output  
<img src="https://github.com/nigel447/servermon/blob/master/remote_smon.gif" width="400">  

Above test on a Oracle Cloud VM with single [Ampere A1](https://www.oracle.com/cloud/compute/arm/) CPU.  
Simple golang code uses a tiny fraction << 1% of the cpu,and approx 300kb of system memory.   
Desktop client with [fyne](https://github.com/fyne-io/fyne)