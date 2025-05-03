docker run -d \
  -e TZ=Etc/UTC \
  -e CONNECTION_TOKEN=123456
  -e SUDO_PASSWORD=123456
  -p 3000:3000 \
  os:test

#  http://<your-ip>:3000/?tkn=supersecrettoken

docker buildx build -f Dockerfile-base -t os:base .

nasm -f bin boot_sect_simple.asm -o boot_sect_simple.bin