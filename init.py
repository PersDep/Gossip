#!/usr/bin/python

import subprocess as subproc
import matplotlib.pyplot as pyplot
import numpy as np
import sys

probabilities = [0.0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6]
ticks = []

iters = 32

nodes_amount = "40"
port = "4400"
min_degree = "2"
max_degree = "8"
ttl = "40"

def plot(figsize, color, width):
    axis = np.zeros(len(probabilities))
    for i in range(0, len(probabilities)):
        axis[i] = i

    pyplot.figure(figsize=figsize)
    pyplot.bar(axis, ticks, color=color, width=width)
    pyplot.xticks(axis + width / 2, probabilities)

    pyplot.title('Nodes: ' + nodes_amount + ', Node degree: ' + min_degree + ' - ' + max_degree + ', TTL: ' + ttl)
    pyplot.xlabel("Loss probability")
    pyplot.ylabel("Time (ticks)")
    pyplot.savefig("plot.jpg")

def my_command(chance):
    return 'INPUT -p udp -i lo -m statistic --mode random --probability ' + str(chance) + ' -j DROP'

def bash_call(cmd):
    result, error = subproc.Popen(cmd.split(), stdout=subproc.PIPE).communicate()
    if error:
        print error
        return error
    return result

if len(sys.argv) == 6:
    nodes_amount = sys.argv[1]
    port = sys.argv[2]
    min_degree = sys.argv[3]
    max_degree = sys.argv[4]
    ttl = sys.argv[5]

for probability in probabilities:
    cur_ticks = 0.0

    bash_call('iptables -A' + my_command(probability))

    cur_iter = iters
    while cur_iter > 0:
        res = bash_call('./main ' + nodes_amount + ' ' + port + ' ' + min_degree + ' ' + max_degree + ' ' + ttl).split()
        if res and res[0] == 'Finished':
            cur_ticks += float(res[2])
            cur_iter -= 1

    bash_call('iptables -D' + my_command(probability))

    ticks.append(cur_ticks / iters)
    print ticks


plot(figsize=(16, 9), color='g', width=0.3)
