#!/bin/bash
#--------------------------------------------
# 此脚本用户测试环境(52.83.232.6)管理后台前端项目发布
# author：wei
#--------------------------------------------
imageName=${1}
imageTag=${2}
argName=${3}

ssh -tt -p 22 ubuntu@52.83.232.6 "bash -s $imageName $imageTag $argName" <<-'EOF'
cd /var/www/laihua_compose/newadmin/design
sh ./build.sh --${3} ${2}
exit
EOF