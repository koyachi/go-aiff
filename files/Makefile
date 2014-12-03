all: a.aiff

a.wav:
	curl -L https://github.com/youpy/go-wav/blob/master/files/a.wav?raw=true -o a.wav

a.aiff: a.wav
	afconvert -f AIFF -d BEI16 a.wav -d BEI16 $@
