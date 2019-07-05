#!/bin/bash

iteration="$(curl -s https://hv59yay1u3.execute-api.eu-central-1.amazonaws.com/live/iteration/current)"

function gen() {
    id="$(echo "${1}" | jq -r '.[0]')"
    name="$(echo "${1}" | jq -r '.[1]')"

    obj="$(jq -n --arg id "${id}" --arg name "${name}" '{ShipID: ($id | tonumber), ShipName: $name, InDivision: '$2'}')"

    res="$(curl -s -H 'Content-Type: application/json' -X POST --data "${obj}" http://localhost:1323/iterations/current/battles)"

    fin_obj="$(echo "${res}" | jq '. + {Status: "finished"} | .Statistics=(.Statistics + {Win: '$3', Damage: {Value: '$4'}, Kills: {Value: '$5'}, Survived: '$6'}) | .')"
    curl -s -H 'Content-Type: application/json' -X POST --data "${fin_obj}" http://localhost:1323/iterations/current/battles/active
}

IFS=$'\n'
for ship in $(echo "${iteration}" | jq -c '.ships[] | [.id,.name]'); do
    num="$(shuf -i 2-30 -n 1)"

    for it in $(seq 1 $num); do
        gen $ship "$(shuf -e true false -n 1)" "$(shuf -e true false -n 1)" "$(shuf -i 10350-133978 -n 1)" "$(shuf -i 0-5 -n 1)" "$(shuf -e true false -n 1)"
    done
done
