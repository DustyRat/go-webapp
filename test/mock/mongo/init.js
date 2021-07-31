db.createUser(
  {
    user: "mock",
    pwd: "password",
    roles: [
      { role: "readWrite", db: "Unit" }
    ]
  }
);

load("./docker-entrypoint-initdb.d/collection/Test.js");
load("./docker-entrypoint-initdb.d/data/Test.js");
