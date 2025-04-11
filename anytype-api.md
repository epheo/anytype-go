

# Anytype API Documentation

## API Overview

- **Base URL**: `http://localhost:31009/v1`
- **Current API Version**: `2025-03-17` (passed in HTTP header `Anytype-Version`)
- **Authentication**: Bearer token (obtained through pairing process)

## Authentication Flow

1. **Display Code**: The desktop app shows a 4-digit verification code
   - Endpoint: `POST /auth/display_code?app_name={appName}`
   - Response: Contains `challenge_id` needed for token request

2. **Get Token**: Exchange the code for an API token
   - Endpoint: `POST /auth/token?challenge_id={challengeId}&code={code}`
   - Response: Returns auth token for subsequent requests

## Spaces

### Get Spaces
- Endpoint: `GET /spaces?offset={offset}&limit={limit}`
- Description: Lists all available spaces
- Pagination: Supports offset/limit pagination

### Get Space
- Endpoint: `GET /spaces/{spaceId}`
- Description: Retrieves details of a specific space

### Create Space
- Endpoint: `POST /spaces`
- Description: Creates a new space
- Payload: JSON with space configuration

## Objects

### Get Objects
- Endpoint: `GET /spaces/{spaceId}/objects?offset={offset}&limit={limit}`
- Description: Lists objects in a space
- Pagination: Supports offset/limit pagination

### Get Object
- Endpoint: `GET /spaces/{spaceId}/objects/{objectId}`
- Description: Retrieves a specific object with its details

### Create Object
- Endpoint: `POST /spaces/{spaceId}/objects`
- Description: Creates a new object in a space
- Payload: JSON with object data

### Delete Object
- Endpoint: `DELETE /spaces/{spaceId}/objects/{objectId}`
- Description: Deletes an object from a space

## Types

### Get Types
- Endpoint: `GET /spaces/{spaceId}/types?offset={offset}&limit={limit}`
- Description: Lists available types in a space
- Pagination: Supports offset/limit pagination

### Get Type
- Endpoint: `GET /spaces/{spaceId}/types/{typeId}`
- Description: Retrieves details of a specific type

## Templates

### Get Templates
- Endpoint: `GET /spaces/{spaceId}/types/{typeId}/templates?offset={offset}&limit={limit}`
- Description: Lists templates for a specific type
- Pagination: Supports offset/limit pagination

### Get Template
- Endpoint: `GET /spaces/{spaceId}/types/{typeId}/templates/{templateId}`
- Description: Retrieves a specific template

## Lists

### Get List Views
- Endpoint: `GET /spaces/{spaceId}/lists/{listId}/views?offset={offset}&limit={limit}`
- Description: Gets available views for a list
- Pagination: Supports offset/limit pagination

### Get Objects in List
- Endpoint: `GET /spaces/{spaceId}/lists/{listId}/{viewId}/objects?offset={offset}&limit={limit}`
- Description: Lists objects contained in a list view
- Pagination: Supports offset/limit pagination

### Add Objects to List
- Endpoint: `POST /spaces/{spaceId}/lists/{listId}/objects`
- Description: Adds objects to a list

### Remove Objects from List
- Endpoint: `DELETE /spaces/{spaceId}/lists/{listId}/objects/{objectId}`
- Description: Removes an object from a list

## Members

### Get Members
- Endpoint: `GET /spaces/{spaceId}/members?offset={offset}&limit={limit}`
- Description: Lists members of a space
- Pagination: Supports offset/limit pagination

### Get Member
- Endpoint: `GET /spaces/{spaceId}/members/{objectId}`
- Description: Retrieves details of a specific member

### Update Member
- Endpoint: `PATCH /spaces/{spaceId}/members/{objectId}`
- Description: Updates member information/permissions

## Search

### Global Search
- Endpoint: `POST /search?offset={offset}&limit={limit}`
- Description: Searches across all spaces
- Payload: JSON with search parameters

### Space Search
- Endpoint: `POST /spaces/{spaceId}/search?offset={offset}&limit={limit}`
- Description: Searches within a specific space
- Payload: JSON with search parameters

## Export

### Export Object
- Endpoint: `GET /spaces/{spaceId}/objects/{objectId}/{format}`
- Description: Exports an object in the specified format

## Notes

- Most endpoints returning collections support pagination with `offset` and `limit` parameters
- A standard response includes both data and pagination information
- The API requires the Anytype desktop application to be running locally
- Error handling includes version compatibility checks between the extension and desktop application


# Anytype API Examples with curl

Below are examples of how to interact with each Anytype API endpoint using curl commands. For all examples, replace:
- `YOUR_TOKEN` with your actual bearer token
- `SPACE_ID`, `OBJECT_ID`, etc. with actual IDs

## Authentication

### Display Code
```bash
curl -X POST "http://localhost:31009/v1/auth/display_code?app_name=RaycastExtension" \
  -H "Anytype-Version: 2025-03-17"
```
Response:
```json
{
  "challenge_id": "abc123xyz"
}
```

### Get Token
```bash
curl -X POST "http://localhost:31009/v1/auth/token?challenge_id=abc123xyz&code=1234" \
  -H "Anytype-Version: 2025-03-17"
```
Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

## Spaces

### Get Spaces
```bash
curl "http://localhost:31009/v1/spaces?offset=0&limit=10" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17"
```
Response:
```json
{
  "items": [
    {
      "id": "space123",
      "name": "Personal Space",
      "icon": "📓"
    }
  ],
  "total": 1,
  "limit": 10,
  "offset": 0
}
```

### Get Space
```bash
curl "http://localhost:31009/v1/spaces/space123" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17"
```
Response:
```json
{
  "id": "space123",
  "name": "Personal Space",
  "icon": "📓",
  "created_at": "2025-01-15T12:30:00Z",
  "updated_at": "2025-03-10T16:45:22Z"
}
```

### Create Space
```bash
curl -X POST "http://localhost:31009/v1/spaces" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Project X",
    "icon": "🚀"
  }'
```
Response:
```json
{
  "id": "space456",
  "name": "Project X",
  "icon": "🚀",
  "created_at": "2025-04-11T14:22:10Z"
}
```

## Objects

### Get Objects
```bash
curl "http://localhost:31009/v1/spaces/space123/objects?offset=0&limit=20" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17"
```
Response:
```json
{
  "items": [
    {
      "id": "obj1",
      "title": "My First Note",
      "type_id": "note"
    },
    {
      "id": "obj2",
      "title": "Project Roadmap",
      "type_id": "document"
    }
  ],
  "total": 42,
  "limit": 20,
  "offset": 0
}
```

### Get Object
```bash
curl "http://localhost:31009/v1/spaces/space123/objects/obj1" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17"
```
Response:
```json
{
  "id": "obj1",
  "title": "My First Note",
  "type_id": "note",
  "created_at": "2025-02-10T08:15:30Z",
  "updated_at": "2025-03-05T16:20:45Z",
  "content": "This is the content of my first note...",
  "relations": {
    "tags": ["important", "personal"]
  }
}
```

### Create Object
```bash
curl -X POST "http://localhost:31009/v1/spaces/space123/objects" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Meeting Notes",
    "type_id": "note",
    "content": "Topics to discuss: 1. Project timeline 2. Resource allocation",
    "relations": {
      "tags": ["meeting", "important"]
    }
  }'
```
Response:
```json
{
  "id": "obj3",
  "title": "Meeting Notes",
  "type_id": "note"
}
```

### Delete Object
```bash
curl -X DELETE "http://localhost:31009/v1/spaces/space123/objects/obj3" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17"
```
Response:
```json
{
  "success": true
}
```

## Types

### Get Types
```bash
curl "http://localhost:31009/v1/spaces/space123/types?offset=0&limit=10" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17"
```
Response:
```json
{
  "items": [
    {
      "id": "note",
      "name": "Note",
      "icon": "📝"
    },
    {
      "id": "task",
      "name": "Task",
      "icon": "✅"
    }
  ],
  "total": 8,
  "limit": 10,
  "offset": 0
}
```

### Get Type
```bash
curl "http://localhost:31009/v1/spaces/space123/types/note" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17"
```
Response:
```json
{
  "id": "note",
  "name": "Note",
  "icon": "📝",
  "relations": [
    {
      "key": "tags",
      "name": "Tags",
      "format": "tag"
    }
  ]
}
```

## Templates

### Get Templates
```bash
curl "http://localhost:31009/v1/spaces/space123/types/note/templates?offset=0&limit=5" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17"
```
Response:
```json
{
  "items": [
    {
      "id": "template1",
      "name": "Basic Note",
      "description": "Simple note template"
    },
    {
      "id": "template2",
      "name": "Meeting Notes",
      "description": "Template for meeting notes"
    }
  ],
  "total": 2,
  "limit": 5,
  "offset": 0
}
```

### Get Template
```bash
curl "http://localhost:31009/v1/spaces/space123/types/note/templates/template1" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17"
```
Response:
```json
{
  "id": "template1",
  "name": "Basic Note",
  "description": "Simple note template",
  "content": "# {{Title}}\n\n_Created on {{Date}}_\n\n## Notes\n\n",
  "default_relations": {
    "tags": []
  }
}
```

## Lists

### Get List Views
```bash
curl "http://localhost:31009/v1/spaces/space123/lists/list1/views?offset=0&limit=5" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17"
```
Response:
```json
{
  "items": [
    {
      "id": "view1",
      "name": "Default View",
      "type": "grid"
    },
    {
      "id": "view2",
      "name": "Timeline",
      "type": "calendar"
    }
  ],
  "total": 2,
  "limit": 5,
  "offset": 0
}
```

### Get Objects in List
```bash
curl "http://localhost:31009/v1/spaces/space123/lists/list1/view1/objects?offset=0&limit=10" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17"
```
Response:
```json
{
  "items": [
    {
      "id": "obj1",
      "title": "My First Note",
      "type_id": "note"
    },
    {
      "id": "obj2",
      "title": "Project Roadmap",
      "type_id": "document"
    }
  ],
  "total": 15,
  "limit": 10,
  "offset": 0
}
```

### Add Objects to List
```bash
curl -X POST "http://localhost:31009/v1/spaces/space123/lists/list1/objects" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17" \
  -H "Content-Type: application/json" \
  -d '{
    "object_ids": ["obj4", "obj5"]
  }'
```
Response:
```json
{
  "success": true,
  "added": 2
}
```

### Remove Objects from List
```bash
curl -X DELETE "http://localhost:31009/v1/spaces/space123/lists/list1/objects/obj4" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17"
```
Response:
```json
{
  "success": true
}
```

## Members

### Get Members
```bash
curl "http://localhost:31009/v1/spaces/space123/members?offset=0&limit=10" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17"
```
Response:
```json
{
  "items": [
    {
      "id": "user1",
      "name": "Alice",
      "role": "admin"
    },
    {
      "id": "user2",
      "name": "Bob",
      "role": "editor"
    }
  ],
  "total": 2,
  "limit": 10,
  "offset": 0
}
```

### Get Member
```bash
curl "http://localhost:31009/v1/spaces/space123/members/user1" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17"
```
Response:
```json
{
  "id": "user1",
  "name": "Alice",
  "role": "admin",
  "joined_at": "2025-01-15T12:30:00Z"
}
```

### Update Member
```bash
curl -X PATCH "http://localhost:31009/v1/spaces/space123/members/user2" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17" \
  -H "Content-Type: application/json" \
  -d '{
    "role": "admin"
  }'
```
Response:
```json
{
  "id": "user2",
  "name": "Bob",
  "role": "admin"
}
```

## Search

### Global Search
```bash
curl -X POST "http://localhost:31009/v1/search?offset=0&limit=20" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "project",
    "types": ["note", "document"]
  }'
```
Response:
```json
{
  "items": [
    {
      "id": "obj2",
      "title": "Project Roadmap",
      "type_id": "document",
      "space_id": "space123"
    },
    {
      "id": "obj6",
      "title": "Project Timeline",
      "type_id": "note",
      "space_id": "space456"
    }
  ],
  "total": 2,
  "limit": 20,
  "offset": 0
}
```

### Space Search
```bash
curl -X POST "http://localhost:31009/v1/spaces/space123/search?offset=0&limit=10" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "meeting",
    "types": ["note"]
  }'
```
Response:
```json
{
  "items": [
    {
      "id": "obj3",
      "title": "Meeting Notes",
      "type_id": "note"
    }
  ],
  "total": 1,
  "limit": 10,
  "offset": 0
}
```

## Export

### Export Object
```bash
curl "http://localhost:31009/v1/spaces/space123/objects/obj1/pdf" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17" \
  -o "exported-note.pdf"
```

This will download the object as a PDF file named "exported-note.pdf" to your current directory.

Alternate formats may include:
```bash
# Export as markdown
curl "http://localhost:31009/v1/spaces/space123/objects/obj1/markdown" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17" \
  -o "exported-note.md"

# Export as HTML
curl "http://localhost:31009/v1/spaces/space123/objects/obj1/html" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17" \
  -o "exported-note.html"
