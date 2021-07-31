db = db.getSiblingDB('Example');
db.createCollection("Model",{});
collection = db.getCollection("Model");

collection.createIndex({ "createdTs": 1 }, { "background": true, "sparse": false, "unique": false });
collection.createIndex({ "updatedTs": 1 }, { "background": true, "sparse": false, "unique": false });

