killall blogServer_linux
rm -rf ./nohup.out
nohup ./blogServer_linux -releaseMod=false -protocol=https &
