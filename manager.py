#!/usr/bin/env python
# -*- coding: UTF-8 -*-

import MySQLdb
import sys
import getopt

def usage():
    print "manager.py usage:"
    print "    -h, --help: print help message"
    print "    --list: 列出某个tag的server，all 为全部"
    print "    --addr: server addr"
    print "    --time: time"

def echo_server(s):
    get_c_sql = "select distinct(cluster) from db_account"
    get_h_sql = "select addr from db_account where cluster = %s"
    if s == "all":
        cursor.execute(get_c_sql)
        cs = cursor.fetchall()
        for row in cs:
            print "\033[0;31;1m%s\033[0m" % row[0]
            cursor.execute(get_h_sql, row[0])
            hs = cursor.fetchall()
            for h in hs:
                print "    %s" % h[0]
    else:
        cursor.execute(get_h_sql, s)
        hs = cursor.fetchall()
        for h in hs:
            print h[0]

def echo_status(addr, t):
    id = get_id(addr)
    sql = "(select from_unixtime(created_time), com_insert,\
                   com_update, com_delete, com_select, qps,\
                   byte_received, byte_sent\
                   from db_status where id = %s\
                   and created_time >= (unix_timestamp() - %s * 3600)\
                   order by created_time)\
           union all\
           (select '                max', max(com_insert), max(com_update),\
                   max(com_delete), max(com_select), max(qps),\
                   max(byte_received),max(byte_sent)\
                   from db_status where id = %s\
                   and created_time > (unix_timestamp() - %s * 3600))\
           union all\
           (select '                avg', avg(com_insert), avg(com_update),\
                   avg(com_delete), avg(com_select), avg(qps),\
                   avg(byte_received),avg(byte_sent)\
                   from db_status where id = %s\
                   and created_time > (unix_timestamp() - %s * 3600))\
           union all\
           (select '                min', min(com_insert), min(com_update),\
                   min(com_delete), min(com_select), min(qps),\
                   min(byte_received),min(byte_sent)\
                   from db_status where id = %s\
                   and created_time > (unix_timestamp() - %s * 3600))" % (id, t,id,t,id,t,id,t)
    cursor.execute(sql)
    i=0
    l = ("time|", "ins", "upd", "del", "sel", "qps|", "recv", "send")
    for row in cursor.fetchall():
        if i % 20 == 0:
            print "\033[0;46;1m%20s%5s%5s%5s%6s%7s%8s%9s\033[0m" % l
        print "\033[0;33;1m%s\033[0m|%5d%5d%5d%6d%6d|%8d%9d" % row 
        i += 1


def get_id(addr):
    cursor.execute("select id from db_account where addr = %s", addr)
    return cursor.fetchone()[0]

def get_conn():
    db = MySQLdb.connect(
            host="10.10.35.205",
            user="test",
            passwd="test",
            port=3306,
            db="dbmonitor",
            charset='utf8',
            )
    return db

def main(args):
    try:
        opts, other = getopt.getopt(args, "h",
                ["addr=", "time=", "list=", "help"])
    except getopt.GetoptError, err:
        print str(err)
        usage()
        sys.exit(2)
    if len(opts) == 0:
        usage()
        sys.exit(1)

    addr = ""
    t = 1
    for o, a in opts:
        if o in ("-h", "--help"):
            usage()
            sys.exit(1)
        elif o == "--list":
            echo_server(a)
            sys.exit(1)
        elif o == "--time":
            t = a
        elif o == "--addr":
            addr = a
        else:
            print "unhandled option"
            sys.exit(3)
    echo_status(addr, t)


if __name__ == '__main__':
    ''' '''

    db = get_conn()
    cursor = db.cursor()
    cursor.execute("set autocommit=1")
    main(sys.argv[1:])
