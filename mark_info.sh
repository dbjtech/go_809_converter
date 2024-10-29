#! /bin/sh
dist_path="git_info.txt"
branch_name="$(git branch --show-current)"
git_hash="$(git rev-parse HEAD)"
touch $dist_path
if [ -n "$git_hash" ]; then
  echo "hash=$git_hash" > $dist_path
fi
if [ -n "$branch_name" ]; then
  echo "branch=$branch_name" >> $dist_path
fi
if [ "$(type date)" ]; then
  timeZone="$(date "+%Z")"
  now=$(date "+%F %T")
  # 若不是北京时间，则加8个小时
  if [ "$timeZone" != "CST" ]; then
    # date 命令不支持+8 hours操作。所以只能转时间戳然后加秒操作
    timestamp=$(date "+%s")
    seconds_new=$(( 28800 + timestamp )) #  加8小时
    now=$( date -d @$seconds_new "+%F %T") # @ 符号能直接把时间戳转成时间
  fi
  echo "build_at=$now" >> $dist_path
fi
