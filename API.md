
# API

The API always expects and replies back JSON-formatted data. Note that if JSON
data is not given to the API, or if expected data is missing, the API will
simply reply with an empty page and a 400 Bad Request error.
It consists of a number of endpoints that each focus on certain aspects of the
todo application.

## `token` endpoint

A token is simply a string that represents a user's "session" - although a session
may be longer than a typical "session".

A primary token is created when a user logs in using the username and password.
A secondary or tertiary token is created when the user authorizes other applications
to create/modify todo's under their name.

A primary token can create only secondary or tertiary tokens, and secondary or
tertiary tokens cannot create any tokens. When a primary token is invalidated,
the tokens created by the primary token will still stay valid.

Anyone can view the detailed informaiton for a public todo.

Tertiary tokens only have permissions to create new todos under the user id of their
owner's user. They cannot modify any todo. Secondary tokens have similar permissions,
but they can also read the data of and modify a todo that is owned by the user id
of their owner's user. Primary tokens have the additional permissions of being
able to list all todos under the user id of their owner's user and delete todos.
Basically, primary endpoints have access to all token-authorized endpoints.

Throughout the API, the token type is represented as an integer:

* `0` is reserved
* `1` is primary
* `2` is secondary
* `3` is tertiary
* `9` is invalid

Some other endpoints require authenication from a token, provided as the `authority`
field.

### Check the type of a token

```
POST /api/token/type
````

#### Parameters

|Name|Type|Description|
|----|----|-----------|
|`token`|`string`|The token to test.|

#### Behaviour

* Return the type of the token.

#### Response

|Name|Type|Description|
|----|----|-----------|
|`type`|`int`|The type of the token.|

### Create a new token

```
POST /api/token/new
````

#### Parameters

|Name|Type|Description|
|----|----|-----------|
|`type`|`int`|The type of the token to be created.|
|`username`|`string`?|The username of the user this token will belong to.|
|`password`|`string`?|The password of the user this token will belong to.|
|`authority`|`string`?|A primary token.|

#### Behaviour

* If `type` is a primary, `username` and `password` are both present and both are valid, create a new primary token.
* Else if `type` is primary or tertiary
    * If `username` and `password` are both present and both are valid, create a new token of the requested type.
    * If `authority` is present and is a valid primary token, create a new token of the requested type.
* Else, return an error 400 or 403.

#### Response

|Name|Type|Description|
|----|----|-----------|
|`error`|`string`?|If an error occurred, this field is present and a friendly error message is filled in appropriately.|
|`token`|`string`?|If no error occurs, this field is present and contains the newly authorized token.|

### Invalidate a token

```
POST /api/token/invalidate
````

#### Parameters

|Name|Type|Description|
|----|----|-----------|
|`token`|`string`|The token to invalidate.|
|`username`|`string`?|The username of the user this token belongs to.|
|`password`|`string`?|The password of the user this token belongs to.|
|`authority`|`string`?|A primary token.|

#### Behaviour

* If `username` and `password` are present and valid, invalidate the token.
* Else if `authority` is present
    * If `token` is a primary token and equal to the contents of `authority`, invalidate the token.
    * Else if `token` is a secondary or tertiary token and `authority` is a primary token, invalidate the token.
    * Else, return an error 400 or 403.
* Else, return an error 400.

#### Response

|Name|Type|Description|
|----|----|-----------|
|`error`|`string`?|If an error occurred, this field is present and a friendly error message is filled in appropriately.|

## `todo` endpoint

A todo is represented in JSON using the following format:

|Name|Type|Description|
|----|----|-----------|
|`id`|`int`|The unique identifier for the todo item.|
|`state`|`int`|The current state of the todo.|
|`tag_id`|`int`|The ID of the tag of this todo.|
|`owner_id`|`int`?|The ID of the owner of this todo. Present only on detailed information.|
|`public`|`boolean`?|Whether or not this todo is public. Present only on detailed information.|
|`name`|`string`|The short name of the todo item. Max 256 characters.|
|`description`|`string`?|The in-depth description of this todo. Present only on detailed information.|

### Create a new or update an existing todo

```
POST /api/todo/update
```

#### Parameters

|Name|Type|Description|
|----|----|-----------|
|`todo`|`todo`|A todo item (see above). Field `owner_id` is ignored.|
|`authority`|`string`|A token.|

#### Behaviour

* If `todo.id == -1` and if `authority` is a present and a valid primary, secondary, or tertiary token, and `todo` is a valid todo, create a new todo under the `owner_id` of the `owner_id` of this token.
* Else if `todo.id >= 0`, `authority` is a present and valid primary or secondary token, `todo` is a valid todo, and a todo that is owned by the owner of the token and that has the given ID exists in the database, update the database.
* Else, return an error.

#### Response

|Name|Type|Description|
|----|----|-----------|
|`error`|`string`?|If an error occurred, this field is present and a friendly error message is filled in appropriately.|

### Remove an existing todo

```
POST /api/todo/remove
```

#### Parameters

|Name|Type|Description|
|----|----|-----------|
|`todo`|`todo`|A todo item. All fields except for `id` are ignored.|
|`authority`|`string`|A primary token.|

#### Behaviour

* If `authority` is a present and a valid primary token, the todo id exists in the database, and the todo associated with the todo id is owned by the token's owner, delete the todo.
* Else, return an error.

#### Response

|Name|Type|Description|
|----|----|-----------|
|`error`|`string`?|If an error occurred, this field is present and a friendly error message is filled in appropriately.|

### Get information on an existing todo

```
POST /api/todo/info
```

#### Parameters

|Name|Type|Description|
|----|----|-----------|
|`todo`|`todo`|A todo item. All fields except for `id` are ignored.|
|`authority`|`string`?|A primary token.|

#### Behaviour

* If `token` is a present and a valid primary or secondary token, the todo id exists in the database, and the todo associated with the todo id is owned by the token's owner, return detailed information on the todo.
* Else If `token` is not present, the todo id exists in the database, and the todo is public, return detailed information on the todo.
* Else, return an error.

#### Response

|Name|Type|Description|
|----|----|-----------|
|`error`|`string`?|If an error occurred, this field is present and a friendly error message is filled in appropriately.|
|`todo`|`todo`?|If no error occurred, this field is present and contains an in-depth todo item.|


## `todos` endpoint

The `todos` endpoint deals with batch fetches of `todos`. It will omit some fields
from the returned results depending on the authority of the possibly supplied token.

### Get a list of todos

```
GET /api/todos/list
POST /api/todos/list
```

#### Parameters

|Name|Type|Description|
|----|----|-----------|
|`authority`|`string`?|A primary token.|

#### Behaviour

* If `token` is a valid primary token, return the todos owned by this user and public todos.
* Else, return public todos.

#### Response

|Name|Type|Description|
|----|----|-----------|
|`todos`|`todo[]`|Todos that match the query.|


## `tags` endpoint

The `tags` endpoint deals with tags. At the moment, the API only supports
listing all tags; addition and removal of tags is not yet implemented.

A valid primary token must be supplied in order to list tokens. This is to
prevent DoS attacks.

A tag is represented in JSON using the following format:

|Name|Type|Description|
|----|----|-----------|
|`id`|`int`|ID of the tag.|
|`name`|`string`|The name of the tag. Max 16 characters.|

### Get a list of tags

```
GET /api/tags/list
```

#### Parameters

None.

#### Behaviour

Return a list of tags.


#### Response

|Name|Type|Description|
|----|----|-----------|
|`tags`|`tag[]`|An array of the tags in the database.|
