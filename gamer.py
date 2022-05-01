import pandas as pd
import numpy as np
from time import sleep


baseTradePrice = 4.97

def dictionary(items, n):
  dicti = {}
  length = len(items)
  for i in range(length):
    if i+n+1 <= length:
      for j in range(i+1,i+n+1):
        if items[i:j] in dicti:
          dicti[items[i:j]] += 1
        else:
          dicti[items[i:j]] = 1
    else:
      for j in range(i+1,length+1):
        if items[i:j] in dicti:
          dicti[items[i:j]] += 1
        else:
          dicti[items[i:j]] = 1
  return dicti

def createString(text):
  sum = ''
  for i in text:
    sum += str(i)
  return sum

def findNext(dicti, ourComb, missed_crash=0, missed1=0, missed2=0, missed3=0, missed4=0, missed5=0):
  fullSum = 0
  max = 0
  maxColor= 0
  crash = [ourComb+'0',0]
  one = [ourComb+'1',0]
  two = [ourComb+'2',0]
  three = [ourComb+'3',0]
  four = [ourComb+'4',0]
  five = [ourComb+'5',0]
  for i in dicti:
    if (ourComb == i[:len(ourComb)]) and (ourComb != i):
      if i[-1] == str(0):
            crash = [i, dicti[i]]
      elif i[-1] == str(1):
            one = [i, dicti[i]]
      elif i[-1] == str(2):
            two = [i, dicti[i]]
      elif i[-1] == str(3):
            three = [i, dicti[i]]
      elif i[-1] == str(4):
            four = [i, dicti[i]]
      elif i[-1] == str(5):
            five = [i, dicti[i]]
      fullSum += dicti[i]
      if dicti[i] >= max:
        max = dicti[i]
        maxColor = i
  if fullSum == 0:
    return [0,0]
  if (crash[1] + 0.19*missed_crash*max >= max):
        return [crash[0],(crash[1] + 0.19*missed_crash*max)/fullSum]
  if (one[1] + 0.32*missed_crash*max >= max):
        return [one[0],(one[1] + 0.32*missed_crash*max)/fullSum]
  if (two[1] + 0.25*missed_crash*max >= max):
        return [two[0],(two[1] + 0.25*missed_crash*max)/fullSum]
  if (three[1] + 0.12*missed_crash*max >= max):
        return [three[0],(three[1] + 0.12*missed_crash*max)/fullSum]
  if (four[1] + 0.08*missed_crash*max >= max):
        return [four[0],(four[1] + 0.08*missed_crash*max)/fullSum]
  if (five[1] + 0.04*missed_crash*max >= max):
        return [five[0],(five[1] + 0.04*missed_crash*max)/fullSum]
  return [maxColor, max/fullSum]


def markovChain(y_color, ourComb, missed_crash=0, missed1=0, missed2=0, missed3=0, missed4=0, missed5=0):
  maxProba = 0
  maxColor = 0
  length = len(ourComb) + 1
  x = createString(y_color)
  for i in range(length-2):
    dicti = dictionary(x, length)
    next, proba = findNext(dicti, ourComb[i:], missed_crash, missed1, missed2, missed3, missed4, missed5)
    if maxProba < proba:
      maxProba = proba
      maxColor = next
    length -= 1
  return [maxColor, maxProba]

def play():
    firstBet = True
    zero = 0
    one = 0
    two = 0
    three = 0
    four = 0
    five = 0
    data = pd.read_excel('data.xlsx')
    y_color = data['Color'] # Цвета
    i = 0
    while True:
        file = open('GOLANG/textOnline.log', 'r') # Parsing by Golang script
        fl = file.read()
        file.close()
        data = ''.join(fl.split())
        pred_x = data[-4:]
        pred, proba = markovChain(y_color, pred_x, zero, one, two, three, four, five)
        val = pred_x[-1]
        if int(val) == 0:
            zero = 0
        else:
            zero += 1
        if int(val) == 1:
            one = 0
        else:
            one += 1
        if int(val) == 2:
            two = 0
        else:
            two += 1
        if int(val) == 3:
            three = 0
        else:
            three += 1
        if int(val) == 4:
            four = 0
        else:
            four += 1
        if int(val) == 5:
            five = 0
        else:
            five += 1
        if (proba > 0.8) and (pred[-1] != '0') and not ('5' in pred_x):
            print("Made a bet.")
            makeBet(firstBet, pred[-1])
            file_object = open('bets.txt', 'a')
            file_object.write(str(pred) + ':' + str(proba)+'\n')
            file_object.close()
            if firstBet:
            	firstBet = False
        else:
        	balance, items = balance_and_items()
        	i += 1
        	print(i, ' ', balance)
	        if items[0] >= baseTradePrice*2:
	            crackItem()
	        if balance >= 70:
	            print("Earned and exited.")
	            break
	        if check_exists_by_xpath('/html/body/div[1]/div[1]/div[2]/div[1]/div[2]/div[2]/div/button[1]'):
	            item = driver.find_element_by_xpath('/html/body/div[1]/div[1]/div[2]/div[1]/div[2]/div[2]/div/button[1]')
	            if 'btn-base drop-preview checked' in item.get_attribute('outerHTML'):
		            item.click()
        while (data[-4:] == pred_x):
            file = open('GOLANG/textOnline.log', 'r')
            fl = file.read()
            file.close()
            data = ''.join(fl.split())

# Steam Code Parser from txt file
# Just copy-paste from your Steam Settings in txt file
def code_parser():
	file = open('codes.txt', 'r')
	fl = file.read()
	file.close()
	fl = fl.split()
	i = int(fl[-1])
	if i <= 30:
		code = fl[i]
		start = code.find('.')
		code = code[start+1:]
		file = open('codes.txt', 'a')
		file.write(' '+str(i+1))
		file.close()
		return code


options = Options()
#options.add_argument('--auto-open-devtools-for-tabs')
options.add_argument('--headless')
options.add_argument('--no-sandbox')
options.add_argument('--disable-gpu')
options.add_argument("window-size=4000,1080")
options.add_argument("--mute-audio")


driver = webdriver.Chrome('./chromedriver', options=options)
driver.get("https://csgorun.gg/")


def main():
	sleep(1)
	sleep(3)
	logIn = driver.find_element_by_xpath('/html/body/div[1]/div[3]/div[1]/div[2]/div/button')
	logIn.click()
	sleep(3)
	login = driver.find_element_by_xpath('/html/body/div[1]/div[7]/div[2]/div/div[2]/div[2]/div/form[1]/input[4]')
	login.send_keys('LOGIN')
	password = driver.find_element_by_xpath('/html/body/div[1]/div[7]/div[2]/div/div[2]/div[2]/div/form[1]/input[5]')
	password.send_keys('PASS')
	signIn = driver.find_element_by_xpath('/html/body/div[1]/div[7]/div[2]/div/div[2]/div[2]/div/form[1]/div[4]/input')
	signIn.click()
	code = code_parser()
	sleep(3)
	codeIn = driver.find_element_by_xpath('/html/body/div[4]/div[3]/div/div/div/form/div[3]/div[1]/div/input')
	codeIn.send_keys(code)
	submit = driver.find_element_by_xpath('/html/body/div[4]/div[3]/div/div/div/form/div[4]/div[1]/div[1]')
	submit.click()
	sleep(7)
	slider = driver.find_element_by_xpath('/html/body/div[1]/div[3]/div[1]/div[2]/div/label/div[1]/div')
	slider.click()
	sleep(1)
	balance, items = balance_and_items()
	print(balance, items)
	if items[0] >= baseTradePrice*2:
		crackItem()
	correctAutoBet()
	play()
	sleep(600)
	driver.close()


main()





