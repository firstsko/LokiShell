# 🔷 LOKISHELL

**Next Generation Log Console**

- **Schema**: Powered by ARMS Product–System–Application CMDB  
- **Filter**: Bash-like syntax, blazing fast CLI filtering  
- **Search**: Advanced keyword queries with Loki-compatible syntax  

---
![image](https://github.com/user-attachments/assets/8c7e73ce-ddac-48d3-9ff4-698b80816ffb)

##  Access Control

```bash
login        # Authenticate your user account


## Service Tree Commands
ls           # List services (flat mode)
ls -l        # List services (detailed mode)
ls -p        # List services by product group (grouped mode)
ls -r        # List services by region (grouped mode)
ls -a        # List services (tree view)

cd /devops   # Navigate to a specific path (format: /ARMS-Product/System/App)

# Quick jump to specific application logs
go hoc-openresty-ser

##  Log Viewing
more         # Paginated forward view of logs
less         # Paginated reverse view from current time
tailt        # Real-time streaming of logs (full tail mode)

## Filter Conditions
### Add Filter
filter -a _ip=172.19.38.71
filter -a filename=/app/docker/fpc/logs/sql.2022-05-06-10.log
filter -a _ip=172.19.38.71,172.19.253.166

### Delete Filter
filter -d _ip=192.168.16.77
filter -d filename=/app/docker/fpc/logs/sql.2022-05-06-10.log
filter -d _ip=172.19.38.71,172.19.253.166

Remove filters by IP or filename.

### Add & Delete Combined
filter -a _ip=192.168.16.77 -d filename=/app/docker/fpc/logs/sql.2022-05-06-10.log



# compiler linux

```
bash build-linux.sh
```

complier windows

```
bash bukld-win.sh
```

execute shell
```
./dist/loki-shell
```
