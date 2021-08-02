db = db.getSiblingDB('Example');
db.createCollection("Document",{});
collection = db.getCollection("Document");

collection.createIndex({ "createdTs": 1 }, { "background": true, "sparse": false, "unique": false });
collection.createIndex({ "updatedTs": 1 }, { "background": true, "sparse": false, "unique": false });
