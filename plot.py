
import matplotlib as mpl
import numpy as np
import matplotlib.pyplot as plt
import matplotlib.cbook as cbook

def read_datafile(file_name):
    # the skiprows keyword is for heading, but I don't know if trailing lines
    # can be specified
    data = np.loadtxt(file_name, delimiter=',')
    return data


data = read_datafile("values.csv")
x = data[:,1]
y = data[:,0]

fig = plt.figure()
ax1 = fig.add_subplot(111)
ax1.plot(x,y)
ax1.legend()
plt.show()
