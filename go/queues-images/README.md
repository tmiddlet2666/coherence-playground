# Queues Image Demo



## Populate the urls.txt file


```bash
for i in $(seq 1 25000); do  echo "http://localhost:8080/image?text=Hello$i${RANDOM}-${RANDOM}"; done > ./urls.txt
```