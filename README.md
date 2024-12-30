

# Compile
`tinygo build -target=pico -opt=1 -stack-size=8kb -size=short -o main.uf2 .`

It might be a good idea to remove the old main.uf2 first 

# Flash to pico
```
cp main.uf2 /media/nchris/RPI-RP2
```