#!/usr/bin/python

import subprocess as subproc
import matplotlib.pyplot as pyplot
import numpy as np

probabilities = [0.0, 0.1, 0.2, 0.3, 0.4, 0.5]#, 0.6]
ticks = []

iters = 4

nodes_amount = 10
port = 4400
min_degree = 2
max_degree = 5
ttl = 40

def plot(figsize, color, width):
    axis = np.zeros(len(probabilities))
    for i in range(0, len(probabilities)):
        axis[i] = i

    pyplot.figure(figsize=figsize)
    pyplot.bar(axis, ticks, color=color, width=width)
    pyplot.xticks(axis + width / 2, probabilities)

    pyplot.title(str(nodes_amount) + ' nodes_amount, ' + str(min_degree) + ' - ' + str(max_degree) + ' degree, ' + str(ttl) + ' TTL')
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


for probability in probabilities:
    cur_ticks = 0.0

    bash_call('iptables -A' + my_command(probability))

    cur_iter = iters
    while cur_iter > 0:
        res = bash_call('./main ' + str(nodes_amount) + ' ' + str(port) + ' ' + str(min_degree) + ' ' + str(max_degree) + ' ' + str(ttl)).split()
        if res and res[0] == 'Finished':
            cur_ticks += float(res[2])
            cur_iter -= 1

    bash_call('iptables -D' + my_command(probability))

    ticks.append(cur_ticks / iters)
    print ticks


plot(figsize=(16, 9), color='g', width=0.3)
