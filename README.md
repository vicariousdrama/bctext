# bctext
A simple utility written in Go for passing text to a Blockclock Mini with convenient formatting control

![image](https://user-images.githubusercontent.com/88121568/219843205-f2cb051d-a857-4f10-bec2-b197fedfe9d8.png)


## Clone, build
```
git clone https://github.com/vicariousdrama/bctext.git
cd bctext
go build
```

If you want a smaller binary, you can build with the following on linux
```
GOOS=linux go build -ldflags="-s -w"
upx --brute bctext
```

## Usage
There are currently a few command line options that can be provided
```c
Usage of bctext:
  -blockclockip string
        Blockclock IP Address (default "21.21.21.21")
  -debugmode
        Show data, dont send to blockclock
  -nopadding
        Controls whether periods should be added to pad out the text in panels to edges. With no padding, data in panels is centered
  -texttoshow string
        Text to Show (default "This is sample output created with bctext by vicariousdrama")
  -wordalign
        Avoid breaking word over panels, starting with next panel for each word
```

If -blockclockip is not provided, then debugmode will be toggled on.  Running without any arguments will process the default value in texttoshow as follows

```
$ bctext

Debug results for this text string
---------------------------------------------------------------------------------------
 slot      over            under       url
    0  THIS...         BCTEX           http://21.21.21.21/api/ou_text/0/THIS.../BCTEX
    1  IS..SA          T..BY..         http://21.21.21.21/api/ou_text/1/IS..SA/T..BY..
    2  MPLE.           VICARI          http://21.21.21.21/api/ou_text/2/MPLE./VICARI
    3  OUTPU           OUSDR           http://21.21.21.21/api/ou_text/3/OUTPU/OUSDR
    4  T..CRE          AMA...          http://21.21.21.21/api/ou_text/4/T..CRE/AMA...
    5  ATED..          ............    http://21.21.21.21/api/ou_text/5/ATED../............
    6  WITH..          ............    http://21.21.21.21/api/ou_text/6/WITH../............
```

The combination of -nopadding and -wordalign provides for 4 different output modes.

Passing -nopadding will result in output tha will be centered when rendered on the actual blockclock

```
$ bctext -nopadding

Debug results for this text string
---------------------------------------------------------------------------------------
 slot      over            under       url
    0  THISIS          TBYVI           http://21.21.21.21/api/ou_text/0/THISIS/TBYVI
    1  SAMP            CARIO           http://21.21.21.21/api/ou_text/1/SAMP/CARIO
    2  LEOUT           USDRA           http://21.21.21.21/api/ou_text/2/LEOUT/USDRA
    3  PUTCR           MA              http://21.21.21.21/api/ou_text/3/PUTCR/MA
    4  EATED                           http://21.21.21.21/api/ou_text/4/EATED/
    5  WITH                            http://21.21.21.21/api/ou_text/5/WITH/
    6  BCTEX                           http://21.21.21.21/api/ou_text/6/BCTEX/
```

Passing -wordalign will force the beginning of each new word to start with the next panel

```
$ bctext -wordalign

Debug results for this text string
---------------------------------------------------------------------------------------
 slot      over            under       url
    0  THIS...         CREAT           http://21.21.21.21/api/ou_text/0/THIS.../CREAT
    1  IS.......       ED......        http://21.21.21.21/api/ou_text/1/IS......./ED......
    2  SAMP            WITH..          http://21.21.21.21/api/ou_text/2/SAMP/WITH..
    3  LE......        BCTEX           http://21.21.21.21/api/ou_text/3/LE....../BCTEX
    4  OUTPU           T........       http://21.21.21.21/api/ou_text/4/OUTPU/T........
    5  T........       BY......        http://21.21.21.21/api/ou_text/5/T......../BY......
    6  ..........      ..........      http://21.21.21.21/api/ou_text/6/........../..........
```

Using both -wordalign and -nopadding will start words with the next panel, and end up centered.

```
$ bctext -wordalign -nopadding

Debug results for this text string
---------------------------------------------------------------------------------------
 slot      over            under       url
    0  THIS            CREAT           http://21.21.21.21/api/ou_text/0/THIS/CREAT
    1  IS              ED              http://21.21.21.21/api/ou_text/1/IS/ED
    2  SAMP            WITH            http://21.21.21.21/api/ou_text/2/SAMP/WITH
    3  LE              BCTEX           http://21.21.21.21/api/ou_text/3/LE/BCTEX
    4  OUTPU           T               http://21.21.21.21/api/ou_text/4/OUTPU/T
    5  T               BY              http://21.21.21.21/api/ou_text/5/T/BY
    6                                  http://21.21.21.21/api/ou_text/6//
```

A more practical example of this feature is for displaying columnar type data

```
$ bctext -wordalign -nopadding -texttoshow "MON TUE WED THU FRI SAT SUN 20 21 22 23 24 25 26"

Debug results for this text string
---------------------------------------------------------------------------------------
 slot      over            under       url
    0  MON             20              http://21.21.21.21/api/ou_text/0/MON/20
    1  TUE             21              http://21.21.21.21/api/ou_text/1/TUE/21
    2  WED             22              http://21.21.21.21/api/ou_text/2/WED/22
    3  THU             23              http://21.21.21.21/api/ou_text/3/THU/23
    4  FRI             24              http://21.21.21.21/api/ou_text/4/FRI/24
    5  SAT             25              http://21.21.21.21/api/ou_text/5/SAT/25
    6  SUN             26              http://21.21.21.21/api/ou_text/6/SUN/26
```
