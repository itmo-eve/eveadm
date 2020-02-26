#!/usr/bin/python3

import sys 
from pexpect import *
import json
import re

apps={}

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
            level = j['level']
            out = funcs[func](func, level, msg, time)
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

def add_app(name):
    if not name in apps:
        apps[name] = {}
        
def handleCreate(name, level, msg, time):
    res = re.search("handleCreate\(\{.*\}\) DONE for (.*)", msg)
    if res:
        app = res.group(1)
        add_app(app)
        print("%s: %s" % (name, app))
        apps[app]['status'] = "created"
        apps[app]['active'] = 0
        return("(%-30s) [%s] '%s' created" % (time, level, app))

def doActivate(name, level, msg, time):
    res = re.search("DomainCreate for (.*)\.\d: (.*)", msg)
    if res:
        app = res.group(1)
        add_app(app)
        print("%s: %s" % (name, app))
        msg = res.group(2)
        apps[app]['message'] = msg
        return("(%-30s) [%s] %s message: %s" % (time, level, app, msg))

def doActivateTail(name, level, msg, time):
    res = re.search("created domainID \d* for (.*)\.\d", msg)
    if res:
        app = res.group(1)
        add_app(app)
        print("%s: %s" % (name, app))
        for d in xl_list(evesh):
            if d["name"].find(app) == 0:
                apps[app]['active'] = 1
                return("(%-30s) [%s] '%s' activated" % (time, level, app))

def doInactivate(name, level, msg, time):
    res = re.search("doInactivate\(\{.*\}\) done for (.*)", msg)
    if res:
        app = res.group(1)
        add_app(app)
        print("%s: %s" % (name, app))
        apps[app]['active'] = 0
        for d in xl_list(evesh):
            if d["name"].find(app) == 0:
                apps[app]['active'] = 1
        if apps[app]['active'] == 0:
            return("(%-30s) [%s] '%s' deactivated" % (time, level, app))

def handleDelete(name, level, msg, time):
    res = re.search("handleDelete\(\{.*\}\) DONE for (.*)", msg)
    if res:
        app = res.group(1)
        add_app(app)
        print("%s: %s" % (name, app))
        for d in xl_list(evesh):
            if d["name"].find(app) == 0:
                app = None
        if app:
            del apps[app]
            return("(%-30s) [%s] '%s' deleted" % (time, level, app))

def info(name, level, msg, time):
    return("(%-30s) [%s] %s: %s" % (time, level, name, msg))

funcs = {
#    "domainmgr.DomainCreate": info,
#    "domainmgr.handleCreate": info,
    "domainmgr.handleCreate": handleCreate,
#    "domainmgr.doActivate": info,
    "domainmgr.doActivate": doActivate,
#    "domainmgr.doActivateTail": info,
    "domainmgr.doActivateTail": doActivateTail,
#    "domainmgr.doInactivate": info,
    "domainmgr.doInactivate": doInactivate,
#    "domainmgr.DomainDestroy": info,
#    "domainmgr.handleDelete": info,
    "domainmgr.handleDelete": handleDelete,
#    "domainmgr.handleModify": info,
#    "domainmgr.xlDomid": info
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
            #vms = xl_list(evesh)
            #print("xl list: %s" % vms)
