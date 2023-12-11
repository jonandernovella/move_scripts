#!/bin/bash
#SBATCH -p core
#SBATCH -n 1
#SBATCH -J transfer_move_scripts.sh
#SBATCH -A YOUR_PROJECT
#SBATCH -t 7-00:00:00

find /Users/jon/Documents/nbis/score/move_scripts -mindepth 1 -maxdepth 2 -not -path '*/.*' -type d -links 100000 > large_directories.txt

xargs -a large_directories.txt -I {} tar -czvf {}.tar.gz {}

rsync -cavz --progress --parallel=3 --exclude-from=large_directories.txt /Users/jon/Documents/nbis/score/move_scripts jon@dardel.pdc.kth.se:/test/fs/di | tee rsync_log.txt
