db = db.getSiblingDB('jeeb');

// Create collections
db.createCollection('users');
db.createCollection('workouts');
db.createCollection('studies');
db.createCollection('sleep');
db.createCollection('finance');
db.createCollection('events');
db.createCollection('integrations');

// Create indexes
db.users.createIndex({ "keycloak_id": 1 }, { unique: true });
db.users.createIndex({ "email": 1 }, { unique: true });

db.workouts.createIndex({ "user_id": 1, "created_at": -1 });
db.workouts.createIndex({ "user_id": 1, "type": 1 });

db.studies.createIndex({ "user_id": 1, "created_at": -1 });
db.studies.createIndex({ "user_id": 1, "subject": 1 });

db.sleep.createIndex({ "user_id": 1, "start_time": -1 });

db.finance.createIndex({ "user_id": 1, "date": -1 });
db.finance.createIndex({ "user_id": 1, "category": 1 });
db.finance.createIndex({ "user_id": 1, "type": 1 });

db.events.createIndex({ "user_id": 1, "start": 1 });
db.events.createIndex({ "user_id": 1, "type": 1 });

db.integrations.createIndex({ "user_id": 1, "provider": 1 }, { unique: true });

print('Jeeb database initialized');
