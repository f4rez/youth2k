
import matplotlib as mpl
import numpy as np
import matplotlib.pyplot as plt
import matplotlib.cbook as cbook

def read_datafile(file_name):
    # the skiprows keyword is for heading, but I don't know if trailing lines
    # can be specified
    data = np.loadtxt(file_name, delimiter=',')
    return data


dataL = read_datafile("valuesL.csv")
x = dataL[:,1]
y = dataL[:,0]

#dataR = read_datafile("valuesL.csv")
#x2 = dataR[:,1]
#y2 = dataR[:,0]

fig = plt.figure()
ax1 = fig.add_subplot(111)
ax1.plot(x,y)
ax1.legend()

# fig2 = plt.figure()
# ax2 = fig2.add_subplot(111)
# ax2.plot(x2,y2)
# ax2.legend()

# dataL2 = read_datafile("valuesL.csv")
# x22 = dataL2[:,1]
# y22 = dataL2[:,0]

# dataR2 = read_datafile("valuesL.csv")
# x23 = dataR2[:,1]
# y23 = dataR2[:,0]

# fig22 = plt.figure()
# ax12 = fig2.add_subplot(111)
# ax12.plot(x,y)
# ax12.legend()

# fig23 = plt.figure()
# ax23 = fig23.add_subplot(111)
# ax23.plot(x23,y23)
# ax23.legend()


# dataL4 = read_datafile("valuesL.csv")
# x4= dataL4[:,1]
# y4 = dataL4[:,0]

# dataR5 = read_datafile("valuesL.csv")
# x25 = dataR5[:,1]
# y25 = dataR5[:,0]

# fig6 = plt.figure()
# ax16 = fig6.add_subplot(111)
# ax16.plot(x4,y4)
# ax16.legend()

# fig27 = plt.figure()
# ax27 = fig27.add_subplot(111)
# ax27.plot(x25,y25)
# ax27.legend()










plt.show()
