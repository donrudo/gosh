rm plugins/*
pushd cmd > /dev/null

echo ""

for f in *; { 
    go build -buildmode=plugin  -o ../plugins/$f.cmd.so $f
    echo "BUILD: $f.cmd.so";}

echo ""
popd > /dev/null
go build
