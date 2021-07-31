db = db.getSiblingDB('Example')
var collection = db.getCollection("Model");

/*
*   _id's starting with 0 - get
*   _id's starting with 1 - find
*   _id's starting with 2 - update
*   _id's starting with 3 - delete
*/

collection.insert({
    "_id": ObjectId("000000000000000000000001"),
    // TODO: Add fields
})

collection.insert({
    "_id": ObjectId("100000000000000000000001"),
    // TODO: Add fields
})

collection.insert({
    "_id": ObjectId("200000000000000000000001"),
    // TODO: Add fields
})

collection.insert({
    "_id": ObjectId("300000000000000000000001"),
    // TODO: Add fields
})
