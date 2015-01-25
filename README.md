# png2vplot

Convert a png image to a vplot file

## Developer Setup

    cd $GOPATH
    git clone git@github.com:ricallinson/png2vplot.git ./src/github.com/ricallinson/png2vplot

## Install

    cd ./src/github.com/ricallinson/png2vplot
    go install

## Use

	png2vplot -h
    png2vplot ./fixtures/test.png ./test.vplot
    png2vplot -x 12000 -y 6000 -p 100 ./fixtures/test.png ./test.vplot
