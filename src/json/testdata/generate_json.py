import json

with open('array_of_int.json', 'w') as file:
    file.write(
        json.dumps(
            [123312 for _ in range(100000)]
        )
    )

with open('array_of_floats.json', 'w') as file:
    file.write(
        json.dumps(
            [123312.12243e5 for _ in range(100000)]
        )
    )


with open('map_of_string.json', 'w') as file:
    file.write(
        json.dumps(
            ['Hello, my name is '+str(i) for i in range(100000)]
        )
    )


with open('map_of_string.json', 'w') as file:
    file.write(
        json.dumps(
            {'k'+str(i): "value"+str(i) for i in range(10000)}
        )
    )


with open('big_string.json', 'w') as file:
    file.write(
        json.dumps(
            "abcdefghijklmnopqrstuvwxyz"*100000
        )
    )
