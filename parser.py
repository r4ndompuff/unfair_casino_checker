import pandas as pd

fl = open('GOLANG/text.log', 'r')
fr = fl.readlines()
gameId = []
ticks = []
money = []
gamers = []
skins = []
moder = []
for i in fr:
	text = i.split(';')
	text[-1] = text[-1][:-1]
	gameId.append(text[0])
	gamers.append(text[1])
	skins.append(text[2])
	money.append(text[3])
	ticks.append(text[4])
	#moder.append(text[6])
	print(text)

dicti = {'ID': gameId, "gamers": gamers, "skins": skins, "money": money, "ticks": ticks}
#"peopleWin": gamers, "peopleLost": skins,
df = pd.DataFrame(dicti)
df.to_csv('dataGO.csv', index=False)