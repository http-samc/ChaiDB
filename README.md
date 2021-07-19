# ChaiDB ☕
⚡ A lightning fast NoSQL database in pure GO!

## Simple ✔️
Chai was designed around the goal of creating the simplest document structure possible. Here's the rundown:
- Imagine the database as your `C` drive. From the root, you can have any amount of either Documents (JSON Objects) or Collections ('Folders' of Documents)
- Collections are special high-level items. By default, Chai upcasts low-level lists into Collections
- Every Document has a unique `_id`. This can be provided client-side. If it's left nil, Chai automatically generates one. `_id`s **cannot** be repeated anywhere inside of a database.

## Getting Started 🌐
To interact with your database, start it on a URL (eg. `http://127.0.0.1:8080/mydb`) and `POST` your initial dataset. Our form can look something like:
```
{
    "2020-21": {
        "Gophers": {
            "wins": 10,
            "losses": 0,
            "games": [
                {"opp": "Pythons", "win": true, "num": 1},
                ...
            ]
        },
        ...
    }
}
```
If you get a `201 Created` response, that means that `mydb` now looks something like:
```
{
    "data": {
        "2020-21": {
            "Gophers": {
                "wins": 10,
                "losses": 0,
                "games": {
                    "60e50f1cfa612f1c864eb2a1": {
                        "opp": "Pythons",
                        "win": true,
                        "num": 1
                    },
                    ...
                }
            },
            ...
        }
    },
    "meta": {} // Metadata for ChaiDB, not part of DB
}
```
Notice that `games` went from a list to a Collection? That's Chai's upcasting. Because we didn't provide an `_id` value in the object body, Chai generated a new one for us. Here's what we'd have to send if we want to use our own `_id`:
```
{
    "opp": "Rubies",
    "win": true,
    "num": 2,
    "_id": "83n4n_ndk2" // custom _id
}
```
To get this into the database, we'll update the record with a `PUT` request to `http://127.0.0.1:8080/mydb/2020-21/Gophers`, using that as our body. Notice how you can query into a specific location, just like typing a path into your file explorer? Also note how the toplevel `data` key isn't called - Chai does this automatically. Here's what `mydb` looks like now:
```
{
    "data": {
        "2020-21": {
            "Gophers": {
                "wins": 10,
                "losses": 0,
                "games": {
                    "60e50f1cfa612f1c864eb2a1": {
                        "opp": "Pythons",
                        "win": true,
                        "num": 1
                    },
                    "83n4n_ndk2": {
                        "opp": "Rubies",
                        "win": true,
                        "num": 2
                    },
                    ...
                }
            },
            ...
        }
    },
    "meta": {} // Metadata for ChaiDB, not part of DB
}
```
The posted ID is 'pulled' out of the posted data and becomes the ID for it. Before we add more games, we should establish a schema so we can keep our dataset clean. Any Collection can have a schema, it can be inserted like a document, but it **must** have the predefined `_id` value of `$SCHEMA`.

We can add ours with a `PUT` request to `http://127.0.0.1:8080/mydb/2020-21/Gophers` using the following body:
```
{
    "_id": "$SCHEMA",
    "opp": "{str} {req}",
    "win": "{bool} {req}",
    "num": "{int} {req}"
}
```
Here's what `mydb` looks like now:
```
{
    "data": {
        "2020-21": {
            "Gophers": {
                "wins": 10,
                "losses": 0,
                "games": {
                    "60e50f1cfa612f1c864eb2a1": {
                        "opp": "Pythons",
                        "win": true,
                        "num": 1
                    },
                    "83n4n_ndk2": {
                        "opp": "Rubies",
                        "win": true,
                        "num": 2
                    },
                    "$SCHEMA": {
                        "opp": "{str} {req}",
                        "win": "{bool} {req}",
                        "num": "{int} {req}"
                    }
                }
            },
            ...
        }
    },
    "meta": {} // Metadata for ChaiDB, not part of DB
}
```
Now, any time a document is created/modified in the `2020-21/Gophers/games` route, the schema will be checked and if an error is found a `406 Not Acceptable` will be sent.

One benefit of lists is how easy they are to sort. Fortunately, the high-level Collections API also features a way to easily sort and cache objects in a collection. To request this serverside, all we have to do is send a `POST` request to `http://127.0.0.1:8080/mydb/$collections/sort` with the following body:
```
{
    "applyTo": "$#/2020-21/Gophers/games",
    "sortBy": "$~/num",
    "reverse": true,
    "return": -1,
    "sortName": "newest"
}
```
Let's breakdown what each of these mean:
- `applyTo` is a reference to a collection. As you saw with the `$SCHEMA` setup, anything starting with `$` is a special call to Chai's API. `$#` means 'put me in the top level data key'. After that, we just use the previously mentioned keypath to get us to the 'games' collection.
- `sortBy` is a reference to a key that will always appear within every Document in the targeted collection. The `$~` operator means 'take me all the way to my parent', and we are specifying that in the parent we want the `num` key when we say `$~/num`. It is highly reccommended you have a schema enforced prior to sorting.
- `reverse` is optional and tells Chai how to sort. By default, it goes low to high (numbers)/A-Z (letters). Since the `num` represent the order of the game in the season, with 1 being the first game they played, 2 being the second, etc. Since we want to get the newest games first (highest `num` value), we want `reverse` to be `true`.
- `return` is a required paramater that tells Chai how much of the list to return. If it is a positive int, the first N items will be returned (or the complete list if the value is larger than the list). To explicitly call all items in the collections, set it equal to -1.
- `sortName` is the name of this sorting configuration.

If you get a `201 Created` response, you can now access the sorted collection with a `GET` request to `http://127.0.0.1:8080/mydb/2020-21/Gophers`, **BUT** we want to use the sort we defined earlier. To do that we need to use the `?newest=true` parameter. Our final URL should look like `http://127.0.0.1:8080/mydb/2020-21/Gophers?newest=true`.

Here's what the response looks like:
```
[
    ... ,
    {
        "opp": "Rubies",
        "win": true,
        "num": 2
    },
    {
        "opp": "Pythons",
        "win": true,
        "num": 1
    }
]
```
Notice the response is a low-level list, sorted by the `num` key from high to low.