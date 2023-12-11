#!/bin/bash

echo "Welcome!"
echo '' 

echo "This tool will not actually do anything to your data. It will create a script that you can submit to Slurm to move your data. Avoid running the script on a login node. 

Please examine the resulting script before using it. If you feel uncertain, ask UPPMAX support to look it over."

max_files_per_dir=1000

# TODO: list possible project directories

read -p 'Move which directory? [default: this one]: ' dir_to_move
dir_to_move=${dir_to_move:-"$PWD"}

echo Moving $dir_to_move

echo ''

echo "This tool will find all subdirectories with more than $max_files_per_dir files in them and package (tar) them before moving."

read -p 'Should we discard the large subdirectories after packaging? [Y/n]: ' keep_dirs
keep_dirs=${keep_dirs:-N}
# ^ sign converts to upper case
keep_dirs=${keep_dirs^}

if [[ $keep_dirs = "N" ]]; then
   echo "We will discard the big directories after packaging."
else
   echo "We will keep the big directories."
fi

echo ''

read -p "Do you wish to automatically delete local files from $HOSTNAME after copying them? [y/N]: " auto_del
auto_del=${auto_del:-N}
auto_del=${auto_del^}


if [[ $auto_del = "Y" ]]; then
   echo "We will delete files that have been copying."
else
   echo "We will keep files here after copying."
fi

echo ''
read -p 'Which system should data be moved to? [default: dardel.pdc.kth.se]' target_host
target_host=${target_host:-"dardel.pdc.kth.se"}


echo ''
read -p "Where on $target_host should data be moved to? " target_dir
# Important: do better sanity checks

if [[ ! $target_dir ]]; then
   echo "You must supply a target directory."
   exit 1
fi


echo ''
username=`whoami`
read -p "What is your user name on $target_host? [$username]: " username
username=${username:-`whoami`}

echo ''
read -p 'How many parallel rsync connections? [10]: ' num_conns
num_conns=${num_conns:-10}


# Time to construct that sbatch script

script_name="transfer_$(basename $dir_to_move).sh"

echo "#!/bin/bash" > $script_name
echo "#SBATCH -p core" >> $script_name
echo "#SBATCH -n 1" >> $script_name
echo "#SBATCH -J $scriptname" >> $script_name
echo "#SBATCH -A YOUR_PROJECT" >> $script_name
echo "#SBATCH -t 7-00:00:00" >> $script_name

echo "find $dir_to_move -mindepth 1 -maxdepth 2 -not -path '*/.*' -type d -links $max_files_per_dir > large_directories.txt" >> $script_name

echo ''
# TO DO: parallelise this. 
# TO DO: consider what happens if directories are changing as this is done. 
echo "xargs -a large_directories.txt -I {} tar -czvf {}.tar.gz {}" >> $script_name

echo ''

if [[ $keep_dirs = "N" ]]; then
   echo "xargs -a large_directories.txt -I{} rm -rf {}" >> $script_name
   echo ''
fi

echo '' 
echo "rsync -cavz --progress --parallel=$num_conns --exclude-from=large_directories.txt $dir_to_move $username@$target_host:$target_dir | tee rsync_log.txt" >> $script_name

# TO DO: if clean-up is required, read the rsync log, parse out successful transfers, and delete that stuff.


# final instructions to the user

echo ""
echo "When you are ready, edit $script_name to set the correct project ID and run \"sbatch $script_name\"."
