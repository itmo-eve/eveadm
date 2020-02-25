#!/usr/bin/python3

import sys 
from pexpect import *
import json


def xl_list(eve):
    vms = []
    eve.sendline('')
    eve.expect('\(ns: pillar\) (.*):/#')
    eve.sendline("xl list")
    eve.expect("xl list\r\r\n")
    eve.expect("Name\s+ID\s+Mem\s+VCPUs\s+State\s+Time\(s\)\r\r\n")
    while 1:
        res = eve.expect(["(\S+)\s+\S*\s+\S*\s+\S*\s+(\S*)\s+(\S*)\r\r\n",'\(ns: pillar\) (.*):/#'])
        if res == 0:
            name = eve.match.groups()[0].decode('utf-8')
            if name != "Domain-0":
                state = eve.match.groups()[1].decode('utf-8')
                time = eve.match.groups()[2].decode('utf-8')
                vms.append({"name":name, "state":state, "time":time})
        else:
            return vms
        

def connect(port):
    eve=spawn('ssh -o StrictHostKeyChecking=no -p %s localhost' % port, timeout=None)
    eve.sendline('')
    eve.expect('\(ns: pillar\) (.*):/#')

    return eve

def setup_log(eve):
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
    
def log_json(data):
    try:
        j=json.loads(data)
        func = j["func"].split('/')[-1]
        if j["file"].find("/pillar/cmd/domainmgr/domainmgr.go") == 0 and func in funcs.keys():
            msg = j["msg"]
            time = j['time']
            out = funcs[func](func, msg, time)
            return(out)
    except json.JSONDecodeError:
        #print("JSONDecodeError:", sys.exc_info())
        #s = eve.readline().decode('utf-8')
        #print(s)
        return(None)
    except:
        print("ERR: '%s'" % s)
        print(sys.exc_info())
        sys.exit()

def info(name, msg, time):
    return("(%-30s) %s: %s" % (time, name, msg))

funcs = {
    "domainmgr.handleCreate": info,
    "domainmgr.doActivate": info,
    "domainmgr.doActivateTail": info,
    "domainmgr.doInactivate": info,
    "domainmgr.DomainDestroy": info,
    "domainmgr.handleDelete": info,
    "domainmgr.handleModify": info,
    "domainmgr.xlDomid": info
    }

if __name__ == '__main__':
    if len(sys.argv) < 2:
        print("Usage: %s port [app-name]" % sys.argv[0])
        sys.exit()

    eve = connect(sys.argv[1])
    setup_log(eve)
    evesh = connect(sys.argv[1])

    while 1:
        s = log_json(eve.readline().decode('utf-8'))
        if s:
            print(s)
            vms = xl_list(evesh)
            print("xl list: %s" % vms)
