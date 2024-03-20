JSON Untuk POST /orders

[
    {
        "id": 1,
        "customerName": "User1",
        "orderedAt": "2021-10-06T16:53:27.675931+07:00",
        "items": [
            {
                "id": 1,
                "itemCode": "Itssem1",
                "description": "ItemDescription1",
                "quantity": 1
            },
            {
                "id": 2,
                "itemCode": "Itessm2",
                "description": "ItemDescription2",
                "quantity": 1
            }
        ]
    },
    {
        "id": 2,
        "customerName": "User2",
        "orderedAt": "2021-10-06T16:53:27.675931+07:00",
        "items": [
            {
                "id": 3,
                "itemCode": "Itemss3",
                "description": "ItemDescription3",
                "quantity": 1
            },
            {
                "id": 4,
                "itemCode": "Itemss4",
                "description": "ItemDescription4",
                "quantity": 1
            }
        ]
    }
]


JSON untuk PATCH /orders

[
    {
        "id": 1,
        "customerName": "Amanda",
        "orderedAt": "2021-10-06T16:53:27.675931+07:00",
        "items": [
            {
                "id": 1,
                "itemCode": "Durhen",
                "description": "Durian Hendri",
                "quantity": 1
            },
            {
                "id": 2,
                "itemCode": "Tasma",
                "description": "Sucita",
                "quantity": 1
            }
        ]
    },
    {
        "id": 2,
        "customerName": "Lex Luthor",
        "orderedAt": "2021-10-06T16:53:27.675931+07:00",
        "items": [
            {
                "id": 3,
                "itemCode": "Dang Dude",
                "description": "Dude Dang",
                "quantity": 1
            },
            {
                "id": 4,
                "itemCode": "AstLuthor",
                "description": "Yeght",
                "quantity": 1
            }
        ]
    }
]

expected output:

[
    {
        "id": 1,
        "orderedAt": "2021-10-06T16:53:27.675931+07:00",
        "customerName": "Amanda",
        "items": [
            {
                "id": 1,
                "itemCode": "Durhen",
                "description": "Durian Hendri",
                "quantity": 1
            },
            {
                "id": 2,
                "itemCode": "Tasma",
                "description": "Sucita",
                "quantity": 1
            }
        ]
    },
    {
        "id": 2,
        "orderedAt": "2021-10-06T16:53:27.675931+07:00",
        "customerName": "Lex Luthor",
        "items": [
            {
                "id": 3,
                "itemCode": "Dang Dude",
                "description": "Dude Dang",
                "quantity": 1
            },
            {
                "id": 4,
                "itemCode": "AstLuthor",
                "description": "Yeght",
                "quantity": 1
            }
        ]
    }
]



JSON untuk POST /items

{
    "items":
[
    {
        "itemCode":"Makanan",
        "description":"Untuk dimakan",
        "quantity":1,
        "ordersId":1

    },
    {
     "itemCode":"Minuman",
        "description":"Untuk diminum",
        "quantity":1,
        "ordersId":2
    }
]
}

expected output:

[
    {
        "itemCode": "Makanan",
        "description": "Untuk dimakan",
        "quantity": 1
    },
    {
        "itemCode": "Minuman",
        "description": "Untuk diminum",
        "quantity": 1
    }
]

JSON untuk /users

[{
"customerName":"Ika"
},
{
"customerName":"Lobi"
}]

expected output:

[
    {
        "id": 3,
        "orderedAt": "2024-03-11T17:07:48.2671286+07:00",
        "customerName": "Ika",
        "items": null
    },
    {
        "id": 4,
        "orderedAt": "2024-03-11T17:07:48.2835737+07:00",
        "customerName": "Lobi",
        "items": null
    }
]

designed as null, only for user creation, please use batch /orders instead

/orders/:id [delete]

to delete using ID

please use /orders/:id [get] in order to get specific ID