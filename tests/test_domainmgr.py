#!/usr/bin/python3

import sys 
from pexpect import *
import json

if len(sys.argv) < 2:
    print("Usage: %s port [app-name]" % sys.argv[0])
    sys.exit()

eve=spawn('ssh -p %s localhost' % sys.argv[1], timeout=None)
eve.sendline('')
eve.expect('\(ns: pillar\) (.*):/#')

eve.sendline("cat /var/run/domainmgr.pid")
res = eve.expect(['([0-9]+)\(ns: pillar\) (.*):/#','cat: /var/run/domainmgr.pid: No such file or directory'])
if res == 0:
    pid = eve.match.groups()[0].decode('utf-8')
    print("domainmgr PID: ", pid)
    eve.sendline("kill " + pid)
if res == 1:
    print(eve.after.decode('utf-8'))
    eve.close()
    sys.exit(-1)

eve.expect('\(ns: pillar\) (.*):/#')
eve.sendline('/opt/zededa/bin/domainmgr -c IMGA -d >/tmp/domainmgr.log 2>&1 &')
eve.expect('\(ns: pillar\) (.*):/#')
eve.sendline('tail -f /tmp/domainmgr.log')
eve.expect('tail -f /tmp/domainmgr.log')

while 1:
    s = eve.readline().decode('utf-8')
    try:
        #print(s)
        j=json.loads(s)
        if j["file"] == "/pillar/pubsub/publish.go:288":
            t=j['time']
            j=json.loads("{" + j['msg'].split(" val {")[1])
            #print(j.keys())
            if 'DisplayName' in j.keys():
                if len(j['LastErr']):
                    out = "(%30s) %s Activated: %s Error: '%s' (%s)" % (t, j['DisplayName'], j['Activated'], j['LastErr'], j['LastErrTime'])
                else:
                    out = "(%30s) %s Activated: %s" % (t, j['DisplayName'], j['Activated'])
                if len(sys.argv) >=3 and j["DisplayName"] == sys.argv[2]:
                    print(out)
                else:
                    print(out)
    except json.JSONDecodeError:
        pass
        #print("JSONDecodeError:", sys.exc_info())
        #print(s)
    except:
        print("ERR: '%s'" % s)
        print(sys.exc_info())
        sys.exit()
