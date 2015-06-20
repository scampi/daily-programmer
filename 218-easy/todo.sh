#! /usr/bin/env bash

## https://www.reddit.com/r/dailyprogrammer/comments/39ws1x/20150615_challenge_218_easy_todo_list_part_1/
## this creates a named pipe through which you send actions
## - add xxx: add xxx to the todo list
## - delete X: deletes the item number X in the list
## - view: prints the todo list

todos=()

mkfifo todo.fifo

while read line; do
    action=$(echo $line | cut -f 1 -d ' ')
    value=$(echo $line | cut -f 2- -d ' ')
    case "$action" in
        add)
            todos=("${todos[@]}" "$value")
            ;;
        delete)
            value=$((value-1))
            todos=("${todos[@]:0:$value}" "${todos[@]:$(($value + 1))}")
            ;;
        view)
            echo -e "\n*** TODO list ***"
            for i in $(seq ${#todos[@]}); do
                echo "$i: ${todos[$((i-1))]}"
            done
            echo
            ;;
    esac
done < todo.fifo

# vim: set et:sw=4:
