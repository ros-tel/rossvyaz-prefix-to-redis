# rossvyaz-prefix-to-redis
---------------------------
Loading prefixes of telecom operators in Radis according to rossvyaz.ru

### Obtaining data from the site [www.rossvyaz.ru](https://www.rossvyaz.ru)
```
wget https://www.rossvyaz.ru/docs/articles/Kody_ABC-3kh.csv
wget https://www.rossvyaz.ru/docs/articles/Kody_ABC-4kh.csv
wget https://www.rossvyaz.ru/docs/articles/Kody_ABC-8kh.csv
wget https://www.rossvyaz.ru/docs/articles/Kody_DEF-9kh.csv
```

### List of telecom operators
```
cat Kody_DEF-9kh.csv | \
iconv -f cp1251 -t utf8 | \
awk -F ";" '{print $5}' | \
sed 's#\t##g' | \
sort | \
uniq
```

### List of regions
```
cat Kody_DEF-9kh.csv | \
iconv -f cp1251 -t utf8 | \
awk -F ";" '{print $6}' | \
sed 's#\t##g' | \
sort | \
uniq
```

### Load testing
```
tail -n+2 Kody_DEF-9kh.csv | \
head | \
iconv -f cp1251 -t utf8 | \
rossvyaz-to-prefix -debug
```

### Full data upload
For complete cleaning, you can add the **-flushall** parametr when loading the first file.
```
tail -n+2 Kody_DEF-9kh.csv | \
iconv -f cp1251 -t utf8 | \
rossvyaz-to-prefix -flushall

tail -n+2 Kody_ABC-3kh.csv | \
iconv -f cp1251 -t utf8 | \
rossvyaz-to-prefix
...
```

### Preview of prefixes
```
redis-cli
> KEYS 79001*
 1) "7900193"
 2) "790010"
 3) "790013"
 4) "790018"
 5) "7900190"
 6) "7900192"
 7) "790015"
 8) "790016"
 9) "7900196"
10) "7900195"
11) "790017"
12) "7900194"
13) "790011"
14) "790014"
15) "790012"
16) "790019"
17) "7900191"
```

##### An example of a search query with a sort by the longest match length with the prefix
```
redis-cli
> MGET 79001812345 7900181234 790018123 79001812 7900181 790018 79001 7900
1) (nil)
2) (nil)
3) (nil)
4) (nil)
5) (nil)
6) "\xd0\x9e\xd0\x9e\xd0\x9e \"\xd0\xa22 \xd0\x9c\xd0\xbe\xd0\xb1\xd0\xb0\xd0\xb9\xd0\xbb\";\xd0\xa1\xd1\x82\xd0\xb0\xd0\xb2\xd1\x80\xd0\xbe\xd0\xbf\xd0\xbe\xd0\xbb\xd1\x8c\xd1\x81\xd0\xba\xd0\xb8\xd0\xb9 \xd0\xba\xd1\x80\xd0\xb0\xd0\xb9"
7) (nil)
8) (nil)
```
